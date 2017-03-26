package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const baseURL = "https://avwx.rest/api/metar/"
const options = "?options=info"

var conditions = map[string]string{
	"RA":   "rain",
	"DZ":   "drizzle",
	"SN":   "snow",
	"SG":   "snow grains",
	"IC":   "ice crystals",
	"PL":   "ice pellets",
	"GR":   "hail",
	"GS":   "small hail/snow pellets",
	"UP":   "unknown precipitaton",
	"BR":   "mist",
	"FG":   "fog",
	"FU":   "smoke",
	"VA":   "volcanic ash",
	"SA":   "sand",
	"HZ":   "haze",
	"PY":   "spray",
	"DU":   "widespread dust",
	"SQ":   "squall",
	"SS":   "sandstorm",
	"DS":   "duststorm",
	"PO":   "well developed dust/sand whirls",
	"FC":   "funnel cloud",
	"VC":   "in vicinity",
	"MI":   "shallow",
	"BC":   "patches",
	"SH":   "showers",
	"PR":   "partial",
	"TS":   "thunderstorm",
	"TSRA": "thunderstorm with heavy rain",
	"BL":   "blowing",
	"DR":   "drifting",
	"FZ":   "freezing",
}

func main() {
	start := time.Now()

	if len(os.Args) < 2 {
		fmt.Println("You know I need a list of airport codes son.")
		os.Exit(1)
	}

	stations := make([]string, 0)
	metars := make(map[string]Metar)

	ch := make(chan *MetarResponse)

	for _, icao := range os.Args[1:] {
		stations = append(stations, formatICAO(icao))
	}

	done := false

	go func() {
		for !done {
			fmt.Print(".")
			time.Sleep(50 * time.Millisecond)
		}
	}()

	for _, station := range stations {
		go fetchMetar(station, ch)
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

func fetchMetar(station string, ch chan<- *MetarResponse) {
	start := time.Now()
	url := baseURL + station + options

	metarResp := new(MetarResponse)
	metarResp.ICAO = station

	resp, err := http.Get(url)
	if err != nil {
		metarResp.Error = err
		ch <- metarResp
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		metarResp.Error = fmt.Errorf("Query failed: %s", resp.Status)
		ch <- metarResp
		return
	}

	var metar Metar
	if err := json.NewDecoder(resp.Body).Decode(&metar); err != nil {
		metarResp.Error = err
		ch <- metarResp
		return
	}

	metarResp.Metar = metar
	fmt.Printf("\nFetched: %s in %.2fs\n", station, time.Since(start).Seconds())
	ch <- metarResp
}

func printMetar(metar Metar) {
	fmt.Println()
	if len(metar.Error) > 0 {
		fmt.Println("Error:", metar.Error)
		return
	}
	fmt.Printf("Station:\t%s --  %s, %s -- %s\n", metar.Station, metar.Info.City, metar.Info.State, metar.Info.Name)
	fmt.Printf("%-10s\t%s\n", "Time:", metar.Time)
	temp, _ := strconv.ParseFloat(metar.Temperature, 64)
	fmt.Printf("Temperature:\t%.1f\u00B0F\n", cToF(temp))

	dewPoint, _ := strconv.ParseFloat(metar.Dewpoint, 64)
	fmt.Printf("Dew Point:\t%.1f\u00B0F\n", cToF(dewPoint))

	fmt.Printf("%-10s\t%s\u00B0 @ %sKT", "Wind:", metar.WindDirection, metar.WindSpeed)
	if metar.WindGust != "" {
		fmt.Printf(" Gusts to %sKT", metar.WindGust)
	}
	fmt.Printf("\n")

	if len(metar.ConditionList) > 0 {
		fmt.Printf("Conditions:\t")
	}

	for i, condition := range metar.ConditionList {
		modifier := ""
		if strings.HasPrefix(condition, "-") {
			modifier = "light "
			condition = condition[1:]
		} else if strings.HasPrefix(condition, "+") {
			modifier = "heavy "
			condition = condition[1:]
		}
		fmt.Printf("%s%s", modifier, conditions[condition])
		if i < len(metar.ConditionList)-1 {
			fmt.Printf(" -- ")
		}
	}

	if len(metar.ConditionList) > 0 {
		fmt.Println()
	}

	if len(metar.CloudLayers) > 0 {
		fmt.Printf("Cloud Layers:")
	}
	for _, layer := range metar.CloudLayers {
		height, _ := strconv.ParseInt(layer[1], 10, 64)
		fmt.Printf("\t%s @ %dFT", layer[0], height*100)
	}
	if len(metar.CloudLayers) > 0 {
		fmt.Println()
	}

	fmt.Printf("Visibility:\t%ssm\n", metar.Visibility)
	pressure, _ := strconv.ParseFloat(metar.Altimeter, 64)
	fmt.Printf("Pressure:\t%.2finHg\n", pressure/100)
	fmt.Printf("Flight Rules:\t%s\n", metar.FlightRules)
	fmt.Printf("Raw Report:\t%s\n", metar.RawReport)
	fmt.Println()
}

func formatICAO(icao string) string {
	len := len(icao)

	if len < 3 || len > 4 {
		fmt.Println("Invalid airport code:", icao)
		os.Exit(1)
	}

	icao = strings.ToUpper(icao)
	if len < 4 {
		icao = "K" + icao
	}

	return icao
}

func cToF(c float64) float64 {
	return c*9/5 + 32
}

type Metar struct {
	Altimeter     string
	Dewpoint      string
	FlightRules   string `json:"Flight-Rules"`
	RawReport     string `json:"Raw-Report"`
	Remarks       string
	Station       string
	Temperature   string
	Time          string
	Visibility    string
	WindDirection string     `json:"Wind-Direction"`
	WindGust      string     `json:"Wind-Gust"`
	WindSpeed     string     `json:"Wind-Speed"`
	CloudLayers   [][]string `json:"Cloud-List"`
	ConditionList []string   `json:"Other-List"`
	Error         string
	Info          LocationInfo
}

type LocationInfo struct {
	City    string
	Country string
	Name    string
	State   string
}

type MetarResponse struct {
	Metar Metar
	Error error
	ICAO  string
}
