package main

import (
	"context"

	"github.com/miladin-dev/health-probe-mux/cron"
)

func runAction(ctx context.Context) error {
	job := cron.New()
	job.RunProbeCron()
	return nil
}

func main() {
	runAction(context.Background())
}
