package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/mjwood10/avwx"
)

func main() {
	start := time.Now()

	if len(os.Args) < 2 {
		fmt.Println("You know I need a list of airport codes son.")
		os.Exit(1)
	}

	stations := make([]string, 0)
	metars := make(map[string]avwx.Metar)

	ch := make(chan *avwx.MetarResponse)

	for _, icao := range os.Args[1:] {
		stations = append(stations, avwx.FormatICAO(icao))
	}

	done := false

	go func() {
		for !done {
			fmt.Print(".")
			time.Sleep(50 * time.Millisecond)
		}
	}()

	for _, station := range stations {
		go func(station string, ch chan *avwx.MetarResponse) {
			ch <- avwx.FetchMetar(station)
		}(station, ch)
	}

	for range stations {
		resp := <-ch
		if resp.Error != nil {
			fmt.Fprintf(os.Stderr, "\nError reading metar for station: %s: %v\n", resp.ICAO, resp.Error)
		} else {
			metars[resp.ICAO] = resp.Metar
		}
	}

	done = true

	fmt.Printf("All stations fetched in %.2fs\n", time.Since(start).Seconds())

	for _, station := range stations {
		if metar, ok := metars[station]; ok {
			printMetar(metar)
		}
	}
}

func printMetar(metar avwx.Metar) {
	fmt.Println()
	if len(metar.Error) > 0 {
		fmt.Println("Error:", metar.Error)
		return
	}
	fmt.Printf("Station:\t%s --  %s, %s -- %s\n", metar.Station, metar.LocationInfo.City, metar.LocationInfo.State, metar.LocationInfo.Name)
	fmt.Printf("%-10s\t%s\n", "Time:", metar.Time)
	fmt.Printf("Temperature:\t%s\u00B0F / %s\u00B0C\n", metar.TemperatureF, metar.Temperature)
	fmt.Printf("Dew Point:\t%s\u00B0F / %s\u00B0C\n", metar.Dewpoint, metar.DewpointF)

	windSpeed, _ := strconv.ParseInt(metar.WindSpeed, 10, 32)
	fmt.Printf("%-10s\t%s\u00B0 (%s) @%dKT", "Wind:", metar.WindDirection, metar.WindDirectionDesc, windSpeed)
	if metar.WindGust != "" {
		fmt.Printf(" Gusts to %sKT", metar.WindGust)
	}
	fmt.Printf("\n")

	if len(metar.ConditionsDec) > 0 {
		fmt.Printf("Conditions:\t")
	}

	for i, condition := range metar.ConditionsDec {
		fmt.Printf("%s %s", condition.Modifier, condition.Desc)
		if len(condition.Other) > 0 {
			fmt.Printf("%s", condition.Other)
		}
		if i < len(metar.ConditionsDec)-1 {
			fmt.Printf(" -- ")
		}
	}

	if len(metar.ConditionsDec) > 0 {
		fmt.Println()
	}

	if len(metar.CloudLayersDec) > 0 {
		fmt.Printf("Cloud Layers:\t")
	}
	for i, layer := range metar.CloudLayersDec {
		fmt.Printf("%s @%sFT", layer.Coverage, layer.HeightFt)
		if len(layer.Type) > 0 {
			fmt.Printf(" (%s)", layer.Type)
		}
		if i < len(metar.CloudLayers)-1 {
			fmt.Printf(" -- ")
		}
	}
	if len(metar.CloudLayers) > 0 {
		fmt.Println()
	}

	fmt.Printf("Visibility:\t%ssm\n", metar.Visibility)
	fmt.Printf("Pressure:\t%sinHg\n", metar.Altimeter)
	fmt.Printf("Flight Rules:\t%s\n", metar.FlightRules)
	fmt.Printf("Raw Report:\t%s\n", metar.RawReport)
	fmt.Println()
}
