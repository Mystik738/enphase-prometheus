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
	registry       *prometheus.Registry
	reported_watts = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "reported_watts",
		Help: "watts reported by individual inverters.",
	}, []string{"serial_number"})
	total_watts = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "total_watts",
		Help: "total watts reported by the system.",
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
	pf = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "frequency",
		Help: "frequency reported by the meter, in hertz.",
	}, []string{"type", "phase"})
	f = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "power_factor_ratio",
		Help: "the power factor ratio of the meter.",
	}, []string{"type", "phase"})
)

type inverter struct {
	SerialNumber    string `json:"serialNumber"`
	LastReportDate  int    `json:"lastReportDate"`
	DevType         int    `json:"devType"`
	LastReportWatts int    `json:"lastReportWatts"`
	MaxReportWatts  int    `json:"maxReportWatts"`
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

func getInverterJson() ([]byte, error) {
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

func getSystemJson() ([]byte, error) {
	log.Println("Getting system json from " + os.Getenv("ENVOY_URL") + "/production.json")
	resp, err := http.Get(os.Getenv("ENVOY_URL") + "/production.json")
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
					log.Println(string(line))
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
	registry.MustRegister(total_watts)
	registry.MustRegister(reported_watts)

	go func() {
		for {
			log.Println("Retrieving metrics.")
			inverterJson, _ := getInverterJson()
			var inverters []inverter
			json.Unmarshal(inverterJson, &inverters)

			log.Println("Received data from", len(inverters), "inverters.")

			for _, inverter := range inverters {
				reported_watts.With(prometheus.Labels{"serial_number": inverter.SerialNumber}).Set(float64(inverter.LastReportWatts))
			}

			systemJson, err := getSystemJson()
			if err == nil {
				var system map[string]interface{}
				json.Unmarshal(systemJson, &system)

				//Some whacky conversion here, but simpler than defining the whole json object
				totalWattage := int(system["production"].([]interface{})[0].(map[string]interface{})["wNow"].(float64))

				log.Println("Received system data, current total watts is", totalWattage)
				total_watts.Set(float64(totalWattage))
			} else {
				total_watts.Set(float64(0))
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
	registry = prometheus.NewRegistry()
	metrics()
	streams()
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	err := http.ListenAndServe(":80", nil)
	checkErr(err)
	log.Println("Server ready to serve.")
}

func checkErr(err error) {
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}
