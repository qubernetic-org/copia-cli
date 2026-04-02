// Package docs generates Jekyll-compatible markdown documentation for Cobra commands.
// Adapted from github.com/cli/cli/v2/internal/docs with GitHub-specific content removed.
//
//nolint:errcheck // fmt.Fprint write errors are intentionally ignored, matching gh CLI upstream style.
package docs

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// formatSlice concatenates elements into a comma-separated string.
// Elements are wrapped with prepend/append strings. If doSort is true, values are sorted.
func formatSlice(values []string, prependWith, appendWith string, doSort bool) string {
	if doSort {
		values = slices.Clone(values)
		slices.Sort(values)
	}

	var parts []string
	for _, v := range values {
		parts = append(parts, prependWith+strings.TrimSpace(v)+appendWith)
	}
	return strings.Join(parts, ", ")
}

// buildAliasList builds full alias paths for a command.
func buildAliasList(cmd *cobra.Command, aliases []string) []string {
	if !cmd.HasParent() {
		return aliases
	}

	parentPath := cmd.Parent().CommandPath()
	var result []string
	for _, a := range aliases {
		result = append(result, parentPath+" "+a)
	}
	return result
}

// CommandGroup represents a group of subcommands with a title.
type CommandGroup struct {
	Title    string
	Commands []*cobra.Command
}

// groupedCommands returns subcommands organized by their Cobra GroupID.
func groupedCommands(cmd *cobra.Command) []CommandGroup {
	var res []CommandGroup

	for _, g := range cmd.Groups() {
		var cmds []*cobra.Command
		for _, c := range cmd.Commands() {
			if c.GroupID == g.ID && c.IsAvailableCommand() {
				cmds = append(cmds, c)
			}
		}
		if len(cmds) > 0 {
			res = append(res, CommandGroup{
				Title:    g.Title,
				Commands: cmds,
			})
		}
	}

	var cmds []*cobra.Command
	for _, c := range cmd.Commands() {
		if c.GroupID == "" && c.IsAvailableCommand() {
			cmds = append(cmds, c)
		}
	}
	if len(cmds) > 0 {
		defaultTitle := "Additional commands"
		if len(cmd.Groups()) == 0 {
			defaultTitle = "Available commands"
		}
		res = append(res, CommandGroup{
			Title:    defaultTitle,
			Commands: cmds,
		})
	}

	return res
}

func printJSONFields(w io.Writer, cmd *cobra.Command) {
	raw, ok := cmd.Annotations["help:json-fields"]
	if !ok {
		return
	}

	fmt.Fprint(w, "### JSON Fields\n\n")
	fmt.Fprint(w, formatSlice(strings.Split(raw, ","), "`", "`", true))
	fmt.Fprint(w, "\n\n")
}

func printAliases(w io.Writer, cmd *cobra.Command) {
	if len(cmd.Aliases) == 0 {
		return
	}

	fmt.Fprint(w, "### Aliases\n\n")
	aliases := buildAliasList(cmd, cmd.Aliases)
	sort.Strings(aliases)
	for _, a := range aliases {
		fmt.Fprintf(w, "  %s\n", a)
	}
	fmt.Fprint(w, "\n\n")
}

func printOptions(w io.Writer, cmd *cobra.Command) error {
	flags := cmd.NonInheritedFlags()
	flags.SetOutput(w)
	if flags.HasAvailableFlags() {
		fmt.Fprint(w, "### Options\n\n")
		if err := printFlagsHTML(w, flags); err != nil {
			return err
		}
		fmt.Fprint(w, "\n\n")
	}

	parentFlags := cmd.InheritedFlags()
	parentFlags.SetOutput(w)
	if hasNonHelpFlags(parentFlags) {
		fmt.Fprint(w, "### Options inherited from parent commands\n\n")
		if err := printFlagsHTML(w, parentFlags); err != nil {
			return err
		}
		fmt.Fprint(w, "\n\n")
	}
	return nil
}

func hasNonHelpFlags(fs *pflag.FlagSet) (found bool) {
	fs.VisitAll(func(f *pflag.Flag) {
		if !f.Hidden && f.Name != "help" {
			found = true
		}
	})
	return
}

var hiddenFlagDefaults = map[string]bool{
	"false": true,
	"":      true,
	"[]":    true,
	"0s":    true,
}

var defaultValFormats = map[string]string{
	"string":   " (default \"%s\")",
	"duration": " (default \"%s\")",
}

