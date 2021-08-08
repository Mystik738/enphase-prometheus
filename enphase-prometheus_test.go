package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestStreamsSuccess(t *testing.T) {
	registry = prometheus.NewRegistry()
	defer initEnvoyServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/stream/meter" {
			rw.WriteHeader(http.StatusOK)
			fmt.Fprintf(rw, "data: {\"production\":{\"ph-a\":{\"p\":-0.0,\"q\":138.135,\"s\":139.586,\"v\":118.313,\"i\":1.18,\"pf\":0.0,\"f\":60.0},\"ph-b\":{\"p\":0.0,\"q\":137.861,\"s\":140.002,\"v\":118.371,\"i\":1.182,\"pf\":0.0,\"f\":60.0},\"ph-c\":{\"p\":0.0,\"q\":0.0,\"s\":0.0,\"v\":0.0,\"i\":0.0,\"pf\":0.0,\"f\":0.0}},\"net-consumption\":{\"ph-a\":{\"p\":0.0,\"q\":0.0,\"s\":17.442,\"v\":118.302,\"i\":0.147,\"pf\":0.0,\"f\":60.0},\"ph-b\":{\"p\":-0.0,\"q\":0.0,\"s\":16.803,\"v\":118.371,\"i\":0.141,\"pf\":0.0,\"f\":60.0},\"ph-c\":{\"p\":0.0,\"q\":0.0,\"s\":0.0,\"v\":0.0,\"i\":0.0,\"pf\":0.0,\"f\":0.0}},\"total-consumption\":{\"ph-a\":{\"p\":-0.0,\"q\":-138.135,\"s\":156.955,\"v\":118.307,\"i\":1.327,\"pf\":-0.0,\"f\":60.0},\"ph-b\":{\"p\":0.0,\"q\":-137.861,\"s\":123.278,\"v\":118.371,\"i\":1.041,\"pf\":0.0,\"f\":60.0},\"ph-c\":{\"p\":0.0,\"q\":0.0,\"s\":0.0,\"v\":0.0,\"i\":0.0,\"pf\":0.0,\"f\":0.0}}}\n\n")
		}
	})).Close()
	streams()
	time.Sleep(time.Duration(250) * time.Millisecond)
	handler := initPrometheus()
	req := httptest.NewRequest("GET", "http://example.com/metrics", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	checkErr(err)

	if len(data) != 864 {
		log.Println(string(data))
		t.Errorf("data should be 864 characters long, is %d", len(data))
	}
}

