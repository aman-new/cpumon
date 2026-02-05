package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	interval := flag.Duration("i", time.Second, "refresh interval (min 100ms)")
	flag.Parse()

	if *interval < 100*time.Millisecond {
		fmt.Fprintln(os.Stderr, "Error: interval must be at least 100ms")
		os.Exit(1)
	}

	m := NewMonitor()
	if err := m.Run(context.Background(), *interval); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
