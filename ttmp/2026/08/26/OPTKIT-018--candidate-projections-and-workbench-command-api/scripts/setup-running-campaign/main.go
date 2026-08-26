package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/the-tree-center/rag-ttc/pkg/ttc/optkitcampaign"
)

func main() {
	store := flag.String("store", "", "campaign store root")
	flag.Parse()
	if *store == "" {
		fmt.Fprintln(os.Stderr, "--store is required")
		os.Exit(2)
	}
	arms, cases, err := optkitcampaign.SemanticFixturePlan()
	if err != nil {
		panic(err)
	}
	summary, err := optkitcampaign.Run(context.Background(), optkitcampaign.RunOptions{
		Root: *store, Reset: true, Arms: arms, Cases: cases, Repeats: 1,
		Executor: optkitcampaign.SemanticFixtureExecutor{},
		StopAfter: optkitcampaign.CheckpointAfterResult,
	})
	if !errors.Is(err, optkitcampaign.ErrCheckpoint) {
		panic(fmt.Sprintf("expected checkpoint interruption, got %v", err))
	}
	fmt.Println(summary.Campaign)
}
