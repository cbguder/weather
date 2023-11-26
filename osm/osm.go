package osm

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type Place struct {
	PlaceId     int      `json:"place_id"`
	Licence     string   `json:"licence"`
	OsmType     string   `json:"osm_type"`
	OsmId       int      `json:"osm_id"`
	Lat         string   `json:"lat"`
	Lon         string   `json:"lon"`
	Category    string   `json:"category"`
	Type        string   `json:"type"`
	PlaceRank   int      `json:"place_rank"`
	Importance  float64  `json:"importance"`
	Addresstype string   `json:"addresstype"`
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Boundingbox []string `json:"boundingbox"`
}

func Search(query string) ([]Place, error) {
	q := url.Values{}
	q.Set("q", query)
	q.Set("format", "jsonv2")
	q.Set("countrycodes", "us")

	loc := "https://nominatim.openstreetmap.org/search?" + q.Encode()

	resp, err := http.Get(loc)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var places []Place
	err = json.NewDecoder(resp.Body).Decode(&places)
	return places, err
}
