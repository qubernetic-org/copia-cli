package main

import (
	"os"

	"github.com/qubernetic/copia-cli/internal/copiacmd"
)

func main() {
	code := copiacmd.Main()
	os.Exit(code)
}
