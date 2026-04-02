package docs

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenMarkdownCustom_FrontMatter(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test command",
	}

	var buf bytes.Buffer
	err := genMarkdownCustom(cmd, &buf, func(s string) string { return s })
	require.NoError(t, err)

	// Should NOT contain front matter — that's added by filePrepender
	// genMarkdownCustom only generates content
	assert.Contains(t, buf.String(), "## test")
	assert.Contains(t, buf.String(), "Test command")
}

func TestGenMarkdownCustom_NoBranding(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "copia-cli",
		Short: "Copia CLI",
		Long:  "Work with Copia repositories.",
	}

	var buf bytes.Buffer
	err := genMarkdownCustom(cmd, &buf, func(s string) string { return s })
	require.NoError(t, err)

	output := buf.String()
	assert.NotContains(t, output, "github.com")
	assert.NotContains(t, output, "GitHub Enterprise")
	assert.NotContains(t, output, " gh ")
}

func TestGenMarkdownCustom_FlagsHTML(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List items",
		RunE:  func(cmd *cobra.Command, args []string) error { return nil },
	}
	cmd.Flags().StringP("state", "s", "open", "Filter by state")
	cmd.Flags().IntP("limit", "L", 30, "Maximum items")

	var buf bytes.Buffer
	err := genMarkdownCustom(cmd, &buf, func(s string) string { return s })
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, `<dl class="flags">`)
	assert.Contains(t, output, `--state`)
	assert.Contains(t, output, `--limit`)
	assert.Contains(t, output, `</dl>`)
}

func TestGenMarkdownCustom_Examples(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "close",
		Short: "Close an issue",
		RunE:  func(cmd *cobra.Command, args []string) error { return nil },
		Example: `# Close an issue
copia issue close 12

# Close with comment
copia issue close 12 --comment "Fixed"`,
	}

	var buf bytes.Buffer
	err := genMarkdownCustom(cmd, &buf, func(s string) string { return s })
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "### Examples")
	assert.Contains(t, output, "Close an issue")
	assert.Contains(t, output, "copia issue close 12")
}

func TestGenMarkdownCustom_Aliases(t *testing.T) {
	parent := &cobra.Command{Use: "issue", Short: "Manage issues"}
	child := &cobra.Command{
		Use:     "list",
		Short:   "List issues",
		Aliases: []string{"ls"},
		RunE:    func(cmd *cobra.Command, args []string) error { return nil },
	}
	parent.AddCommand(child)

	var buf bytes.Buffer
	err := genMarkdownCustom(child, &buf, func(s string) string { return s })
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Aliases")
	assert.Contains(t, output, "ls")
}

func TestGenMarkdownCustom_JSONFields(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List items",
		RunE:  func(cmd *cobra.Command, args []string) error { return nil },
		Annotations: map[string]string{
			"help:json-fields": "number,title,state,labels",
		},
	}

	var buf bytes.Buffer
	err := genMarkdownCustom(cmd, &buf, func(s string) string { return s })
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "### JSON Fields")
	assert.Contains(t, output, "`number`")
	assert.Contains(t, output, "`title`")
}

func TestGenMarkdownCustom_SubcommandGroups(t *testing.T) {
	parent := &cobra.Command{Use: "issue", Short: "Manage issues"}
	parent.AddGroup(&cobra.Group{ID: "general", Title: "General commands"})
	parent.AddGroup(&cobra.Group{ID: "targeted", Title: "Targeted commands"})

	create := &cobra.Command{Use: "create", Short: "Create an issue", GroupID: "general",
		RunE: func(cmd *cobra.Command, args []string) error { return nil }}
	list := &cobra.Command{Use: "list", Short: "List issues", GroupID: "general",
		RunE: func(cmd *cobra.Command, args []string) error { return nil }}
	close := &cobra.Command{Use: "close", Short: "Close an issue", GroupID: "targeted",
		RunE: func(cmd *cobra.Command, args []string) error { return nil }}
	view := &cobra.Command{Use: "view", Short: "View an issue", GroupID: "targeted",
		RunE: func(cmd *cobra.Command, args []string) error { return nil }}

	parent.AddCommand(create, list, close, view)

	var buf bytes.Buffer
	err := genMarkdownCustom(parent, &buf, func(s string) string { return s })
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "### General commands")
	assert.Contains(t, output, "### Targeted commands")
	assert.True(t, strings.Index(output, "General commands") < strings.Index(output, "Targeted commands"))
}

func TestGenMarkdownCustom_SeeAlso(t *testing.T) {
	parent := &cobra.Command{Use: "issue", Short: "Manage issues"}
	child := &cobra.Command{
		Use:   "close",
		Short: "Close an issue",
		RunE:  func(cmd *cobra.Command, args []string) error { return nil },
	}
	parent.AddCommand(child)

	var buf bytes.Buffer
	err := genMarkdownCustom(child, &buf, func(s string) string { return s })
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "### See also")
	assert.Contains(t, output, "issue")
}

func TestCmdManualPath(t *testing.T) {
	parent := &cobra.Command{Use: "copia-cli"}
	child := &cobra.Command{Use: "issue"}
	sub := &cobra.Command{Use: "close"}
	parent.AddCommand(child)
	child.AddCommand(sub)

	assert.Equal(t, "copia-cli_issue_close.md", cmdManualPath(sub))
}

func TestFormatSlice(t *testing.T) {
	result := formatSlice([]string{"b", "a", "c"}, "`", "`", true)
	assert.Contains(t, result, "`a`")
	assert.Contains(t, result, "`b`")
	assert.Contains(t, result, "`c`")
}
