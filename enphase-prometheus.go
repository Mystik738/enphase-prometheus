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

func metrics(w http.ResponseWriter, req *http.Request) {
	inverterJson, _ := getInverterJson()
	var inverters []inverter
	json.Unmarshal(inverterJson, &inverters)

	log.Println("Received data from", len(inverters), "inverters.")

	fmt.Fprintf(w, "# TYPE reported_watts gauge\n")
	for _, inverter := range inverters {
		fmt.Fprintf(w, "reported_watts{serial_number=\"%s\"} %d\n", inverter.SerialNumber, inverter.LastReportWatts)
	}

	systemJson, err := getSystemJson()
	if err == nil {
		var system map[string]interface{}
		json.Unmarshal(systemJson, &system)

		//Some whacky conversion here, but simpler than defining the whole json object
		totalWattage := int(system["production"].([]interface{})[0].(map[string]interface{})["wNow"].(float64))

		log.Println("Received system data, current total watts is", totalWattage)
		fmt.Fprintf(w, "\n# TYPE total_watts gauge\n")
		fmt.Fprintf(w, "total_watts %d\n", totalWattage)
	} else {
		log.Println("Error retrieving system data.")
	}
}

func main() {
	http.HandleFunc("/metrics", metrics)
	http.ListenAndServe(":80", nil)
	log.Println("Server ready to serve.")
}

func checkErr(err error) {
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}