func getDefaultValueDisplayString(f *pflag.Flag) string {
	if hiddenFlagDefaults[f.DefValue] || hiddenFlagDefaults[f.Value.Type()] {
		return ""
	}

	if dvf, found := defaultValFormats[f.Value.Type()]; found {
		return fmt.Sprintf(dvf, f.Value)
	}
	return fmt.Sprintf(" (default %s)", f.Value)
}

type flagView struct {
	Name      string
	Varname   string
	Shorthand string
	DefValue  string
	Usage     string
}

var flagsTemplate = `
<dl class="flags">{{ range . }}
	<dt>{{ if .Shorthand }}<code>-{{.Shorthand}}</code>, {{ end }}
		<code>--{{.Name}}{{ if .Varname }} &lt;{{.Varname}}&gt;{{ end }}{{.DefValue}}</code></dt>
	<dd>{{.Usage}}</dd>
{{ end }}</dl>
`

var tpl = template.Must(template.New("flags").Parse(flagsTemplate))

func printFlagsHTML(w io.Writer, fs *pflag.FlagSet) error {
	var flags []flagView
	fs.VisitAll(func(f *pflag.Flag) {
		if f.Hidden || f.Name == "help" {
			return
		}
		varname, usage := pflag.UnquoteUsage(f)

		flags = append(flags, flagView{
			Name:      f.Name,
			Varname:   varname,
			Shorthand: f.Shorthand,
			DefValue:  getDefaultValueDisplayString(f),
			Usage:     usage,
		})
	})
	return tpl.Execute(w, flags)
}

// genMarkdownCustom creates custom markdown output for a single command.
func genMarkdownCustom(cmd *cobra.Command, w io.Writer, linkHandler func(string) string) error {
	fmt.Fprint(w, "{% raw %}")
	fmt.Fprintf(w, "## %s\n\n", cmd.CommandPath())

	hasLong := cmd.Long != ""
	if !hasLong {
		fmt.Fprintf(w, "%s\n\n", cmd.Short)
	}
	if cmd.Runnable() {
		fmt.Fprintf(w, "```\n%s\n```\n\n", cmd.UseLine())
	}
	if hasLong {
		fmt.Fprintf(w, "%s\n\n", cmd.Long)
	}

	for _, g := range groupedCommands(cmd) {
		fmt.Fprintf(w, "### %s\n\n", g.Title)
		for _, subcmd := range g.Commands {
			fmt.Fprintf(w, "* [%s](%s)\n", subcmd.CommandPath(), linkHandler(cmdManualPath(subcmd)))
		}
		fmt.Fprint(w, "\n\n")
	}

	if err := printOptions(w, cmd); err != nil {
		return err
	}
	printAliases(w, cmd)
	printJSONFields(w, cmd)
	fmt.Fprint(w, "{% endraw %}\n")

	if len(cmd.Example) > 0 {
		fmt.Fprint(w, "### Examples\n\n{% highlight bash %}{% raw %}\n")
		fmt.Fprint(w, cmd.Example)
		fmt.Fprint(w, "{% endraw %}{% endhighlight %}\n\n")
	}

	if cmd.HasParent() {
		p := cmd.Parent()
		fmt.Fprint(w, "### See also\n\n")
		fmt.Fprintf(w, "* [%s](%s)\n", p.CommandPath(), linkHandler(cmdManualPath(p)))
	}

	return nil
}

// GenMarkdownTreeCustom generates markdown for the entire command tree.
func GenMarkdownTreeCustom(cmd *cobra.Command, dir string, filePrepender, linkHandler func(string) string) error {
	for _, c := range cmd.Commands() {
		_, forceGeneration := c.Annotations["markdown:generate"]
		if c.Hidden && !forceGeneration {
			continue
		}

		if err := GenMarkdownTreeCustom(c, dir, filePrepender, linkHandler); err != nil {
			return err
		}
	}

	filename := filepath.Join(dir, cmdManualPath(cmd))
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	if _, err := io.WriteString(f, filePrepender(filename)); err != nil {
		return err
	}
	if err := genMarkdownCustom(cmd, f, linkHandler); err != nil {
		return err
	}
	return nil
}

func cmdManualPath(c *cobra.Command) string {
	if basenameOverride, found := c.Annotations["markdown:basename"]; found {
		return basenameOverride + ".md"
	}
	return strings.ReplaceAll(c.CommandPath(), " ", "_") + ".md"
}
