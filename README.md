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

- `ge-enabled`: Enable/disable monitoring for a service
- `ge-timeoutSeconds`: Maximum time to wait for a response
- `ge-intervalSeconds`: How frequently to check the service
- `ge-paths`: List of paths to monitor with HTTP methods (e.g., `GET/health`, `POST/api/v1/test`)

## Installation

### Prerequisites

- Kubernetes cluster (v1.16+)
- `kubectl` configured to access your cluster
- Go 1.16+ (for building from source)

### Quick Start

1. Clone the repository:

   ```bash
   git clone https://github.com/greye-monitoring/greye.git
   cd greye
   ```