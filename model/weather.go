package model

import "fmt"

type WeatherRequest struct {
	Country string
	Region  string
	Ip      string
}

func (w WeatherRequest) GetNearestLocation() string {
	if len(w.Ip) > 0 {
		return w.Ip
	}

	if len(w.Region) > 0 {
		return w.Region
	}

	if len(w.Country) > 0 {
		return w.Country
	}

	return ""
}

func (w WeatherRequest) GetPreciseLocation() string {
	if len(w.Country) == 0 || len(w.Region) == 0 {
		return w.GetNearestLocation()
	}

	return fmt.Sprintf("%s %s", w.Country, w.Region)
}
