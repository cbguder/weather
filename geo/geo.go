package geo

import "math"

type Coordinates struct {
	Lat float64
	Lon float64
}

func Distance(p1, p2 Coordinates) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	piRad := math.Pi / 180
	la1 := p1.Lat * piRad
	lo1 := p1.Lon * piRad
	la2 := p2.Lat * piRad
	lo2 := p2.Lon * piRad

	r := 6378100.0 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	meters := 2 * r * math.Asin(math.Sqrt(h))
	kilometers := meters / 1000
	miles := kilometers * 0.621371
	return miles
}

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}
