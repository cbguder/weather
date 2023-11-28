package cmd

import (
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
	rootCmd.AddCommand(trendsCmd)
}

func trendsE(_ *cobra.Command, args []string) error {
	records, err := noaa.RecordsForStation(args[0], startDate, endDate)
	if err != nil {
		return err
	}

	groups := make(map[string][]noaa.DailyRecord)
	years := make(map[string]struct{})

	for _, record := range records {
		date, err := time.Parse("20060102", record.Date)
		if err != nil {
			return err
		}

		year := date.Format("2006")
		years[year] = struct{}{}

		keys := []string{
			year,
			date.Format("Jan"),
			date.Format("2006/Jan"),
		}

		for _, key := range keys {
			groups[key] = append(groups[key], record)
		}
	}

	sortedYears := make([]string, 0, len(years))
	for year := range years {
		sortedYears = append(sortedYears, year)
	}
	sort.Strings(sortedYears)

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

	for _, year := range sortedYears {
		row := table.Row{year}

		for _, month := range monthNames {
			key := year + "/" + month
			row = append(row, cellValue(groups[key]))
		}

		row = append(row, cellValue(groups[year]))
		t.AppendRow(row)
	}

	t.AppendSeparator()

	row := table.Row{""}
	for _, month := range monthNames {
		row = append(row, cellValue(groups[month]))
	}

	row = append(row, cellValue(records))

	t.AppendRow(row)

	t.Render()

	return nil
}

func cellValue(records []noaa.DailyRecord) any {
	if len(records) == 0 {
		return ""
	}

	scorecard := model.Score(records)

	if scorecard.Records > 0 {
		return scorecard.Score
	} else {
		return ""
	}
}
