package parser

import (
	"encoding/json"
	"fmt"

	"github.com/miladin-dev/health-probe-mux/probe"
	"sigs.k8s.io/yaml"
)

const yamlExample = `
probe:
  - tcpSocket:
      port: 8080
    failureThreshold: 1
    periodSeconds: 1   
  - httpGet:
      port: 8081
    failureThreshold: 1
    periodSeconds: 1
`

type convertedProbe struct {
	// ExposeOnPort int            `json:"exposeOnPort,omitempty"`
	Probes []*probe.Probe `json:"probe,omitempty"`
}

func ParseYAML(content []byte) ([]*probe.Probe, error) {
	probeByte, err := yaml.YAMLToJSON(content)
	if err != nil {
		return nil, fmt.Errorf("error converting yaml to json: %w", err)
	}

	probe := &convertedProbe{}
	err = json.Unmarshal(probeByte, probe)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal to struct: %v", err)
	}
	return probe.Probes, nil
}