func TestMetricsSuccess(t *testing.T) {
	registry = prometheus.NewRegistry()
	defer initEnvoyServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/api/v1/production/inverters" {
			rw.WriteHeader(http.StatusOK)
			fmt.Fprintf(rw, "[\r\n  {\r\n    \"serialNumber\": \"482125062378\",\r\n    \"lastReportDate\": 1627599602,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 14,\r\n    \"maxReportWatts\": 240\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061710\",\r\n    \"lastReportDate\": 1627599608,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 223,\r\n    \"maxReportWatts\": 245\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062686\",\r\n    \"lastReportDate\": 1627599597,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 16,\r\n    \"maxReportWatts\": 244\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062528\",\r\n    \"lastReportDate\": 1627599606,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 219,\r\n    \"maxReportWatts\": 245\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061458\",\r\n    \"lastReportDate\": 1627599613,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 221,\r\n    \"maxReportWatts\": 245\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062610\",\r\n    \"lastReportDate\": 1627599613,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 220,\r\n    \"maxReportWatts\": 244\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062558\",\r\n    \"lastReportDate\": 1627599623,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 43,\r\n    \"maxReportWatts\": 243\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062650\",\r\n    \"lastReportDate\": 1627599617,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 213,\r\n    \"maxReportWatts\": 240\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061975\",\r\n    \"lastReportDate\": 1627599619,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 214,\r\n    \"maxReportWatts\": 242\r\n  },\r\n  {\r\n    \"serialNumber\": \"202117037990\",\r\n    \"lastReportDate\": 1627599601,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 18,\r\n    \"maxReportWatts\": 245\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062554\",\r\n    \"lastReportDate\": 1627599607,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 13,\r\n    \"maxReportWatts\": 242\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061677\",\r\n    \"lastReportDate\": 1627599604,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 16,\r\n    \"maxReportWatts\": 239\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061455\",\r\n    \"lastReportDate\": 1627599611,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 16,\r\n    \"maxReportWatts\": 241\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061240\",\r\n    \"lastReportDate\": 1627599619,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 217,\r\n    \"maxReportWatts\": 243\r\n  }\r\n]")
		}
		if req.URL.String() == "/production.json" {
			rw.WriteHeader(http.StatusOK)
			fmt.Fprintf(rw, "{\"production\":[{\"type\":\"inverters\",\"activeCount\":14,\"readingTime\":0,\"wNow\":10,\"whLifetime\":248190},{\"type\":\"eim\",\"activeCount\":0,\"measurementType\":\"production\",\"readingTime\":1627874283,\"wNow\":0.0,\"whLifetime\":0.0,\"varhLeadLifetime\":0.0,\"varhLagLifetime\":0.0,\"vahLifetime\":0.0,\"rmsCurrent\":2.365,\"rmsVoltage\":236.641,\"reactPwr\":276.755,\"apprntPwr\":279.56,\"pwrFactor\":0.0,\"whToday\":0.0,\"whLastSevenDays\":0.0,\"vahToday\":0.0,\"varhLeadToday\":0.0,\"varhLagToday\":0.0}],\"consumption\":[{\"type\":\"eim\",\"activeCount\":0,\"measurementType\":\"total-consumption\",\"readingTime\":1627874283,\"wNow\":0.0,\"whLifetime\":0.0,\"varhLeadLifetime\":0.0,\"varhLagLifetime\":0.0,\"vahLifetime\":0.0,\"rmsCurrent\":2.074,\"rmsVoltage\":236.704,\"reactPwr\":-276.755,\"apprntPwr\":490.918,\"pwrFactor\":0.0,\"whToday\":0.0,\"whLastSevenDays\":0.0,\"vahToday\":0.0,\"varhLeadToday\":0.0,\"varhLagToday\":0.0},{\"type\":\"eim\",\"activeCount\":0,\"measurementType\":\"net-consumption\",\"readingTime\":1627874283,\"wNow\":-0.0,\"whLifetime\":0.0,\"varhLeadLifetime\":0.0,\"varhLagLifetime\":0.0,\"vahLifetime\":0.0,\"rmsCurrent\":0.291,\"rmsVoltage\":236.766,\"reactPwr\":0.0,\"apprntPwr\":34.442,\"pwrFactor\":0.0,\"whToday\":0,\"whLastSevenDays\":0,\"vahToday\":0,\"varhLeadToday\":0,\"varhLagToday\":0}],\"storage\":[{\"type\":\"acb\",\"activeCount\":0,\"readingTime\":0,\"wNow\":0,\"whNow\":0,\"state\":\"idle\"}]}")
		}
	})).Close()
	metrics()
	//we need to allow the metrics to collect
	time.Sleep(time.Duration(250) * time.Millisecond)
	handler := initPrometheus()
	req := httptest.NewRequest("GET", "http://example.com/metrics", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	checkErr(err)

	if len(data) != 864 {
		log.Println(string(data))
		t.Errorf("data should be 864 characters long, is %d", len(data))
	}
}

func TestEnvoyAuthFailure(t *testing.T) {
	defer initEnvoyServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/api/v1/production/inverters" {
			rw.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(rw, "{\r\n    \"status\": 401,\r\n    \"error\": \"\",\r\n    \"info\": \"Authentication required\",\r\n    \"moreInfo\": \"\"\r\n}")
		}
	})).Close()
	os.Setenv("PASSWORD", "12345")

	data, err := getInverterJson()
	log.Printf("Returned %s", data)
	if err == nil {
		t.Errorf("expected error to not be nil")
	}
}

