package cmd

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/cbguder/weather/geo"
	"github.com/cbguder/weather/model"
	"github.com/cbguder/weather/noaa"
	"github.com/cbguder/weather/osm"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

const radius = 10.0

type stationWithDistance struct {
	Station  noaa.Station
	Distance float64
}

var rootCmd = &cobra.Command{
	Use:  "weather <location>",
	RunE: rootE,
	Args: cobra.ExactArgs(1),
}

func Execute() error {
	return rootCmd.Execute()
}

func rootE(_ *cobra.Command, args []string) error {
	loc := locationForQuery(args[0])

	stations, err := noaa.Stations()
	if err != nil {
		return err
	}

	nearby := getNearbyStations(stations, loc, radius)
	sort.Slice(nearby, func(i, j int) bool {
		return nearby[i].Distance < nearby[j].Distance
	})

	t := newTableWriter()
	t.AppendHeader(table.Row{"ID", "Name", "Elevation", "Distance", "Records", "Score"})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 3, Transformer: text.NewNumberTransformer("%.2fm")},
		{Number: 4, Transformer: text.NewNumberTransformer("%.2fmi")},
		{Number: 6, Transformer: text.NewNumberTransformer("%.2f")},
	})

	for _, sd := range nearby {
		records, err := noaa.RecordsForStation(sd.Station.Id)
		if err != nil {
			return err
		}

		numRecords, stationScore := model.Score(records)

		if numRecords == 0 {
			continue
		}

		t.AppendRow(table.Row{
			sd.Station.Id,
			sd.Station.Name,
			sd.Station.Elev,
			sd.Distance,
			numRecords,
			stationScore,
		})
	}

	t.Render()

	return nil
}

func newTableWriter() table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)

	t.Style().Format.Header = text.FormatDefault

	return t
}

func getNearbyStations(stations []noaa.Station, p geo.Coordinates, radius float64) []stationWithDistance {
	var nearby []stationWithDistance

	for _, station := range stations {
		d := geo.Distance(p, geo.Coordinates{station.Lat, station.Lon})
		if d < radius {
			nearby = append(nearby, stationWithDistance{station, d})
		}
	}

	return nearby
}

func locationForQuery(query string) geo.Coordinates {
	places, err := osm.Search(query)
	if err != nil {
		log.Fatalln(err)
	}

	if len(places) == 0 {
		log.Fatalln("No results found")
	}

	fmt.Println(places[0].DisplayName)

	lat, err := strconv.ParseFloat(places[0].Lat, 64)
	if err != nil {
		log.Fatalln(err)
	}

	lng, err := strconv.ParseFloat(places[0].Lon, 64)
	if err != nil {
		log.Fatalln(err)
	}

	return geo.Coordinates{lat, lng}
}
