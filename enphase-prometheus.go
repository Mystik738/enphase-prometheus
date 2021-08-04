package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
		Help: "Watts reported by individual inverters.",
	}, []string{"serial_number"})
	total_watts = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "total_watts",
		Help: "Total watts reported by the system.",
	})
)

type inverter struct {
	SerialNumber    string `json:"serialNumber"`
	LastReportDate  int    `json:"lastReportDate"`
	DevType         int    `json:"devType"`
	LastReportWatts int    `json:"lastReportWatts"`
	MaxReportWatts  int    `json:"maxReportWatts"`
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

func metrics() {
	log.Println("Initializing metrics.")
	registry = prometheus.NewRegistry()
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
	metrics()
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	http.ListenAndServe(":80", nil)
	log.Println("Server ready to serve.")
}

func checkErr(err error) {
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}
