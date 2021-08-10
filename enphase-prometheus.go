package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dac "github.com/xinsnake/go-http-digest-auth-client"
)

var (
	arrayLocations map[string]geo

	registry      *prometheus.Registry
	reportedWatts = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "reported_watts",
		Help: "watts reported by individual inverters.",
	}, []string{"serial_number", "x", "y"})
	totalWatts = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "total_watts",
		Help: "total watts reported by the system.",
	})
	wattHoursToday = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "watt_hours_today",
		Help: "total watt hours today",
	})
	wattHoursSevenDays = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "watt_hours_seven_days",
		Help: "total watt hours past seven days",
	})
	wattHoursLifetime = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "watt_hours_lifetime",
		Help: "watt hours produced over the lifetime of this device",
	})
	p = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "active_power",
		Help: "active power reported by the meter, in watts.",
	}, []string{"type", "phase"})
	q = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "reactive_power",
		Help: "reactive power reported by the meter, in watts.",
	}, []string{"type", "phase"})
	s = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "apparent_power",
		Help: "apparent power reported by the meter, in watts.",
	}, []string{"type", "phase"})
	v = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "voltage",
		Help: "voltage reported by the meter, in volts.",
	}, []string{"type", "phase"})
	i = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "amperage",
		Help: "current reported by the meter, in amperes.",
	}, []string{"type", "phase"})
	f = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "frequency",
		Help: "frequency reported by the meter, in hertz.",
	}, []string{"type", "phase"})
	pf = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "power_factor_ratio",
		Help: "the power factor ratio of the meter.",
	}, []string{"type", "phase"})
)

type production struct {
	WattHoursToday     int `json:"wattHoursToday"`
	WattHoursSevenDays int `json:"wattHoursSevenDays"`
	WattHoursLifetime  int `json:"wattHoursLifetime"`
	WattsNow           int `json:"wattsNow"`
}

type inverter struct {
	SerialNumber    string `json:"serialNumber"`
	LastReportDate  int    `json:"lastReportDate"`
	DevType         int    `json:"devType"`
	LastReportWatts int    `json:"lastReportWatts"`
	MaxReportWatts  int    `json:"maxReportWatts"`
}

//Generated with https://mholt.github.io/json-to-go/
type arrayLayout struct {
	SystemID   int `json:"system_id"`
	Rotation   int `json:"rotation"`
	Dimensions struct {
		XMin int `json:"x_min"`
		XMax int `json:"x_max"`
		YMin int `json:"y_min"`
		YMax int `json:"y_max"`
	} `json:"dimensions"`
	Arrays []struct {
		ArrayID int    `json:"array_id"`
		Label   string `json:"label"`
		X       int    `json:"x"`
		Y       int    `json:"y"`
		Azimuth int    `json:"azimuth"`
		Modules []struct {
			ModuleID int `json:"module_id"`
			Rotation int `json:"rotation"`
			X        int `json:"x"`
			Y        int `json:"y"`
			Inverter struct {
				InverterID int    `json:"inverter_id"`
				SerialNum  string `json:"serial_num"`
			} `json:"inverter"`
		} `json:"modules"`
		Dimensions struct {
			XMin int `json:"x_min"`
			XMax int `json:"x_max"`
			YMin int `json:"y_min"`
			YMax int `json:"y_max"`
		} `json:"dimensions"`
		Tilt            int    `json:"tilt"`
		ArrayTypeName   string `json:"array_type_name"`
		PcuCount        int    `json:"pcu_count"`
		PvModuleDetails struct {
			Manufacturer string      `json:"manufacturer"`
			Model        string      `json:"model"`
			Type         interface{} `json:"type"`
			PowerRating  interface{} `json:"power_rating"`
		} `json:"pv_module_details"`
	} `json:"arrays"`
	Haiku string `json:"haiku"`
}

type geo struct {
	X int
	Y int
}

type phase struct {
	P  float64 `json:"p"`
	Q  float64 `json:"q"`
	S  float64 `json:"s"`
	V  float64 `json:"v"`
	I  float64 `json:"i"`
	Pf float64 `json:"pf"`
	F  float64 `json:"f"`
}

type threePhase struct {
	PhA phase `json:"ph-a"`
	PhB phase `json:"ph-b"`
	PhC phase `json:"ph-c"`
}

