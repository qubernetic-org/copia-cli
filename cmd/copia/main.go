package main

import (
	"fmt"
	"os"

	"github.com/qubernetic-org/copia-cli/internal/build"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("copia version %s (%s)\n", build.Version, build.Date)
		os.Exit(0)
	}
	fmt.Println("copia: command-line interface for Copia")
	os.Exit(0)
}
