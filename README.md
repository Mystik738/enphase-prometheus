# enphase-prometheus
Reads from an Enphase Envoy S endpoint and converts the JSON into a scrapable Prometheus metrics HTTP endpoint.

## Environment Variables

| Name | Description |
| ------------- | ----------- |
| SLEEP_SECONDS | Amount of seconds the scraper waits between scrapes. Defaults to 10 seconds if not set |
| ENVOY_URL | The base url of your envoy system, i.e. http://envoy.local |
| USERNAME        | Username to log into the envoy system. Considered sensitive.       |
| PASSWORD     | Password to log into the envoy system. Considered sensitive.        |

## Monitoring with Kubernetes
This repository contains a deployment, service, and service monitor definition that can be deployed to an existing Kubernetes cluster with Prometheus. Values may need to be adjusted to match your Kubernetes configuration. A sample secret is provided, but the credentials will need to be encoded into the file before application.