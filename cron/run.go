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

func (c *Cron) RunProbeCron(args []string) error {
	fileContent, err := os.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("unable to read configuration file: %v", err)
	}
	probes, err := parser.ParseYAML(fileContent)
	if err != nil {
		return fmt.Errorf("unable to parse yaml: %v", err)
	}
	for _, probe := range probes {
		go func(p *pb.Probe) {
			status := c.prober.RunProbe(p)
			<-status
			// If any of probes finish, kill the binary
			os.Exit(1)
		}(probe)
	}

	// Block to infinity
	select {}
}