func getInverterJSON() ([]byte, error) {
	log.Println("Getting system json from " + os.Getenv("ENVOY_URL") + "/api/v1/production/inverters")
	t := dac.NewTransport(os.Getenv("USERNAME"), os.Getenv("PASSWORD"))
	req, err := http.NewRequest("GET", os.Getenv("ENVOY_URL")+"/api/v1/production/inverters", nil)
	checkErr(err)
	resp, err := t.RoundTrip(req)
	if err != nil {
		return []byte("[]"), err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	return body, nil
}

func getSystemJSON() ([]byte, error) {
	log.Println("Getting system json from " + os.Getenv("ENVOY_URL") + "/api/v1/production")
	resp, err := http.Get(os.Getenv("ENVOY_URL") + "/api/v1/production")
	checkErr(err)
	if resp.StatusCode != http.StatusOK {
		return []byte("[]"), fmt.Errorf("received http status %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	return body, nil
}

func streams() {
	log.Println("Initializing stream.")
	registry.MustRegister(p)
	registry.MustRegister(q)
	registry.MustRegister(s)
	registry.MustRegister(v)
	registry.MustRegister(i)
	registry.MustRegister(pf)
	registry.MustRegister(f)

	var gauges []*prometheus.GaugeVec = []*prometheus.GaugeVec{p, q, s, v, i, pf, f}

	go func() {
		t := dac.NewTransport(os.Getenv("USERNAME"), os.Getenv("PASSWORD"))
		t.HTTPClient = &http.Client{
			Timeout: time.Second * 3600,
		}
		retries := 1
		for {
			log.Println("Reading from stream.")
			req, err := http.NewRequest("GET", os.Getenv("ENVOY_URL")+"/stream/meter", nil)
			checkErr(err)
			resp, err := t.RoundTrip(req)
			if err == nil {
				retries = 1
				reader := bufio.NewReader(resp.Body)
				var stream map[string]threePhase
				line, err := reader.ReadBytes('\n')
				for err == nil {
					if len(line) > 2 {
						line = line[6:]
						json.Unmarshal(line, &stream)

						for phaseType := range stream {
							for i, gauge := range gauges {
								vA := reflect.ValueOf(stream[phaseType].PhA)
								(*gauge).With(prometheus.Labels{"type": phaseType, "phase": "ph-a"}).Set(vA.Field(i).Interface().(float64))

								vB := reflect.ValueOf(stream[phaseType].PhB)
								(*gauge).With(prometheus.Labels{"type": phaseType, "phase": "ph-b"}).Set(vB.Field(i).Interface().(float64))

								vC := reflect.ValueOf(stream[phaseType].PhC)
								(*gauge).With(prometheus.Labels{"type": phaseType, "phase": "ph-c"}).Set(vC.Field(i).Interface().(float64))
							}
						}
					}

					line, err = reader.ReadBytes('\n')
				}
			} else {
				log.Println("Error reading from stream.")
				log.Println(err.Error())
				retries *= 2
				time.Sleep(time.Duration(retries) * 100 * time.Millisecond)
				if retries > 300 {
					retries = 300
				}
			}
		}
	}()
}

func metrics() {
	log.Println("Initializing metrics.")

	if os.Getenv("ARRAY_LAYOUT") != "" {
		log.Println("Reading layout information.")
		arrayLocations = make(map[string]geo)
		var arrayDefinition arrayLayout
		json.Unmarshal([]byte(os.Getenv("ARRAY_LAYOUT")), &arrayDefinition)

		for _, solarArray := range arrayDefinition.Arrays {
			for _, module := range solarArray.Modules {
				arrayLocations[module.Inverter.SerialNum] = geo{
					X: module.X,
					Y: module.Y,
				}
			}
		}
	}
	registry.MustRegister(reportedWatts)
	registry.MustRegister(totalWatts)
	registry.MustRegister(wattHoursToday)
	registry.MustRegister(wattHoursSevenDays)
	registry.MustRegister(wattHoursLifetime)

	go func() {
		for {
			log.Println("Retrieving metrics.")
			inverterJSON, _ := getInverterJSON()
			var inverters []inverter
			json.Unmarshal(inverterJSON, &inverters)

			log.Println("Received data from", len(inverters), "inverters.")

			for _, inverter := range inverters {
				if val, ok := arrayLocations[inverter.SerialNumber]; ok {
					reportedWatts.With(prometheus.Labels{"serial_number": inverter.SerialNumber, "x": strconv.Itoa(val.X), "y": strconv.Itoa(val.Y)}).Set(float64(inverter.LastReportWatts))
				} else {
					reportedWatts.With(prometheus.Labels{"serial_number": inverter.SerialNumber, "x": "0", "y": "0"}).Set(float64(inverter.LastReportWatts))
				}
			}

			systemJSON, err := getSystemJSON()
			if err == nil {
				var system production
				json.Unmarshal(systemJSON, &system)

				log.Println("Received system data, current total watts is", system.WattsNow)
				totalWatts.Set(float64(system.WattsNow))
				wattHoursSevenDays.Set(float64(system.WattHoursSevenDays))
				wattHoursLifetime.Set(float64(system.WattHoursLifetime))
				wattHoursToday.Set(float64(system.WattHoursToday))
			} else {
				totalWatts.Set(float64(0))
				log.Println("Error retrieving system data.")
			}

			sleep := 10
			if os.Getenv("SLEEP_SECONDS") != "" {
				sleep, err = strconv.Atoi(os.Getenv("SLEEP_SECONDS"))
				checkErr(err)
			}
			time.Sleep(time.Duration(sleep) * time.Second)
		}
	}()
}

func main() {
	http.Handle("/metrics", initPrometheus())
	metrics()
	streams()
	port := "80"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	err := http.ListenAndServe(":"+port, nil)
	checkErr(err)
	log.Println("Server ready to serve.")
}

func checkErr(err error) {
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}

func initPrometheus() http.Handler {
	registry = prometheus.NewRegistry()
	return promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
}
