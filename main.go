package main

import (
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/cbguder/weather/geo"
	"github.com/cbguder/weather/model"
	"github.com/cbguder/weather/noaa"
	"github.com/cbguder/weather/osm"
)

const (
	radius = 10.0
)

type stationWithDistance struct {
	Station  noaa.Station
	Distance float64
}

func main() {
	stations := readAllStations()

	if len(os.Args) < 2 {
		log.Fatalln("Please provide a search query")
	}

	loc := locationForQuery(os.Args[1])
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
		records := readRecordsForStation(sd.Station.Id)
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

func readAllStations() []noaa.Station {
	stationsFile, err := noaa.OpenDataFile("ghcnd-stations.txt")
	if err != nil {
		log.Fatalln(err)
	}

	defer stationsFile.Close()

	stations, err := noaa.ReadStations(stationsFile)
	if err != nil {
		log.Fatalln(err)
	}

	return stations
}

func readRecordsForStation(stationId string) []noaa.DailyRecord {
	path := fmt.Sprintf("by_station/%s.csv.gz", stationId)

	recordsFile, err := noaa.OpenDataFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	defer recordsFile.Close()

	gz, err := gzip.NewReader(recordsFile)
	if err != nil {
		log.Fatalln(err)
	}

	defer gz.Close()

	records, err := noaa.ReadDailyRecords(gz)
	if err != nil {
		log.Fatalln(err)
	}

	return records
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
