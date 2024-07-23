package probe

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	v1 "k8s.io/api/core/v1"
	k8sprobe "k8s.io/kubernetes/pkg/probe"
	httpprobe "k8s.io/kubernetes/pkg/probe/http"
	tcpprobe "k8s.io/kubernetes/pkg/probe/tcp"
)

const maxRetries = 3

type Probe struct {
	*v1.Probe
	// TODO: Add option to expose the status on port
	// ExposedPort int `json:"exposeOnPort,omitempty"`
}

type Prober struct {
	TCP  tcpprobe.Prober
	HTTP httpprobe.Prober
}

func NewProber() *Prober {
	return &Prober{
		TCP:  tcpprobe.New(),
		HTTP: httpprobe.New(false),
	}
}

func (p *Prober) RunProbe(probe *Probe) <-chan struct{} {
	end := make(chan struct{})
	go func(probe *Probe) {
		period := time.Duration(probe.PeriodSeconds) * time.Second
		periodTicker := time.NewTicker(time.Duration(period))
		failureThreshold := probe.FailureThreshold
		for {
			result := p.run(probe)
			if result == k8sprobe.Failure {
				failureThreshold--
				if failureThreshold <= 0 {
					end <- struct{}{}
				}
			} else {
				failureThreshold = probe.FailureThreshold
			}
			<-periodTicker.C
		}
	}(probe)
	return end
}

func (p *Prober) run(probe *Probe) k8sprobe.Result {
	result, output, err := p.runWithRetries(probe)
	if err != nil {
		fmt.Printf("[runWithRetries]: %v\n", err)
	}
	switch result {
	case k8sprobe.Success:
	case k8sprobe.Failure:
		fmt.Printf("[FAILURE] probe execution: %v\n", output)
	case k8sprobe.Warning:
		fmt.Printf("[WARNING] probe execution: %s\n", output)
	case k8sprobe.Unknown:
		fmt.Printf("[UNKNOWN] probe execution: %s\n", output)
	}
	return result
}

func (p *Prober) runWithRetries(probe *Probe) (k8sprobe.Result, string, error) {
	var result k8sprobe.Result
	var output string
	var err error
	for i := 0; i < maxRetries; i++ {
		result, output, err = p.runProbe(probe)
		if err == nil {
			return result, output, nil
		}
	}
	return result, output, err
}

func (p *Prober) runProbe(probe *Probe) (k8sprobe.Result, string, error) {
	if probe.TCPSocket != nil {
		return p.runTCPProbe(probe)
	}
	if probe.HTTPGet != nil {
		return p.runHTTPProbe(probe)
	}
	return k8sprobe.Unknown, "", fmt.Errorf("unknown probe format")
}

func (p *Prober) runTCPProbe(probe *Probe) (k8sprobe.Result, string, error) {
	timeoutSeconds := probe.TimeoutSeconds
	if timeoutSeconds == 0 {
		timeoutSeconds = 1
	}
	timeoutDur := time.Duration(timeoutSeconds) * time.Second
	host := probe.TCPSocket.Host
	if host == "" {
		host = "localhost"
	}
	port := int(probe.TCPSocket.Port.IntVal)
	return p.TCP.Probe(host, port, timeoutDur)
}

func (p *Prober) runHTTPProbe(probe *Probe) (k8sprobe.Result, string, error) {
	timeoutSeconds := probe.TimeoutSeconds
	if timeoutSeconds == 0 {
		timeoutSeconds = 1
	}
	timeoutDur := time.Duration(timeoutSeconds) * time.Second
	port := probe.HTTPGet.Port.String()
	host := probe.HTTPGet.Host
	if host == "" {
		host = "localhost"
	}
	path := probe.HTTPGet.Path
	scheme := "http"
	if probe.HTTPGet.Scheme == v1.URISchemeHTTPS {
		scheme = "https"
	}
	url := &url.URL{
		Scheme: scheme,
		Host:   net.JoinHostPort(host, port),
		Path:   path,
	}
	req := &http.Request{
		Method: http.MethodGet,
		URL:    url,
	}
	return p.HTTP.Probe(req, timeoutDur)
}
