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

var cloudTypes = map[string]string{
	"CB":    "cumulonimbus",
	"TCU":   "towering cumulus",
	"CBMAM": "cumulonimbus mammatus",
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
	temp, _ := strconv.ParseFloat(metar.Temperature, 32)
	fmt.Printf("Temperature:\t%.1f\u00B0F / %.1f\u00B0C\n", cToF(temp), temp)

	dewPoint, _ := strconv.ParseFloat(metar.Dewpoint, 64)
	fmt.Printf("Dew Point:\t%.1f\u00B0F / %.1f\u00B0C\n", cToF(dewPoint), dewPoint)

	windDegrees, _ := strconv.ParseInt(metar.WindDirection, 10, 32)
	windSpeed, _ := strconv.ParseInt(metar.WindSpeed, 10, 32)
	fmt.Printf("%-10s\t%s\u00B0 (%s) @%dKT", "Wind:", metar.WindDirection, getDirection(windDegrees), windSpeed)
	if metar.WindGust != "" {
		fmt.Printf(" Gusts to %sKT", metar.WindGust)
	}
	fmt.Printf("\n")

	if len(metar.ConditionList) > 0 {
		fmt.Printf("Conditions:\t")
	}

	for i, condition := range metar.ConditionList {
		modifier := ""
		vicinity := false

		if strings.HasPrefix(condition, "VC") {
			vicinity = true
			condition = condition[2:]
		}
		if strings.HasPrefix(condition, "-") {
			modifier = "light "
			condition = condition[1:]
		} else if strings.HasPrefix(condition, "+") {
			modifier = "heavy "
			condition = condition[1:]
		}

		fmt.Printf("%s%s", modifier, conditions[condition])
		if vicinity {
			fmt.Printf(" in vicinity")
		}
		if i < len(metar.ConditionList)-1 {
			fmt.Printf(" -- ")
		}
	}

	if len(metar.ConditionList) > 0 {
		fmt.Println()
	}

	if len(metar.CloudLayers) > 0 {
		fmt.Printf("Cloud Layers:\t")
	}
	for i, layer := range metar.CloudLayers {
		height, _ := strconv.ParseInt(layer[1], 10, 32)
		fmt.Printf("%s @%dFT", layer[0], height*100)
		if len(layer) > 2 {
			fmt.Printf(" (%s)", cloudTypes[layer[2]])
		}
		if i < len(metar.CloudLayers)-1 {
			fmt.Printf(" -- ")
		}
	}
	if len(metar.CloudLayers) > 0 {
		fmt.Println()
	}

	fmt.Printf("Visibility:\t%ssm\n", metar.Visibility)
	pressure, _ := strconv.ParseFloat(metar.Altimeter, 32)
	fmt.Printf("Pressure:\t%.2finHg\n", pressure/100)
	fmt.Printf("Flight Rules:\t%s\n", metar.FlightRules)
	fmt.Printf("Raw Report:\t%s\n", metar.RawReport)
	fmt.Println()
}

func getDirection(degrees int64) string {
	switch {
	case (degrees > 349 && degrees <= 360) || (degrees >= 0 && degrees <= 11):
		return "N"
	case degrees > 11 && degrees <= 34:
		return "NNE"
	case degrees > 34 && degrees <= 56:
		return "NE"
	case degrees > 56 && degrees <= 79:
		return "ENE"
	case degrees > 79 && degrees <= 101:
		return "E"
	case degrees > 101 && degrees <= 124:
		return "ESE"
	case degrees > 124 && degrees <= 146:
		return "SE"
	case degrees > 146 && degrees <= 169:
		return "SSE"
	case degrees > 169 && degrees <= 191:
		return "S"
	case degrees > 191 && degrees <= 214:
		return "SSW"
	case degrees > 214 && degrees <= 236:
		return "SW"
	case degrees > 236 && degrees <= 259:
		return "WSW"
	case degrees > 259 && degrees <= 281:
		return "W"
	case degrees > 281 && degrees <= 304:
		return "WNW"
	case degrees > 304 && degrees <= 326:
		return "NW"
	case degrees > 326 && degrees <= 349:
		return "NNW"
	default:
		return ""
	}
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
