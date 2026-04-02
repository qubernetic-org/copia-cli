package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/qubernetic/copia-cli/internal/config"
	"github.com/qubernetic/copia-cli/internal/copiacmd"
	"github.com/qubernetic/copia-cli/internal/docs"
	"github.com/qubernetic/copia-cli/pkg/cmdutil"
	"github.com/qubernetic/copia-cli/pkg/iostreams"
)

func main() {
	outDir := "docs/site/manual"
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
		return "./" + strings.TrimSuffix(name, ".md")
	}

	filePrepender := func(_ string) string {
		return "---\nlayout: manual\npermalink: /:path/:basename\n---\n\n"
	}

	if err := docs.GenMarkdownTreeCustom(rootCmd, outDir, filePrepender, linkHandler); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating docs: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated docs in %s\n", outDir)
}
