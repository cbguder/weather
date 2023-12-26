package model

import (
	"github.com/cbguder/weather/noaa"
)

type Scorecard struct {
	Records uint
	Score   float32
}

type TrendsReport struct {
	Overall Scorecard
	Monthly map[string]Scorecard

	YearlyAvg  map[string]Scorecard
	MonthlyAvg map[string]Scorecard
}

type keyType int

const (
	keyTypeOverall keyType = iota
	keyTypeYearlyAvg
	keyTypeMonthlyAvg
	keyTypeMonthly
)

type groupKey struct {
	typ keyType
	key string
}

func Score(records []noaa.DailyRecord) Scorecard {
	m := &basicModel{}

	for _, record := range records {
		m.Add(record)
	}

	return m.Scorecard()
}

func Trends(records []noaa.DailyRecord) (*TrendsReport, error) {
	groups := make(map[groupKey]*basicModel)

	for _, record := range records {
		d := record.Date

		keys := []groupKey{
			{typ: keyTypeOverall},
			{typ: keyTypeMonthly, key: d.Format("2006/Jan")},
			{typ: keyTypeYearlyAvg, key: d.Format("2006")},
			{typ: keyTypeMonthlyAvg, key: d.Format("Jan")},
		}

		for _, key := range keys {
			m, ok := groups[key]
			if !ok {
				m = &basicModel{}
				groups[key] = m
			}
			m.Add(record)
		}
	}

	report := TrendsReport{
		Monthly:    make(map[string]Scorecard),
		YearlyAvg:  make(map[string]Scorecard),
		MonthlyAvg: make(map[string]Scorecard),
	}

	for key, group := range groups {
		scorecard := group.Scorecard()

		switch key.typ {
		case keyTypeOverall:
			report.Overall = scorecard
		case keyTypeMonthly:
			report.Monthly[key.key] = scorecard
		case keyTypeYearlyAvg:
			report.YearlyAvg[key.key] = scorecard
		case keyTypeMonthlyAvg:
			report.MonthlyAvg[key.key] = scorecard
		}
	}

	return &report, nil
}
