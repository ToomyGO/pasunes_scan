package cli

import (
	"flag"
	"fmt"
)

type Config struct {
	Target  string
	Verbose bool
}

func ParsFlags() *Config {

	target := flag.String("target", "", "Target domain or IP")
	verbose := flag.Bool("v", false, "Enable verbose output")

	flag.Parse()

	if *target == "" {
		fmt.Println("❌ Please provide a target using -target flag.")
		flag.Usage()
		return nil
	}
	return &Config{
		Target:  *target,
		Verbose: *verbose,
	}
}