func TestSystemJsonFailure(t *testing.T) {
	registry = prometheus.NewRegistry()
	defer initEnvoyServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/api/v1/production/inverters" {
			rw.WriteHeader(http.StatusOK)
			fmt.Fprintf(rw, "[\r\n  {\r\n    \"serialNumber\": \"482125062378\",\r\n    \"lastReportDate\": 1627599602,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 14,\r\n    \"maxReportWatts\": 240\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061710\",\r\n    \"lastReportDate\": 1627599608,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 223,\r\n    \"maxReportWatts\": 245\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062686\",\r\n    \"lastReportDate\": 1627599597,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 16,\r\n    \"maxReportWatts\": 244\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062528\",\r\n    \"lastReportDate\": 1627599606,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 219,\r\n    \"maxReportWatts\": 245\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061458\",\r\n    \"lastReportDate\": 1627599613,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 221,\r\n    \"maxReportWatts\": 245\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062610\",\r\n    \"lastReportDate\": 1627599613,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 220,\r\n    \"maxReportWatts\": 244\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062558\",\r\n    \"lastReportDate\": 1627599623,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 43,\r\n    \"maxReportWatts\": 243\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062650\",\r\n    \"lastReportDate\": 1627599617,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 213,\r\n    \"maxReportWatts\": 240\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061975\",\r\n    \"lastReportDate\": 1627599619,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 214,\r\n    \"maxReportWatts\": 242\r\n  },\r\n  {\r\n    \"serialNumber\": \"202117037990\",\r\n    \"lastReportDate\": 1627599601,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 18,\r\n    \"maxReportWatts\": 245\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062554\",\r\n    \"lastReportDate\": 1627599607,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 13,\r\n    \"maxReportWatts\": 242\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061677\",\r\n    \"lastReportDate\": 1627599604,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 16,\r\n    \"maxReportWatts\": 239\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061455\",\r\n    \"lastReportDate\": 1627599611,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 16,\r\n    \"maxReportWatts\": 241\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061240\",\r\n    \"lastReportDate\": 1627599619,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 217,\r\n    \"maxReportWatts\": 243\r\n  }\r\n]")
		}
		if req.URL.String() == "/production.json" {
			rw.WriteHeader(http.StatusInternalServerError)
		}
		if req.URL.String() == "/stream/meter" {
			rw.WriteHeader(http.StatusOK)
			fmt.Fprintf(rw, "data: {\"production\":{\"ph-a\":{\"p\":-0.0,\"q\":138.135,\"s\":139.586,\"v\":118.313,\"i\":1.18,\"pf\":0.0,\"f\":60.0},\"ph-b\":{\"p\":0.0,\"q\":137.861,\"s\":140.002,\"v\":118.371,\"i\":1.182,\"pf\":0.0,\"f\":60.0},\"ph-c\":{\"p\":0.0,\"q\":0.0,\"s\":0.0,\"v\":0.0,\"i\":0.0,\"pf\":0.0,\"f\":0.0}},\"net-consumption\":{\"ph-a\":{\"p\":0.0,\"q\":0.0,\"s\":17.442,\"v\":118.302,\"i\":0.147,\"pf\":0.0,\"f\":60.0},\"ph-b\":{\"p\":-0.0,\"q\":0.0,\"s\":16.803,\"v\":118.371,\"i\":0.141,\"pf\":0.0,\"f\":60.0},\"ph-c\":{\"p\":0.0,\"q\":0.0,\"s\":0.0,\"v\":0.0,\"i\":0.0,\"pf\":0.0,\"f\":0.0}},\"total-consumption\":{\"ph-a\":{\"p\":-0.0,\"q\":-138.135,\"s\":156.955,\"v\":118.307,\"i\":1.327,\"pf\":-0.0,\"f\":60.0},\"ph-b\":{\"p\":0.0,\"q\":-137.861,\"s\":123.278,\"v\":118.371,\"i\":1.041,\"pf\":0.0,\"f\":60.0},\"ph-c\":{\"p\":0.0,\"q\":0.0,\"s\":0.0,\"v\":0.0,\"i\":0.0,\"pf\":0.0,\"f\":0.0}}}")
		}
	})).Close()
	metrics()
	//we need to allow the metrics to collect
	time.Sleep(time.Duration(250) * time.Millisecond)
	handler := initPrometheus()

	data, err := getSystemJson()
	log.Printf("Returned %s", data)
	if err == nil {
		t.Errorf("expected error to not be nil")
	}
	req := httptest.NewRequest("GET", "http://example.com/metrics", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err = ioutil.ReadAll(res.Body)
	checkErr(err)

	if len(data) != 863 {
		log.Println(string(data))
		t.Errorf("data should be 863 characters long, is %d", len(data))
	}
}

func initPrometheus() http.Handler {
	return promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
}

func initEnvoyServer(handler http.HandlerFunc) *httptest.Server {
	os.Setenv("USERNAME", "envoy")
	os.Setenv("PASSWORD", "123456")
	os.Setenv("SLEEP_SECONDS", "10")

	server := httptest.NewServer(handler)

	os.Setenv("ENVOY_URL", server.URL)

	return server
}
