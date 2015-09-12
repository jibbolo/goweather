package main

import (
	"flag"
	"fmt"
	"os"
)

var token string
var latlon string

func init() {
	flag.StringVar(&latlon, "l", "", "Latitude,Longitude")
	flag.StringVar(&token, "t", "", "Token for Weather.com API")
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	if token == "" {
		fmt.Fprintf(os.Stderr, "Missing -token")
		os.Exit(1)
	}

	wa := newWeatherApi(token)

	loc := "autoip"
	if latlon != "" {
		loc = latlon
	}
	gl, err := wa.getGeoLookup(loc)
	checkError(err)

	hd, err := wa.getHourly(gl)
	checkError(err)

	fmt.Printf("%v (%v)\n", gl.Location.City, gl.Location.Country)
	for _, h := range hd.HourlyForecast {
		hour := h.FCTTIME.getTime().Hour()
		if hour == 0 {
			fmt.Println("========")
		}
		fmt.Printf("%2d: ", hour)
		for i := int32(0); i < (h.Pop/10)+1; i++ {
			fmt.Print("|")
		}
		fmt.Printf(" %d%%\n", h.Pop)
	}

}
