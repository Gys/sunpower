package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	date2 := time.Now().AddDate(0, 0, -1) // yesterday
	date1 := date2.AddDate(-10, 0, 0)     // 10 years before yesterday

	lat := "38.722501"
	lon := "-9.4323331"
	if len(os.Args) == 3 {
		lat = os.Args[1]
		lon = os.Args[2]
	} else {
		fmt.Printf("\nUsage: ./sunpower <lat> <lon>\n\n")
	}
	fmt.Printf("Request data for period %s until %s for latlon %s, %s...\n", date1.Format("Jan 2, 2006"), date2.Format("Jan 2, 2006"), lat, lon)

	// get the insolation from NASA
	// docs: https://power.larc.nasa.gov/docs/v1/
	// ALLSKY_SFC_SW_DWN = All Sky Insolation Incident on a Horizontal Surface
	// result is kW-hr/m^2/day
	resp, err := http.Get("https://power.larc.nasa.gov/cgi-bin/v1/DataAccess.py?request=execute&identifier=SinglePoint&parameters=ALLSKY_SFC_SW_DWN&startDate=" + date1.Format("20060102") + "&endDate=" + date2.Format("20060102") + "&userCommunity=SSE&tempAverage=DAILY&outputList=CSV&lat=" + lat + "&lon=" + lon + "&user=anonymous")
	if err != nil {
		println("could not get an api result", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println("could not read api result", err)
		return
	}
	pwr := powerResult{}
	err = json.Unmarshal(body, &pwr)
	if err != nil {
		fmt.Printf("unmarshal error: %s", err)
		return
	}
	fmt.Printf("Downloading data as CSV from %s...\n", pwr.Outputs.Csv)

	// get the csv data and calc the average value
	resp, err = http.Get(pwr.Outputs.Csv)
	if err != nil {
		println("could not download the csv data", err)
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		println("cound not read csv data", err)
		return
	}
	csv := string(body)
	start := len(csv)
	count := 0
	sum := 0.0
	for i, r := range strings.Split(csv, "\n") {
		if strings.HasPrefix(r, "LAT,LON,") {
			start = i
		}
		if i > start {
			// example line: 38.72251,-9.43232,2018,03,21,              5.91
			fields := strings.Split(r, ",")
			if len(fields) == 6 {
				s := strings.TrimSpace(fields[5])
				if s != "-999" { // ignore days without valid values
					val, err := strconv.ParseFloat(s, 64)
					if err == nil {
						count++
						sum += val
					}
				}
			}
		}
	}
	fmt.Printf("Calculated average of %0.2f kiloWatt Hour per square meter per day\n", sum/float64(count))
}

type powerResult struct {
	Features []struct {
		Geometry struct {
			Coordinates []float64 `json:"coordinates"`
			Type        string    `json:"type"`
		} `json:"geometry"`
		Properties struct {
			Parameter struct {
				// skip the data
			} `json:"parameter"`
		} `json:"properties"`
		Type string `json:"type"`
	} `json:"features"`
	Header struct {
		APIVersion string `json:"api_version"`
		EndDate    string `json:"endDate"`
		FillValue  string `json:"fillValue"`
		StartDate  string `json:"startDate"`
		Title      string `json:"title"`
	} `json:"header"`
	Messages []interface{} `json:"messages"`
	Outputs  struct {
		Csv string `json:"csv"`
	} `json:"outputs"`
	ParameterInformation struct {
		ALLSKYSFCSWDWN struct {
			Longname string `json:"longname"`
			Units    string `json:"units"`
		} `json:"ALLSKY_SFC_SW_DWN"`
	} `json:"parameterInformation"`
	Time [][]interface{} `json:"time"`
	Type string          `json:"type"`
}
