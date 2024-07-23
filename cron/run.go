package cron

import (
	"fmt"
	"os"

	"github.com/miladin-dev/health-probe-mux/parser"
	pb "github.com/miladin-dev/health-probe-mux/probe"
)

type Cron struct {
	prober *pb.Prober
}

func New() *Cron {
	return &Cron{
		prober: pb.NewProber(),
	}
}

func (c *Cron) RunProbeCron() {
	probes, err := parser.ParseYAML()
	if err != nil {
		fmt.Printf("RunProbeCron: %v", err)
		return
	}
	for _, probe := range probes {
		go func(p *pb.Probe) {
			status := c.prober.RunProbe(p)
			<-status
			// If any of probes finished, kill the binary
			os.Exit(1)
		}(probe)
	}

	// Block to infinity
	select {}
}
