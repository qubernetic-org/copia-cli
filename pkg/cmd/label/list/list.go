package list

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/qubernetic/copia-cli/pkg/cmdutil"
	"github.com/qubernetic/copia-cli/pkg/iostreams"
)

var validJSONFields = []string{"name", "color", "description"}

// ListOptions holds all inputs for the label list command.
type ListOptions struct {
	IO         *iostreams.IOStreams
	HTTPClient *http.Client
	Host       string
	Token      string
	Owner      string
	Repo       string
	JSON       cmdutil.JSONFlags
}

type labelEntry struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

// NewCmdList creates the `copia label list` command.
func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	opts := &ListOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List labels in a repository",
		Long:    "List labels in a Copia repository. Labels are displayed with their name, color, and description.",
		Aliases: []string{"ls"},
		Example: `  $ copia-cli label list
  $ copia-cli label list --json name,color`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.IO = f.IOStreams
			host, token, err := f.ResolveAuth()
			if err != nil {
				return err
			}
			opts.Host = host
			opts.Token = token

			owner, repo, err := f.ResolveRepo()
			if err != nil {
				return err
			}
			opts.Owner = owner
			opts.Repo = repo
			opts.HTTPClient = &http.Client{}
			return ListRun(opts)
		},
	}

	cmdutil.AddJSONFlags(cmd, &opts.JSON, validJSONFields)

	return cmd
}

func ListRun(opts *ListOptions) error {
	url := fmt.Sprintf("https://%s/api/v1/repos/%s/%s/labels", opts.Host, opts.Owner, opts.Repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+opts.Token)

	resp, err := opts.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("connecting to %s: %w", opts.Host, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error (HTTP %d)", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var labels []labelEntry
	if err := json.Unmarshal(body, &labels); err != nil {
		return fmt.Errorf("parsing response: %w", err)
	}

	if opts.JSON.IsJSON() {
		return cmdutil.PrintJSON(opts.IO.Out, labels)
	}

	w := tabwriter.NewWriter(opts.IO.Out, 0, 0, 2, ' ', 0)
	for _, l := range labels {
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", l.Name, l.Color, l.Description)
	}
	return w.Flush()
}
