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
)

func TestMetricsSuccess(t *testing.T) {
	handler := initPrometheus()
	defer initEnvoyServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/api/v1/production/inverters" {
			rw.WriteHeader(http.StatusOK)
			fmt.Fprintf(rw, "[\r\n  {\r\n    \"serialNumber\": \"482125062378\",\r\n    \"lastReportDate\": 1627599602,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 14,\r\n    \"maxReportWatts\": 240\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061710\",\r\n    \"lastReportDate\": 1627599608,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 223,\r\n    \"maxReportWatts\": 245\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062686\",\r\n    \"lastReportDate\": 1627599597,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 16,\r\n    \"maxReportWatts\": 244\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062528\",\r\n    \"lastReportDate\": 1627599606,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 219,\r\n    \"maxReportWatts\": 245\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061458\",\r\n    \"lastReportDate\": 1627599613,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 221,\r\n    \"maxReportWatts\": 245\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062610\",\r\n    \"lastReportDate\": 1627599613,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 220,\r\n    \"maxReportWatts\": 244\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062558\",\r\n    \"lastReportDate\": 1627599623,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 43,\r\n    \"maxReportWatts\": 243\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062650\",\r\n    \"lastReportDate\": 1627599617,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 213,\r\n    \"maxReportWatts\": 240\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061975\",\r\n    \"lastReportDate\": 1627599619,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 214,\r\n    \"maxReportWatts\": 242\r\n  },\r\n  {\r\n    \"serialNumber\": \"202117037990\",\r\n    \"lastReportDate\": 1627599601,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 18,\r\n    \"maxReportWatts\": 245\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062554\",\r\n    \"lastReportDate\": 1627599607,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 13,\r\n    \"maxReportWatts\": 242\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061677\",\r\n    \"lastReportDate\": 1627599604,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 16,\r\n    \"maxReportWatts\": 239\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061455\",\r\n    \"lastReportDate\": 1627599611,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 16,\r\n    \"maxReportWatts\": 241\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061240\",\r\n    \"lastReportDate\": 1627599619,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 217,\r\n    \"maxReportWatts\": 243\r\n  }\r\n]")
		}
		if req.URL.String() == "/production.json" {
			rw.WriteHeader(http.StatusOK)
			fmt.Fprintf(rw, "{\"production\":[{\"type\":\"inverters\",\"activeCount\":14,\"readingTime\":0,\"wNow\":10,\"whLifetime\":248190},{\"type\":\"eim\",\"activeCount\":0,\"measurementType\":\"production\",\"readingTime\":1627874283,\"wNow\":0.0,\"whLifetime\":0.0,\"varhLeadLifetime\":0.0,\"varhLagLifetime\":0.0,\"vahLifetime\":0.0,\"rmsCurrent\":2.365,\"rmsVoltage\":236.641,\"reactPwr\":276.755,\"apprntPwr\":279.56,\"pwrFactor\":0.0,\"whToday\":0.0,\"whLastSevenDays\":0.0,\"vahToday\":0.0,\"varhLeadToday\":0.0,\"varhLagToday\":0.0}],\"consumption\":[{\"type\":\"eim\",\"activeCount\":0,\"measurementType\":\"total-consumption\",\"readingTime\":1627874283,\"wNow\":0.0,\"whLifetime\":0.0,\"varhLeadLifetime\":0.0,\"varhLagLifetime\":0.0,\"vahLifetime\":0.0,\"rmsCurrent\":2.074,\"rmsVoltage\":236.704,\"reactPwr\":-276.755,\"apprntPwr\":490.918,\"pwrFactor\":0.0,\"whToday\":0.0,\"whLastSevenDays\":0.0,\"vahToday\":0.0,\"varhLeadToday\":0.0,\"varhLagToday\":0.0},{\"type\":\"eim\",\"activeCount\":0,\"measurementType\":\"net-consumption\",\"readingTime\":1627874283,\"wNow\":-0.0,\"whLifetime\":0.0,\"varhLeadLifetime\":0.0,\"varhLagLifetime\":0.0,\"vahLifetime\":0.0,\"rmsCurrent\":0.291,\"rmsVoltage\":236.766,\"reactPwr\":0.0,\"apprntPwr\":34.442,\"pwrFactor\":0.0,\"whToday\":0,\"whLastSevenDays\":0,\"vahToday\":0,\"varhLeadToday\":0,\"varhLagToday\":0}],\"storage\":[{\"type\":\"acb\",\"activeCount\":0,\"readingTime\":0,\"wNow\":0,\"whNow\":0,\"state\":\"idle\"}]}")
		}
		if req.URL.String() == "/api/v1/production" {
			rw.WriteHeader(http.StatusOK)
			fmt.Fprintf(rw, "{\r\n  \"wattHoursToday\": 7883,\r\n  \"wattHoursSevenDays\": 140276,\r\n  \"wattHoursLifetime\": 396713,\r\n  \"wattsNow\": 10\r\n}")
		}
	})).Close()
	metrics()
	//we need to allow the metrics to collect
	time.Sleep(time.Duration(250) * time.Millisecond)

	getMetrics(t, handler, 1462)
}

