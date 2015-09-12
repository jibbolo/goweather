package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	baseUrl     string = "http://api.wunderground.com/api/%s%%s"
	geoLookupEP string = "/geolookup/q/%s.json"
	hourlyEP    string = "/hourly%s.json"
)

type weatherApi struct {
	token string
}

func newWeatherApi(token string) *weatherApi {
	return &weatherApi{token}
}

func (wa *weatherApi) getGeoLookup(q string) (*geoLookup, error) {
	res, err := wa.getEndpoint(geoLookupEP, q)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	gl := &geoLookup{}

	enc := json.NewDecoder(res.Body)
	err = enc.Decode(gl)
	if err != nil {
		return nil, err
	}
	if gl.Location.L == "" {
		return nil, fmt.Errorf("unknown location: %s", q)
	}
	return gl, nil
}

func (wa *weatherApi) getHourly(gl *geoLookup) (*hourlyData, error) {

	res, err := wa.getEndpoint(hourlyEP, gl.Location.L)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	hd := &hourlyData{}
	err = hd.decodeFromJson(res.Body)
	if err != nil {
		return nil, err
	}
	return hd, nil
}

func (wa *weatherApi) getEndpoint(endpoint string, params ...interface{}) (*http.Response, error) {
	url := fmt.Sprintf(baseUrl, wa.token)
	if len(params) != 0 {
		endpoint = fmt.Sprintf(endpoint, params...)
	}
	url = fmt.Sprintf(url, endpoint)
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return r, nil
}

type geoLookup struct {
	Location struct {
		Country string `json:"country_iso3166"`
		City    string
		L       string
	}
}

type FCTTime struct {
	Epoch int64 `json:"epoch,string"`
}

func (e *FCTTime) getTime() time.Time {
	return time.Unix(e.Epoch, 0)
}

type hourlyData struct {
	HourlyForecast []struct {
		FCTTIME FCTTime `json="FCTTIME"`
		Pop     int32   `json:",string"`
	} `json:"hourly_forecast"`
}

func (hd *hourlyData) decodeFromJson(r io.Reader) (err error) {
	enc := json.NewDecoder(r)
	err = enc.Decode(hd)
	return
}

func readFromFile(filename string) (*hourlyData, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	hd := &hourlyData{}
	err = hd.decodeFromJson(f)
	if err != nil {
		return nil, err
	}
	return hd, nil
}
