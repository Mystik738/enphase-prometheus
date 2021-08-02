# enphase-prometheus
Reads from an Enphase Envoy S endpoint and converts the JSON into a scrapable Prometheus metrics HTTP endpoint.

## Monitoring with Kubernetes
This repository contains a deployment, service, and service monitor definition that can be deployed to an existing Kubernetes cluster with Prometheus. Values may need to be adjusted to match your Kubernetes configuration.