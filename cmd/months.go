package cmd

import (
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/cbguder/weather/model"
	"github.com/cbguder/weather/noaa"
)

var monthsCmd = &cobra.Command{
	Use:   "months <station-id>",
	Short: "Show score by month for a single station",
	Args:  cobra.ExactArgs(1),
	RunE:  monthsE,
}

func init() {
	rootCmd.AddCommand(monthsCmd)
}

func monthsE(_ *cobra.Command, args []string) error {
	records, err := noaa.RecordsForStation(args[0], startDate, endDate)
	if err != nil {
		return err
	}

	groupedRecords := groupRecordsByKey(records, func(record noaa.DailyRecord) string {
		date, _ := time.Parse("20060102", record.Date)
		return date.Month().String()
	})

	t := newTableWriter()
	t.AppendHeader(table.Row{"Month", "Records", "Score"})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 3, Transformer: scoreTransformer},
	})

	for _, group := range groupedRecords {
		scorecard := model.Score(group.records)

		if scorecard.Records == 0 {
			continue
		}

		t.AppendRow(table.Row{
			group.key,
			scorecard.Records,
			scorecard.Score,
		})
	}

	t.Render()

	return nil
}
