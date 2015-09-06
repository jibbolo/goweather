package main

import (
	"flag"
	"fmt"
	"os"
)

var token string

func init() {
	flag.StringVar(&token, "token", "", "Token for Weather.com API")
	flag.Parse()
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	wa := newWeatherApi(token)

	gl, err := wa.getGeoLookup()
	checkError(err)

	hd, err := wa.getHourly(gl)
	checkError(err)

	fmt.Printf("%v (%v)\n", gl.Location.City, gl.Location.Country)
	for _, h := range hd.HourlyForecast {
		fmt.Printf("%2d: ", h.FCTTIME.getTime().Hour())
		for i := int32(0); i < (h.Pop/10)+1; i++ {
			fmt.Print("|")
		}
		fmt.Printf(" %d%%\n", h.Pop)
	}

}
