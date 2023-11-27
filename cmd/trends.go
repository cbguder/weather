package cmd

import (
	"strconv"
	"strings"
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

	groupedRecords := groupRecordsByKey(records, func(record noaa.DailyRecord) string {
		date, _ := time.Parse("20060102", record.Date)
		return date.Format("2006/01")
	})

	var columnConfigs []table.ColumnConfig
	header := table.Row{"Year"}
	for i, month := range monthNames {
		header = append(header, month)

		columnConfigs = append(columnConfigs, table.ColumnConfig{
			Number:      i + 2,
			Transformer: scoreTransformer,
			AlignHeader: text.AlignRight,
		})
	}

	t := newTableWriter()
	t.AppendHeader(header)
	t.SetColumnConfigs(columnConfigs)

	row := table.Row{""}
	for _, group := range groupedRecords {
		parts := strings.Split(group.key, "/")
		year := parts[0]
		month, _ := strconv.Atoi(parts[1])

		if row[0] != year {
			if len(row) > 1 {
				t.AppendRow(row)
			}

			row = table.Row{year, "", "", "", "", "", "", "", "", "", "", "", ""}
		}

		scorecard := model.Score(group.records)
		if scorecard.Records > 0 {
			row[month] = scorecard.Score
		}
	}

	t.AppendRow(row)
	t.Render()

	return nil
}
