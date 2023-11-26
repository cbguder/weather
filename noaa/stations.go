package noaa

import (
	"bufio"
	"io"
	"log"
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

func ReadStations(r io.Reader) ([]Station, error) {
	scanner := bufio.NewScanner(r)

	var stations []Station

	for scanner.Scan() {
		line := scanner.Text()

		station := Station{
			Id:    parseString(line[0:11]),
			Lat:   parseFloat(line[12:20]),
			Lon:   parseFloat(line[21:30]),
			Elev:  parseFloat(line[31:37]),
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

func parseFloat(s string) float64 {
	f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		log.Fatalln(err)
	}
	return f
}