func TestEnvoyAuthFailure(t *testing.T) {
	defer initEnvoyServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/api/v1/production/inverters" {
			rw.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(rw, "{\r\n    \"status\": 401,\r\n    \"error\": \"\",\r\n    \"info\": \"Authentication required\",\r\n    \"moreInfo\": \"\"\r\n}")
		}
	})).Close()
	os.Setenv("AUTH_TOKEN", "1kd4js")

	data, _ := getInverterJSON()
	log.Printf("Returned %s", data)
	if string(data) != "{\r\n    \"status\": 401,\r\n    \"error\": \"\",\r\n    \"info\": \"Authentication required\",\r\n    \"moreInfo\": \"\"\r\n}" {
		t.Errorf("expected error to not be nil")
	}
}

func TestSystemJsonFailure(t *testing.T) {
	handler := initPrometheus()
	defer initEnvoyServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/api/v1/production/inverters" {
			rw.WriteHeader(http.StatusOK)
			fmt.Fprintf(rw, "[\r\n  {\r\n    \"serialNumber\": \"482125062378\",\r\n    \"lastReportDate\": 1627599602,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 14,\r\n    \"maxReportWatts\": 240\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061710\",\r\n    \"lastReportDate\": 1627599608,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 223,\r\n    \"maxReportWatts\": 245\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062686\",\r\n    \"lastReportDate\": 1627599597,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 16,\r\n    \"maxReportWatts\": 244\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062528\",\r\n    \"lastReportDate\": 1627599606,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 219,\r\n    \"maxReportWatts\": 245\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061458\",\r\n    \"lastReportDate\": 1627599613,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 221,\r\n    \"maxReportWatts\": 245\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062610\",\r\n    \"lastReportDate\": 1627599613,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 220,\r\n    \"maxReportWatts\": 244\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062558\",\r\n    \"lastReportDate\": 1627599623,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 43,\r\n    \"maxReportWatts\": 243\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062650\",\r\n    \"lastReportDate\": 1627599617,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 213,\r\n    \"maxReportWatts\": 240\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061975\",\r\n    \"lastReportDate\": 1627599619,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 214,\r\n    \"maxReportWatts\": 242\r\n  },\r\n  {\r\n    \"serialNumber\": \"202117037990\",\r\n    \"lastReportDate\": 1627599601,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 18,\r\n    \"maxReportWatts\": 245\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062554\",\r\n    \"lastReportDate\": 1627599607,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 13,\r\n    \"maxReportWatts\": 242\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061677\",\r\n    \"lastReportDate\": 1627599604,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 16,\r\n    \"maxReportWatts\": 239\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061455\",\r\n    \"lastReportDate\": 1627599611,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 16,\r\n    \"maxReportWatts\": 241\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061240\",\r\n    \"lastReportDate\": 1627599619,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 217,\r\n    \"maxReportWatts\": 243\r\n  }\r\n]")
		}
		if req.URL.String() == "/production.json" {
			rw.WriteHeader(http.StatusInternalServerError)
		}
		if req.URL.String() == "/api/v1/production" {
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

	data, err := getSystemJSON()
	log.Printf("Returned %s", data)
	if err == nil {
		t.Errorf("expected error to not be nil")
	}
	getMetrics(t, handler, 1448)
}

func TestStreamsSuccess(t *testing.T) {
	handler := initPrometheus()
	defer initEnvoyServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/stream/meter" {
			rw.WriteHeader(http.StatusOK)
			fmt.Fprintf(rw, "data: {\"production\":{\"ph-a\":{\"p\":-0.0,\"q\":138.135,\"s\":139.586,\"v\":118.313,\"i\":1.18,\"pf\":0.0,\"f\":60.0},\"ph-b\":{\"p\":0.0,\"q\":137.861,\"s\":140.002,\"v\":118.371,\"i\":1.182,\"pf\":0.0,\"f\":60.0},\"ph-c\":{\"p\":0.0,\"q\":0.0,\"s\":0.0,\"v\":0.0,\"i\":0.0,\"pf\":0.0,\"f\":0.0}},\"net-consumption\":{\"ph-a\":{\"p\":0.0,\"q\":0.0,\"s\":17.442,\"v\":118.302,\"i\":0.147,\"pf\":0.0,\"f\":60.0},\"ph-b\":{\"p\":-0.0,\"q\":0.0,\"s\":16.803,\"v\":118.371,\"i\":0.141,\"pf\":0.0,\"f\":60.0},\"ph-c\":{\"p\":0.0,\"q\":0.0,\"s\":0.0,\"v\":0.0,\"i\":0.0,\"pf\":0.0,\"f\":0.0}},\"total-consumption\":{\"ph-a\":{\"p\":-0.0,\"q\":-138.135,\"s\":156.955,\"v\":118.307,\"i\":1.327,\"pf\":-0.0,\"f\":60.0},\"ph-b\":{\"p\":0.0,\"q\":-137.861,\"s\":123.278,\"v\":118.371,\"i\":1.041,\"pf\":0.0,\"f\":60.0},\"ph-c\":{\"p\":0.0,\"q\":0.0,\"s\":0.0,\"v\":0.0,\"i\":0.0,\"pf\":0.0,\"f\":0.0}}}\n\n")
		}
	})).Close()
	streams()

	time.Sleep(time.Duration(250) * time.Millisecond)
	getMetrics(t, handler, 3944)
}

