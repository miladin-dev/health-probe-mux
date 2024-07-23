# health-probe-mux
health-probe-mux is a CLI tool that extends Kubernetes' native health checks. It allows for multiple health checks to be part of the same health probe. 

### Usage
Suppose you run an application that depends on an external service or database, and your application cannot function properly without it. With native health checks, it's not possible to check both the container's health and the health of an external source simultaneously.

### Status:
Not ready for production use.\
At the moment, only the CLI is available without the option to expose the results on a port. 
The plan is to make the binary run inside a sidecar container and expose the results of the defined configuration on a specified port. The application container can then monitor the traffic on this port and act accordingly.

### Roadmap:
- Add github actions for tests
- Make binary available as docker image
- Improve various things

### License
This project was implemented as part of my master’s thesis on the topic "The Analysis of Containerization Principles and Internal Architecture of Kubernetes."\
Distributed under the MIT License. See the `LICENSE` file for details.
