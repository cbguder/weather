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

type csvRecord struct {
	StationId string
	Date      string
	Element   string
	Value     int
	MFlag     string
	QFlag     string
	SFlag     string
	ObsTime   string
}

type DailyRecord struct {
	StationId string
	Date      string
	Elements  []ElementRecord
}

func (r DailyRecord) Element(element string) (ElementRecord, bool) {
	for _, e := range r.Elements {
		if e.Element == element {
			return e, true
		}
	}

	return ElementRecord{}, false
}

type ElementRecord struct {
	Element string
	Value   int
	MFlag   string
	QFlag   string
	SFlag   string
	ObsTime string
}

func PreloadRecordsForStations(stationIds []string) error {
	paths := make([]string, len(stationIds))

	for i, stationId := range stationIds {
		paths[i] = fmt.Sprintf("by_station/%s.csv.gz", stationId)
	}

	return preloadDataFiles(paths)
}

func RecordsForStation(stationId string, startDate, endDate *time.Time) ([]DailyRecord, error) {
	path := fmt.Sprintf("by_station/%s.csv.gz", stationId)

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

	records, err := readCsvRecords(gz, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return groupDailyRecords(records), nil
}

func readCsvRecords(r io.Reader, startDate, endDate *time.Time) ([]csvRecord, error) {
	var records []csvRecord

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

		record := csvRecord{
			StationId: rec[0],
			Date:      rec[1],
			Element:   rec[2],
			Value:     value,
			MFlag:     rec[4],
			QFlag:     rec[5],
			SFlag:     rec[6],
			ObsTime:   rec[7],
		}

		records = append(records, record)
	}

	sort.Slice(records, func(i, j int) bool {
		if records[i].Date == records[j].Date {
			return records[i].Element < records[j].Element
		}

		return records[i].Date < records[j].Date
	})

	return records, nil
}

func groupDailyRecords(rawRecords []csvRecord) []DailyRecord {
	var records []DailyRecord

	for _, raw := range rawRecords {
		var rec DailyRecord

		if len(records) == 0 || records[len(records)-1].Date != raw.Date {
			rec = DailyRecord{
				StationId: raw.StationId,
				Date:      raw.Date,
			}

			records = append(records, rec)
		} else {
			rec = records[len(records)-1]
		}

		rec.Elements = append(rec.Elements, ElementRecord{
			Element: raw.Element,
			Value:   raw.Value,
			MFlag:   raw.MFlag,
			QFlag:   raw.QFlag,
			SFlag:   raw.SFlag,
			ObsTime: raw.ObsTime,
		})

		records[len(records)-1] = rec
	}

	return records
}
