# enphase-prometheus
Reads from an Enphase Envoy S endpoint and converts the JSON into a scrapable Prometheus metrics HTTP endpoint.

## Environment Variables

| Name | Description |
| ------------- | ----------- |
| ENVOY_URL | The base url of your envoy system, i.e. http://envoy.local |
| AUTH_TOKEN   | JWT to log into the envoy system. Considered sensitive.       |
| ARRAY_LAYOUT | Layout of array as defined by https://enlighten.enphaseenergy.com/pv/systems/${SYSTEM_ID}/array_layout_x.json |
| SLEEP_SECONDS | Amount of seconds the scraper waits between scrapes. Defaults to 10 seconds if not set |
| PORT | Port used by webserver. Defaults to 80 if not set |

## Monitoring with Kubernetes
This repository contains a deployment, service, and service monitor definition that can be deployed to an existing Kubernetes cluster with Prometheus. Values may need to be adjusted to match your Kubernetes configuration. A sample secret is provided, but the credentials will need to be encoded into the file before application.

## Building this image

Ensure you can run multi-architecture builds, then run the following:

docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 --push -t mystik738/enphase-prometheus:latest -t mystik738/enphase-prometheus:${VERSION} .