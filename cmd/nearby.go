package cmd

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"
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
	places, err := osm.Search(args[0])
	if err != nil {
		return err
	}

	if len(places) == 0 {
		return fmt.Errorf("no locations found")
	}

	loc, err := locationForPlace(places[0])
	if err != nil {
		return err
	}

	stations, err := noaa.Stations()
	if err != nil {
		return err
	}

	radius, _ := cmd.Flags().GetFloat64("radius")

	nearby := getNearbyStations(stations, loc, radius)
	sort.Slice(nearby, func(i, j int) bool {
		return nearby[i].Distance < nearby[j].Distance
	})

	stationIds := make([]string, len(nearby))
	for i, sd := range nearby {
		stationIds[i] = sd.Station.Id
	}

	err = noaa.PreloadDailyRecords(stationIds)
	if err != nil {
		return err
	}

	t := newTableWriter()
	t.SetTitle("Stations near %s", places[0].DisplayName)
	t.AppendHeader(table.Row{"ID", "Name", "Elevation", "Distance", "Records", "Score"})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 3, Transformer: sprintfTransformer("%.2fm")},
		{Number: 4, Transformer: sprintfTransformer("%.2fmi")},
		{Number: 6, Transformer: scoreTransformer},
	})

	for _, sd := range nearby {
		records, err := noaa.DailyRecords(sd.Station.Id, startDate, endDate)
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
		coords := geo.Coordinates{Lat: station.Lat, Lon: station.Lon}
		d := geo.Distance(p, coords)
		if d < radius {
			nearby = append(nearby, stationWithDistance{station, d})
		}
	}

	return nearby
}

func locationForPlace(place osm.Place) (geo.Coordinates, error) {
	lat, err := strconv.ParseFloat(place.Lat, 64)
	if err != nil {
		return geo.Coordinates{}, err
	}

	lng, err := strconv.ParseFloat(place.Lon, 64)
	if err != nil {
		return geo.Coordinates{}, err
	}

	return geo.Coordinates{
		Lat: lat,
		Lon: lng,
	}, nil
}
