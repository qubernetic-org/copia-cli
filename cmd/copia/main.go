package main

import (
	"os"

	"github.com/qubernetic-org/copia-cli/internal/copiacmd"
)

func main() {
	code := copiacmd.Main()
	os.Exit(code)
}
