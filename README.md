# greye

<p align="center">
  <img src="assets/logo.jpeg" alt="logo" width="500">
</p>

## Introduction

Greye is a Kubernetes-native tool designed to monitor the availability and performance of services and the availability of other clusters.
It provides real-time insights into service health, supports customizable alerting, and offers flexible integration options to ensure your infrastructure remains reliable and performant.

## Features

- **Service Monitoring**: Track availability and response times of Kubernetes services
- **Annotation-Based Configuration**: Configure monitoring parameters directly via Kubernetes annotations
- **Multi-Cluster Support**: Monitor multiple Kubernetes clusters
- **Customizable Alerting**: Configure thresholds and notification channels
- **Flexible Path Monitoring**: Test specific API endpoints with custom HTTP methods
- **High Performance**: Optimized for minimal resource usage and maximum reliability

## How It Works

Greye watches Kubernetes Service resources and uses annotations to determine monitoring parameters:

| Annotation               | Description                                                                                  |
|--------------------------|----------------------------------------------------------------------------------------------|
| `ge-enabled`             | Enable/disable monitoring for a service                                                      |
| `ge-timeoutSeconds`      | Maximum time to wait for a response                                                          |
| `ge-intervalSeconds`     | How frequently to check the service                                                          |
| `ge-paths`               | List of paths to monitor with HTTP methods. If no methods are specified, use the default one |
| `ge-port`                | Port to monitor                                                                              |
| `ge-protocol`            | Protocol to use for monitoring                                                               |
| `ge-headers`             | Custom headers to include in requests                                                        |
| `ge-body`                | Custom body to include in requests                                                           |
| `ge-maxFailedRequests`   | Maximum number of failed requests before alerting                                            |
| `ge-stopMonitoringUntil` | Stop monitoring until a specific time                                                        |
| `ge-forcePodMonitor`     | Force monitoring of pods instead of services                                                 |
| `ge-authentication`      | Authentication method to use for monitoring                                                  |
| `ge-authUsername`        | Username for authentication                                                                  |
| `ge-authPassword`        | Password for authentication                                                                  |



## Example

```yaml
apiVersion: v1
kind: Service
metadata:
   annotations:
      ge-enabled: 'true'
      ge-intervalSeconds: '60'
      ge-paths: |-
         POST/example
         /example/1
      ge-body: '{"key": "value"}'
...
```


## Installation

There is a helm chart available to install Greye in your cluster. You can find it [here](https://github.com/greye-monitoring/helm-charts)


## Contributing
