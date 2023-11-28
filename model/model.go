package model

import (
	"github.com/cbguder/weather/noaa"
)

const (
	minTempF = 45
	maxTempF = 85
	maxPrcp  = 10
)

var (
	minTempC int
	maxTempC int
)

func init() {
	minTempC, maxTempC = f2c(minTempF), f2c(maxTempF)
}

type Scorecard struct {
	Records int
	Score   float32
}

func Score(records []noaa.DailyRecord) Scorecard {
	numRecords := 0
	goodDays := 0

	for _, record := range records {
		if !isValidDay(record) {
			continue
		}

		numRecords++

		if isGoodDay(record) {
			goodDays++
		}
	}

	if numRecords == 0 {
		return Scorecard{}
	}

	score := 100.0 * float32(goodDays) / float32(numRecords)
	return Scorecard{
		Records: numRecords,
		Score:   score,
	}
}

// Convert Fahrenheit to Celsius (tenths of a degree)
func f2c(f int) int {
	c := float32(f-32) / 1.8
	return int(c * 10)
}

// Convert Celsius (tenths of a degree) to Fahrenheit
func c2f(c int) int {
	f := float32(c)/10*1.8 + 32
	return int(f)
}

func isValidDay(record noaa.DailyRecord) bool {
	if _, ok := record.Element("TMIN"); !ok {
		return false
	}

	if _, ok := record.Element("TMAX"); !ok {
		return false
	}

	return true
}

func isGoodDay(record noaa.DailyRecord) bool {
	tmin, _ := record.Element("TMIN")
	tmax, _ := record.Element("TMAX")
	prcp, _ := record.Element("PRCP")

	if tmin.Value < minTempC {
		return false
	}

	if tmax.Value > maxTempC {
		return false
	}

	if prcp.Value > maxPrcp {
		return false
	}

	return true
}
