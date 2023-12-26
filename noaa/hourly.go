package noaa

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type HourlyRecord struct {
	StationId   string
	StationName string
	Time        time.Time
	ReportType  string
	Temp        int
	Prcp        int
}

func PreloadHourlyRecords(stationId string, startDate, endDate time.Time) error {
	var paths []string

	for year := startDate.Year(); year <= endDate.Year(); year++ {
		path := hourlyDataPath(year, stationId)
		paths = append(paths, path)
	}

	return preloadDataFiles(paths)
}

func HourlyRecords(stationId string, startDate, endDate time.Time, loc *time.Location) ([]DailyRecord, error) {
	var records []DailyRecord

	for year := startDate.Year(); year <= endDate.Year(); year++ {
		path := hourlyDataPath(year, stationId)

		recordsFile, err := openDataFile(path)
		if err != nil {
			return nil, err
		}

		defer recordsFile.Close()

		page, err := readHourlyRecords(recordsFile, startDate, endDate, loc)
		if err != nil {
			return nil, err
		}

		grouped := groupHourlyRecords(stationId, page)

		records = append(records, grouped...)
	}

	return records, nil
}

func hourlyDataPath(year int, stationId string) string {
	return fmt.Sprintf("/data/global-hourly/access/%d/%s.csv", year, stationId)
}

func groupHourlyRecords(stationId string, hourlies []HourlyRecord) []DailyRecord {
	var (
		records  []DailyRecord
		cur      DailyRecord
		lastDate time.Time
	)

	for _, record := range hourlies {
		date := time.Date(
			record.Time.Year(),
			record.Time.Month(),
			record.Time.Day(),
			0, 0, 0, 0,
			record.Time.Location(),
		)

		if date == lastDate {
			cur.HourlyRecords = append(cur.HourlyRecords, record)
			continue
		}

		lastDate = date

		if len(cur.HourlyRecords) > 0 {
			records = append(records, cur)
		}

		cur = DailyRecord{
			StationId:     stationId,
			Date:          date,
			HourlyRecords: []HourlyRecord{record},
		}
	}

	if len(cur.HourlyRecords) > 0 {
		records = append(records, cur)
	}

	return records
}

func readHourlyRecords(r io.Reader, startDate, endDate time.Time, loc *time.Location) ([]HourlyRecord, error) {
	var records []HourlyRecord

	cr := csv.NewReader(r)

	header, err := cr.Read()
	if err != nil {
		return nil, err
	}

	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		var record HourlyRecord

		for i, value := range rec {
			switch header[i] {
			case "STATION":
				record.StationId = value
			case "NAME":
				record.StationName = value
			case "REPORT_TYPE":
				record.ReportType = value
			case "DATE":
				record.Time, err = time.ParseInLocation("2006-01-02T15:04:05", value, loc)
				if err != nil {
					return nil, err
				}
			case "TMP":
				record.Temp, err = parseTemp(value)
				if err != nil {
					return nil, err
				}
			case "AA1":
				record.Prcp, err = parsePrcp(value)
				if err != nil {
					return nil, err
				}
			}
		}

		if record.ReportType != "FM-13" && record.ReportType != "FM-15" {
			continue
		}

		if record.Time.IsZero() {
			continue
		}

		if record.Temp == 9999 {
			continue
		}

		if startDate.After(record.Time) || endDate.Before(record.Time) {
			continue
		}

		records = append(records, record)
	}

	return records, nil
}

func parseTemp(value string) (int, error) {
	if value == "+9999,9" {
		return 9999, nil
	}

	parts := strings.Split(value, ",")
	if len(parts) == 0 {
		return 0, nil
	}

	return strconv.Atoi(parts[0])
}

func parsePrcp(value string) (int, error) {
	parts := strings.Split(value, ",")
	if len(parts) == 0 || parts[0] != "01" {
		return 0, nil
	}

	return strconv.Atoi(parts[1])
}
