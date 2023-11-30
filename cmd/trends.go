package cmd

import (
	"sort"

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
	rootCmd.AddCommand(trendsCmd)
}

func trendsE(_ *cobra.Command, args []string) error {
	records, err := noaa.RecordsForStation(args[0], startDate, endDate)
	if err != nil {
		return err
	}

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
