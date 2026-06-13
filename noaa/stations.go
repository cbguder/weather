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

		lat, err := parseFloat(field(line, 12, 20))
		if err != nil {
			return nil, err
		}

		lon, err := parseFloat(field(line, 21, 30))
		if err != nil {
			return nil, err
		}

		elev, err := parseFloat(field(line, 31, 37))
		if err != nil {
			return nil, err
		}

		station := Station{
			Id:    parseString(field(line, 0, 11)),
			Lat:   lat,
			Lon:   lon,
			Elev:  elev,
			State: parseString(field(line, 38, 40)),
			Name:  parseString(field(line, 41, 71)),
			GSN:   parseString(field(line, 72, 75)),
			HCN:   parseString(field(line, 76, 79)),
			WMO:   parseString(field(line, 80, 85)),
		}

		stations = append(stations, station)
	}

	return stations, scanner.Err()
}

// field returns line[start:end], clamped to the line's length so that
// records shorter than the full fixed-width format don't cause a panic.
func field(line string, start, end int) string {
	if start > len(line) {
		return ""
	}
	if end > len(line) {
		end = len(line)
	}
	return line[start:end]
}

func parseString(s string) string {
	return strings.TrimSpace(s)
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}
