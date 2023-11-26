package model

import (
	"time"

	"github.com/cbguder/weather/noaa"
)

const (
	tminF = 50
	tmaxF = 80
)

type dayRecord struct {
	Tmin int
	Tmax int
	Prcp int
}

func Score(records []noaa.DailyRecord) (int, float32) {
	numRecords := 0
	goodDays := 0

	dayRecords := make(map[time.Time]dayRecord)

	for _, record := range records {
		if _, ok := dayRecords[record.Date]; !ok {
			dayRecords[record.Date] = dayRecord{
				Tmin: -9999,
				Tmax: -9999,
				Prcp: -9999,
			}
		}

		dayRecord := dayRecords[record.Date]

		if record.Element == "TMAX" {
			dayRecord.Tmax = record.Value
		} else if record.Element == "TMIN" {
			dayRecord.Tmin = record.Value
		} else if record.Element == "PRCP" {
			dayRecord.Prcp = record.Value
		} else {
			continue
		}

		dayRecords[record.Date] = dayRecord
	}

	for _, dayRecord := range dayRecords {
		if dayRecord.Tmin == -9999 || dayRecord.Tmax == -9999 {
			continue
		}

		numRecords++

		if isGoodDay(dayRecord) {
			goodDays++
		}
	}

	if numRecords == 0 {
		return 0, 0
	}

	score := 100.0 * float32(goodDays) / float32(numRecords)
	return numRecords, score
}

// Convert Fahrenheit to Celsius (tenths of a degree)
func f2c(f int) int {
	c := float32(f-32) / 1.8
	return int(c * 10)
}

func isGoodDay(day dayRecord) bool {
	tmin, tmax := f2c(tminF), f2c(tmaxF)

	if day.Tmin < tmin {
		return false
	}

	if day.Tmax > tmax {
		return false
	}

	if day.Prcp > 0 {
		return false
	}

	return true
}
