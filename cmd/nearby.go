package cmd

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github.com/cbguder/weather/geo"
	"github.com/cbguder/weather/model"
	"github.com/cbguder/weather/noaa"
	"github.com/cbguder/weather/osm"
)

type stationWithDistance struct {
	Station  noaa.Station
	Distance float64
}

var nearbyCmd = &cobra.Command{
	Use:   "nearby <location>",
	Short: "Show score for nearby weather stations",
	Args:  cobra.ExactArgs(1),
	RunE:  nearbyE,
}

func init() {
	nearbyCmd.Flags().Float64("radius", 10, "radius in miles")

	rootCmd.AddCommand(nearbyCmd)
}

func nearbyE(cmd *cobra.Command, args []string) error {
	loc := locationForQuery(args[0])

	stations, err := noaa.Stations()
	if err != nil {
		return err
	}

	radius, _ := cmd.Flags().GetFloat64("radius")

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
		records, err := noaa.RecordsForStation(sd.Station.Id, startDate, endDate)
		if err != nil {
			return err
		}

		scorecard := model.Score(records)

		if scorecard.Records == 0 {
			continue
		}

		t.AppendRow(table.Row{
			sd.Station.Id,
			sd.Station.Name,
			sd.Station.Elev,
			sd.Distance,
			scorecard.Records,
			scorecard.Score,
		})
	}

	t.Render()

	return nil
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
