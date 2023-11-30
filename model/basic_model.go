package model

import "github.com/cbguder/weather/noaa"

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

type basicModel struct {
	numRecords uint
	goodDays   uint
}

func (m *basicModel) Add(record noaa.DailyRecord) {
	if !m.isValidDay(record) {
		return
	}

	m.numRecords++

	if m.isGoodDay(record) {
		m.goodDays++
	}
}

func (m *basicModel) Scorecard() Scorecard {
	if m.numRecords == 0 {
		return Scorecard{}
	}

	score := 100.0 * float32(m.goodDays) / float32(m.numRecords)

	return Scorecard{
		Records: m.numRecords,
		Score:   score,
	}
}

func (m *basicModel) isValidDay(record noaa.DailyRecord) bool {
	if _, ok := record.Element("TMIN"); !ok {
		return false
	}

	if _, ok := record.Element("TMAX"); !ok {
		return false
	}

	return true
}

func (m *basicModel) isGoodDay(record noaa.DailyRecord) bool {
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

// Convert Fahrenheit to Celsius (tenths of a degree)
func f2c(f int) int {
	c := float32(f-32) / 1.8
	return int(c * 10)
}
