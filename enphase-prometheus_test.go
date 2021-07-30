package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetEnvoyJson(t *testing.T) {
	os.Setenv("ENVOY_URL", "http://localhost:2983")
	os.Setenv("USERNAME", "envoy")
	os.Setenv("PASSWORD", "123456")

	server := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/api/v1/production/inverters" {
			user, password, ok := req.BasicAuth()

			if ok != true {
				t.Error("Invalid http basic auth credentials")
			}
			if user != "envoy" {
				t.Error("Incorrect http basic username")
			}
			if password != "123456" {
				t.Error("Incorrect http basic password")
			}

			rw.WriteHeader(http.StatusOK)
			fmt.Fprintf(rw, "[\r\n  {\r\n    \"serialNumber\": \"482125062378\",\r\n    \"lastReportDate\": 1627599602,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 14,\r\n    \"maxReportWatts\": 240\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061710\",\r\n    \"lastReportDate\": 1627599608,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 223,\r\n    \"maxReportWatts\": 245\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062686\",\r\n    \"lastReportDate\": 1627599597,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 16,\r\n    \"maxReportWatts\": 244\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062528\",\r\n    \"lastReportDate\": 1627599606,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 219,\r\n    \"maxReportWatts\": 245\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061458\",\r\n    \"lastReportDate\": 1627599613,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 221,\r\n    \"maxReportWatts\": 245\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062610\",\r\n    \"lastReportDate\": 1627599613,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 220,\r\n    \"maxReportWatts\": 244\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062558\",\r\n    \"lastReportDate\": 1627599623,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 43,\r\n    \"maxReportWatts\": 243\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062650\",\r\n    \"lastReportDate\": 1627599617,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 213,\r\n    \"maxReportWatts\": 240\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061975\",\r\n    \"lastReportDate\": 1627599619,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 214,\r\n    \"maxReportWatts\": 242\r\n  },\r\n  {\r\n    \"serialNumber\": \"202117037990\",\r\n    \"lastReportDate\": 1627599601,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 18,\r\n    \"maxReportWatts\": 245\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125062554\",\r\n    \"lastReportDate\": 1627599607,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 13,\r\n    \"maxReportWatts\": 242\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061677\",\r\n    \"lastReportDate\": 1627599604,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 16,\r\n    \"maxReportWatts\": 239\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061455\",\r\n    \"lastReportDate\": 1627599611,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 16,\r\n    \"maxReportWatts\": 241\r\n  },\r\n  {\r\n    \"serialNumber\": \"482125061240\",\r\n    \"lastReportDate\": 1627599619,\r\n    \"devType\": 1,\r\n    \"lastReportWatts\": 217,\r\n    \"maxReportWatts\": 243\r\n  }\r\n]")
		}
	}))

	l, err := net.Listen("tcp", "127.0.0.1:2983")
	if err != nil {
		log.Fatal(err)
	}
	server.Listener.Close()
	server.Listener = l
	server.Start()
	defer server.Close()

	json := getEnvoyJson()

	log.Println(json)

}