func getMetrics(t *testing.T, handler http.Handler, length int) {

	req := httptest.NewRequest("GET", "http://example.com/metrics", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	checkErr(err)

	if len(data) != length {
		log.Println(string(data))
		t.Errorf("data should be %d characters long, is %d", length, len(data))
	}
}

func initEnvoyServer(handler http.HandlerFunc) *httptest.Server {
	os.Setenv("AUTH_TOKEN", "1kdi394js")
	os.Setenv("SLEEP_SECONDS", "10")
	os.Setenv("ARRAY_LAYOUT", "{\"system_id\":2335303,\"rotation\":0,\"dimensions\":{\"x_min\":30,\"x_max\":430,\"y_min\":0,\"y_max\":700},\"arrays\":[{\"array_id\":3871525,\"label\":\"array 1\",\"x\":230,\"y\":350,\"azimuth\":270,\"modules\":[{\"module_id\":48968985,\"rotation\":0,\"x\":300,\"y\":100,\"inverter\":{\"inverter_id\":51116942,\"serial_num\":\"482125061710\"}},{\"module_id\":48968986,\"rotation\":0,\"x\":200,\"y\":100,\"inverter\":{\"inverter_id\":51116946,\"serial_num\":\"482125061458\"}},{\"module_id\":48968987,\"rotation\":0,\"x\":100,\"y\":100,\"inverter\":{\"inverter_id\":51116938,\"serial_num\":\"482125062528\"}},{\"module_id\":48968988,\"rotation\":0,\"x\":0,\"y\":100,\"inverter\":{\"inverter_id\":51116956,\"serial_num\":\"482125062558\"}},{\"module_id\":48968989,\"rotation\":0,\"x\":-100,\"y\":100,\"inverter\":{\"inverter_id\":51116940,\"serial_num\":\"482125062554\"}},{\"module_id\":48968990,\"rotation\":0,\"x\":-200,\"y\":100,\"inverter\":{\"inverter_id\":51116933,\"serial_num\":\"202117037990\"}},{\"module_id\":48968991,\"rotation\":0,\"x\":-300,\"y\":100,\"inverter\":{\"inverter_id\":51116932,\"serial_num\":\"482125062686\"}},{\"module_id\":48968992,\"rotation\":0,\"x\":300,\"y\":-100,\"inverter\":{\"inverter_id\":51116950,\"serial_num\":\"482125061240\"}},{\"module_id\":48968993,\"rotation\":0,\"x\":200,\"y\":-100,\"inverter\":{\"inverter_id\":51116948,\"serial_num\":\"482125062610\"}},{\"module_id\":48968994,\"rotation\":0,\"x\":100,\"y\":-100,\"inverter\":{\"inverter_id\":51116952,\"serial_num\":\"482125061975\"}},{\"module_id\":48968995,\"rotation\":0,\"x\":0,\"y\":-100,\"inverter\":{\"inverter_id\":51116949,\"serial_num\":\"482125062650\"}},{\"module_id\":48968996,\"rotation\":0,\"x\":-100,\"y\":-100,\"inverter\":{\"inverter_id\":51116944,\"serial_num\":\"482125061455\"}},{\"module_id\":48968997,\"rotation\":0,\"x\":-200,\"y\":-100,\"inverter\":{\"inverter_id\":51116936,\"serial_num\":\"482125061677\"}},{\"module_id\":48968998,\"rotation\":0,\"x\":-300,\"y\":-100,\"inverter\":{\"inverter_id\":51116935,\"serial_num\":\"482125062378\"}}],\"dimensions\":{\"x_min\":30,\"x_max\":430,\"y_min\":0,\"y_max\":700},\"tilt\":20,\"array_type_name\":\"\",\"pcu_count\":14,\"pv_module_details\":{\"manufacturer\":\"SunSpark\",\"model\":\"SST-320M3B\",\"type\":null,\"power_rating\":null}}],\"haiku\":\"Put upon the roof I am waiting for the sun All I see is clouds\"}")

	server := httptest.NewServer(handler)

	os.Setenv("ENVOY_URL", server.URL)

	return server
}
