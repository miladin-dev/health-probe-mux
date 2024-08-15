# 🚀 **health-probe-mux** &middot; ![Go Version](https://img.shields.io/badge/Go-1.20-blue) ![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)

**health-probe-mux** is a CLI tool designed to extend Kubernetes' native health checks by allowing multiple health checks to be part of a single health probe. 

## **Usage**
Imagine you're running an application that relies on an external service or database, or requires multiple resources to be considered healthy. In this situations native health checks might not be enough.

## **Status**
**Not ready for production use**\
Currently, only the CLI is available, without the ability to expose results on a port. The roadmap includes plans to enable the binary to run inside a sidecar container and expose health check results on a specified port, allowing your application to monitor this data and act accordingly.

## **Roadmap**
- [ ] Add GitHub Actions
- [ ] Make binary available as a Docker image
- [ ] Enhance functionality and performance improvements

## **License**
This project was developed as part of my master’s thesis on "The Analysis of Containerization Principles and Internal Architecture of Kubernetes."\
Distributed under the MIT License. See the `LICENSE` file for more details.

