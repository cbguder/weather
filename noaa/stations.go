package noaa

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

type Station struct {
	Id    string
	Lat   float64
	Lon   float64
	Elev  float64
	State string
	Name  string
	GSN   string
	HCN   string
	WMO   string
}

func Stations() ([]Station, error) {
	stationsFile, err := openDataFile("/pub/data/ghcn/daily/ghcnd-stations.txt")
	if err != nil {
		return nil, err
	}

	defer stationsFile.Close()

	return readStations(stationsFile)
}

func readStations(r io.Reader) ([]Station, error) {
	scanner := bufio.NewScanner(r)

	var stations []Station

	for scanner.Scan() {
		line := scanner.Text()

		lat, err := parseFloat(line[12:20])
		if err != nil {
			return nil, err
		}

		lon, err := parseFloat(line[21:30])
		if err != nil {
			return nil, err
		}

		elev, err := parseFloat(line[31:37])
		if err != nil {
			return nil, err
		}

		station := Station{
			Id:    parseString(line[0:11]),
			Lat:   lat,
			Lon:   lon,
			Elev:  elev,
			State: parseString(line[38:40]),
			Name:  parseString(line[41:71]),
			GSN:   parseString(line[72:75]),
			HCN:   parseString(line[76:79]),
			WMO:   parseString(line[80:85]),
		}

		stations = append(stations, station)
	}

	return stations, scanner.Err()
}

func parseString(s string) string {
	return strings.TrimSpace(s)
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}
