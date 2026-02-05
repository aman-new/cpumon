package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	m := NewMonitor()
	if err := m.Run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
