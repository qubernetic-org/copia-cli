package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra/doc"
	"github.com/qubernetic/copia-cli/internal/config"
	"github.com/qubernetic/copia-cli/internal/copiacmd"
	"github.com/qubernetic/copia-cli/pkg/cmdutil"
	"github.com/qubernetic/copia-cli/pkg/iostreams"
)

func main() {
	outDir := "docs/manual/src/commands"
	if len(os.Args) > 1 {
		outDir = os.Args[1]
	}

	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output dir: %v\n", err)
		os.Exit(1)
	}

	ios := iostreams.System()
	f := &cmdutil.Factory{
		IOStreams: ios,
		Config: func() (*config.Config, error) {
			return &config.Config{Hosts: map[string]*config.HostConfig{}}, nil
		},
	}

	rootCmd := copiacmd.NewRootCmd(f)
	rootCmd.DisableAutoGenTag = true

	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, ".md")
		return base + ".md"
	}

	filePrepender := func(filename string) string {
		name := filepath.Base(filename)
		name = strings.TrimSuffix(name, ".md")
		name = strings.ReplaceAll(name, "_", " ")
		return fmt.Sprintf("# %s\n\n", name)
	}

	if err := doc.GenMarkdownTreeCustom(rootCmd, outDir, filePrepender, linkHandler); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating docs: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated command docs in %s\n", outDir)
}
