package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ToomyGO/pasunes_scan/internal/cli"
	"github.com/ToomyGO/pasunes_scan/internal/pipeline"

	"github.com/ToomyGO/pasunes_scan/internal/scanner"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	//Test and run PipeLine

	in := make(chan int)

	doubleStage := func(in <-chan int) <-chan int {
		out := make(chan int)

		go func() {
			defer close(out)

			for v := range in {
				time.Sleep(100 * time.Millisecond)
				out <- v * 2
			}
		}()
		return out
	}

	outs := pipeline.FanOut(in, 3, doubleStage)

	merged := pipeline.FanIn(ctx, outs...)

	go func() {
		for i := 1; i <= 6; i++ {
			in <- i
		}
		close(in)
	}()

	for v := range merged {
		fmt.Println(v)
	}

	config := cli.ParsFlags()
	if config == nil {
		return
	}

	fmt.Printf("🎯 Target: %s\n", config.Target)
	if config.Verbose {
		fmt.Println("🔍 Verbose mode: enabled")
	}

	result, err := scanner.ResolvedTarget(config.Target, config.Verbose)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("Resolved: %s -> %s (isIP=%v)\n", result.Input, result.Resolved, result.IsIP)
}
