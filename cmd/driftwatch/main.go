package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/manifest"
	"github.com/yourorg/driftwatch/internal/reporter"
)

func main() {
	manifestPath := flag.String("manifest", "", "Path to the service manifest YAML file (required)")
	outputFormat := flag.String("format", "text", "Output format: text or json")
	flag.Parse()

	if *manifestPath == "" {
		fmt.Fprintln(os.Stderr, "error: --manifest flag is required")
		flag.Usage()
		os.Exit(1)
	}

	m, err := manifest.Load(*manifestPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading manifest: %v\n", err)
		os.Exit(1)
	}

	drifts, err := drift.Detect(m)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error detecting drift: %v\n", err)
		os.Exit(1)
	}

	w := reporter.NewWriter(os.Stdout)

	switch *outputFormat {
	case "json":
		if err := w.WriteJSON(drifts); err != nil {
			fmt.Fprintf(os.Stderr, "error writing JSON report: %v\n", err)
			os.Exit(1)
		}
	case "text":
		if err := w.WriteText(drifts); err != nil {
			fmt.Fprintf(os.Stderr, "error writing text report: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "error: unknown format %q, must be 'text' or 'json'\n", *outputFormat)
		os.Exit(1)
	}

	if len(drifts) > 0 {
		os.Exit(2)
	}
}
