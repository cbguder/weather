package model

import (
	"github.com/cbguder/weather/noaa"
)

const (
	minDailyTempF = 45
	maxDailyTempF = 85

	minHourlyTempF = 50
	maxHourlyTempF = 85

	maxDailyPrcp  = 10
	maxHourlyPrcp = 2
)

var (
	minDailyTempC  int
	maxDailyTempC  int
	minHourlyTempC int
	maxHourlyTempC int
)

func init() {
	minDailyTempC, maxDailyTempC = f2c(minDailyTempF), f2c(maxDailyTempF)
	minHourlyTempC, maxHourlyTempC = f2c(minHourlyTempF), f2c(maxHourlyTempF)
}

type basicModel struct {
	validRecords uint
	goodRecords  uint
}

func (m *basicModel) Add(record noaa.DailyRecord) {
	if !m.isValidDay(record) {
		return
	}

	m.validRecords++

	if m.isGoodDay(record) {
		m.goodRecords++
	}
}

func (m *basicModel) Scorecard() Scorecard {
	if m.validRecords == 0 {
		return Scorecard{}
	}

	score := 100.0 * float32(m.goodRecords) / float32(m.validRecords)

	return Scorecard{
		Records: m.validRecords,
		Score:   score,
	}
}

func (m *basicModel) isValidDay(record noaa.DailyRecord) bool {
	if len(record.HourlyRecords) > 0 {
		return true
	}

	if _, ok := record.Elements["TMIN"]; !ok {
		return false
	}

	if _, ok := record.Elements["TMAX"]; !ok {
		return false
	}

	return true
}

func (m *basicModel) isGoodDay(record noaa.DailyRecord) bool {
	if len(record.HourlyRecords) > 0 {
		return m.isGoodDayHourly(record.HourlyRecords)
	}

	return m.isGoodDayDaily(record)
}

func (m *basicModel) isGoodDayDaily(record noaa.DailyRecord) bool {
	tmin, _ := record.Elements["TMIN"]
	tmax, _ := record.Elements["TMAX"]
	prcp, _ := record.Elements["PRCP"]

	if tmin.Value < minDailyTempC {
		return false
	}

	if tmax.Value > maxDailyTempC {
		return false
	}

	if prcp.Value > maxDailyPrcp {
		return false
	}

	return true
}

func (m *basicModel) isGoodDayHourly(records []noaa.HourlyRecord) bool {
	seenHours := make(map[int]struct{})
	goodRun := 0

	for _, record := range records {
		hour := record.Time.Hour()
		if _, ok := seenHours[hour]; ok {
			continue
		}

		seenHours[hour] = struct{}{}

		if hour >= 8 && hour <= 19 && m.isGoodHour(record) {
			goodRun++
			if goodRun == 6 {
				return true
			}
		} else {
			goodRun = 0
		}
	}

	return false
}

func (m *basicModel) isGoodHour(record noaa.HourlyRecord) bool {
	if record.Temp < minHourlyTempC {
		return false
	}

	if record.Temp > maxHourlyTempC {
		return false
	}

	if record.Prcp > maxHourlyPrcp {
		return false
	}

	return true
}

// Convert Fahrenheit to Celsius (tenths of a degree)
func f2c(f int) int {
	c := float32(f-32) / 1.8
	return int(c * 10)
}
