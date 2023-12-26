package cmd

import (
	"errors"
	"sort"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github.com/cbguder/weather/model"
	"github.com/cbguder/weather/noaa"
)

var monthNames = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

var trendsCmd = &cobra.Command{
	Use:   "trends <station-id>",
	Short: "Show historic trends for a single station",
	Args:  cobra.ExactArgs(1),
	RunE:  trendsE,
}

func init() {
	trendsCmd.Flags().Bool("hourly", false, "Use hourly data (slower)")

	rootCmd.AddCommand(trendsCmd)
}

func trendsE(cmd *cobra.Command, args []string) error {
	hourly, _ := cmd.Flags().GetBool("hourly")
	stationId := args[0]

	var (
		records []noaa.DailyRecord
		station *noaa.Station
		err     error
	)

	if hourly {
		records, err = hourlyRecords(stationId)
		station = &noaa.Station{
			Id:   stationId,
			Name: stationId,
		}
	} else {
		records, station, err = dailyRecords(stationId)
	}

	if err != nil {
		return err
	}

	return renderTrends(records, station)
}

func hourlyRecords(stationId string) ([]noaa.DailyRecord, error) {
	if endDate == nil {
		now := time.Now()
		endDate = &now
	}

	if startDate == nil {
		oneYearAgo := endDate.AddDate(-1, 0, 0)
		startDate = &oneYearAgo
	}

	err := noaa.PreloadHourlyRecords(stationId, *startDate, *endDate)
	if err != nil {
		return nil, err
	}

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		return nil, err
	}

	return noaa.HourlyRecords(stationId, *startDate, *endDate, loc)
}

func dailyRecords(stationId string) ([]noaa.DailyRecord, *noaa.Station, error) {
	stations, err := noaa.Stations()
	if err != nil {
		return nil, nil, err
	}

	var station *noaa.Station
	for _, s := range stations {
		if s.Id == stationId {
			station = &s
			break
		}
	}

	if station == nil {
		return nil, nil, errors.New("station not found")
	}

	records, err := noaa.DailyRecords(stationId, startDate, endDate)
	if err != nil {
		return nil, nil, err
	}

	return records, station, nil
}

func renderTrends(records []noaa.DailyRecord, station *noaa.Station) error {
	trends, err := model.Trends(records)
	if err != nil {
		return err
	}

	var years []string
	for year := range trends.YearlyAvg {
		years = append(years, year)
	}
	sort.Strings(years)

	header := table.Row{"Year"}
	for _, month := range monthNames {
		header = append(header, month)
	}

	t := newTableWriter()
	t.SetTitle("Historic Trends for %s", station.Name)
	t.AppendHeader(header)

	columnConfigs := make([]table.ColumnConfig, 13)
	for i := 0; i < 13; i++ {
		columnConfigs[i] = table.ColumnConfig{
			Number:      i + 2,
			Transformer: scoreTransformer,
			AlignHeader: text.AlignRight,
			Align:       text.AlignRight,
		}
	}
	t.SetColumnConfigs(columnConfigs)

	for _, year := range years {
		row := table.Row{year}

		for _, month := range monthNames {
			key := year + "/" + month
			row = append(row, cellValue(trends.Monthly[key]))
		}

		row = append(row, cellValue(trends.YearlyAvg[year]))
		t.AppendRow(row)
	}

	t.AppendSeparator()

	row := table.Row{""}
	for _, month := range monthNames {
		row = append(row, cellValue(trends.MonthlyAvg[month]))
	}

	row = append(row, cellValue(trends.Overall))

	t.AppendRow(row)
	t.Render()

	return nil
}

func cellValue(scorecard model.Scorecard) any {
	if scorecard.Records == 0 {
		return ""
	}

	return scorecard.Score
}
