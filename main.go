package main

import (
	"context"
	"fmt"
	"os"

	"github.com/miladin-dev/health-probe-mux/cron"
)

func runAction(ctx context.Context) {
	job := cron.New()
	if len(os.Args) <= 1 {
		fmt.Println("not enough arguments")
		return
	}
	if err := job.RunProbeCron(os.Args[1:]); err != nil {
		fmt.Printf("unable to run cron job: %v", err)
		return
	}
}

func main() {
	ctx := context.Background()
	runAction(ctx)
}
