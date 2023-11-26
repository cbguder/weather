package noaa

import (
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"
)

type DailyRecord struct {
	StationId string
	Date      time.Time
	Element   string
	Value     int
	MFlag     string
	QFlag     string
	SFlag     string
	ObsTime   string
}

func RecordsForStation(stationId string) ([]DailyRecord, error) {
	path := fmt.Sprintf("by_station/%s.csv.gz", stationId)

	recordsFile, err := openDataFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	defer recordsFile.Close()

	gz, err := gzip.NewReader(recordsFile)
	if err != nil {
		log.Fatalln(err)
	}

	defer gz.Close()

	return readDailyRecords(gz)
}

func readDailyRecords(r io.Reader) ([]DailyRecord, error) {
	var records []DailyRecord

	cr := csv.NewReader(r)
	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		record := DailyRecord{
			StationId: rec[0],
			Date:      parseDate(rec[1]),
			Element:   rec[2],
			Value:     parseInt(rec[3]),
			MFlag:     rec[4],
			QFlag:     rec[5],
			SFlag:     rec[6],
			ObsTime:   rec[7],
		}

		records = append(records, record)
	}

	return records, nil
}

func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalln(err)
	}
	return i
}

func parseDate(s string) time.Time {
	t, err := time.Parse("20060102", s)
	if err != nil {
		log.Fatalln(err)
	}
	return t
}
