package noaa

import (
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"sort"
	"strconv"
	"time"
)

type dailyCsvRecord struct {
	StationId string
	Date      time.Time
	Element   string
	Value     int
}

type DailyRecord struct {
	StationId     string
	Date          time.Time
	Elements      map[string]ElementRecord
	HourlyRecords []HourlyRecord
}

type ElementRecord struct {
	Element string
	Value   int
}

func PreloadDailyRecords(stationIds []string) error {
	paths := make([]string, len(stationIds))

	for i, stationId := range stationIds {
		paths[i] = dailyDataPath(stationId)
	}

	return preloadDataFiles(paths)
}

func DailyRecords(stationId string, startDate, endDate *time.Time) ([]DailyRecord, error) {
	path := dailyDataPath(stationId)
	recordsFile, err := openDataFile(path)
	if err != nil {
		return nil, err
	}

	defer recordsFile.Close()

	gz, err := gzip.NewReader(recordsFile)
	if err != nil {
		return nil, err
	}

	defer gz.Close()

	records, err := readDailyCsvRecords(gz, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return groupDailyRecords(records), nil
}

func dailyDataPath(stationId string) string {
	return fmt.Sprintf("/pub/data/ghcn/daily/by_station/%s.csv.gz", stationId)
}

func readDailyCsvRecords(r io.Reader, startDate, endDate *time.Time) ([]dailyCsvRecord, error) {
	var records []dailyCsvRecord

	cr := csv.NewReader(r)
	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		value, err := strconv.Atoi(rec[3])
		if err != nil {
			return nil, err
		}

		if value == -9999 {
			continue
		}

		date, err := time.Parse("20060102", rec[1])
		if err != nil {
			return nil, err
		}

		if startDate != nil && startDate.After(date) {
			continue
		}

		if endDate != nil && endDate.Before(date) {
			continue
		}

		record := dailyCsvRecord{
			StationId: rec[0],
			Date:      date,
			Element:   rec[2],
			Value:     value,
		}

		records = append(records, record)
	}

	sort.Slice(records, func(i, j int) bool {
		cmp := records[i].Date.Compare(records[j].Date)

		if cmp == 0 {
			return records[i].Element < records[j].Element
		}

		return cmp == -1
	})

	return records, nil
}

func groupDailyRecords(rawRecords []dailyCsvRecord) []DailyRecord {
	var records []DailyRecord

	for _, raw := range rawRecords {
		rec := DailyRecord{
			Elements: make(map[string]ElementRecord),
		}

		if len(records) == 0 || !records[len(records)-1].Date.Equal(raw.Date) {
			rec = DailyRecord{
				StationId: raw.StationId,
				Date:      raw.Date,
				Elements:  make(map[string]ElementRecord),
			}

			records = append(records, rec)
		} else {
			rec = records[len(records)-1]
		}

		rec.Elements[raw.Element] = ElementRecord{
			Element: raw.Element,
			Value:   raw.Value,
		}

		records[len(records)-1] = rec
	}

	return records
}
