package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	dac "github.com/xinsnake/go-http-digest-auth-client"
)

type inverter struct {
	SerialNumber    string `json:"serialNumber"`
	LastReportDate  int    `json:"lastReportDate"`
	DevType         int    `json:"devType"`
	LastReportWatts int    `json:"lastReportWatts"`
	MaxReportWatts  int    `json:"maxReportWatts"`
}

type inverterList struct {
	inverters []inverter
}

func getEnvoyJson() ([]byte, error) {
	log.Println("Getting Envoy json from " + os.Getenv("ENVOY_URL") + "/api/v1/production/inverters")
	t := dac.NewTransport(os.Getenv("USERNAME"), os.Getenv("PASSWORD"))
	req, err := http.NewRequest("GET", os.Getenv("ENVOY_URL")+"/api/v1/production/inverters", nil)
	resp, err := t.RoundTrip(req)
	if err != nil {
		return []byte("[]"), err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	return body, nil
}

func metrics(w http.ResponseWriter, req *http.Request) {
	inverterJson, _ := getEnvoyJson()
	var inverters []inverter
	json.Unmarshal(inverterJson, &inverters)

	log.Println("Received data from", len(inverters), "inverters.")

	fmt.Fprintf(w, "# TYPE reported_wattage gauge\n")
	for _, inverter := range inverters {
		fmt.Fprintf(w, "reported_wattage{serial_number=\"%s\"} %d\n", inverter.SerialNumber, inverter.LastReportWatts)
	}
}

func main() {
	http.HandleFunc("/metrics", metrics)
	http.ListenAndServe(":80", nil)
}

func checkErr(err error) {
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}
