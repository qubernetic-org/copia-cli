package list

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
)

var validJSONFields = []string{"number", "title", "state", "labels", "updatedAt"}

// ListOptions holds all inputs for the issue list command.
type ListOptions struct {
	IO         *iostreams.IOStreams
	HTTPClient *http.Client
	Host       string
	Token      string
	Owner      string
	Repo       string
	State      string
	Limit      int
	JSON       cmdutil.JSONFlags
}

type labelRef struct {
	Name string `json:"name"`
}

type issueEntry struct {
	Number    int64      `json:"number"`
	Title     string     `json:"title"`
	State     string     `json:"state"`
	UpdatedAt string     `json:"updated_at"`
	Labels    []labelRef `json:"labels"`
}

// NewCmdList creates the `copia issue list` command.
func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	opts := &ListOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List issues in a repository",
		Aliases: []string{"ls"},
		Example: `  copia issue list
  copia issue list --state closed
  copia issue list --json number,title,state`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.IO = f.IOStreams
			host, token, err := f.ResolveAuth()
			if err != nil {
				return err
			}
			opts.Host = host
			opts.Token = token

			if f.BaseRepo == nil {
				return fmt.Errorf("could not determine repository. Run from inside a git repository")
			}
			owner, repo, err := f.BaseRepo()
			if err != nil {
				return err
			}
			opts.Owner = owner
			opts.Repo = repo
			opts.HTTPClient = &http.Client{}
			return listRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.State, "state", "s", "open", "Filter by state: {open|closed|all}")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of issues")
	cmdutil.AddJSONFlags(cmd, &opts.JSON, validJSONFields)

	return cmd
}

func listRun(opts *ListOptions) error {
	url := fmt.Sprintf("https://%s/api/v1/repos/%s/%s/issues?state=%s&limit=%d&type=issues",
		opts.Host, opts.Owner, opts.Repo, opts.State, opts.Limit)

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

	var issues []issueEntry
	if err := json.Unmarshal(body, &issues); err != nil {
		return fmt.Errorf("parsing response: %w", err)
	}

	if opts.JSON.IsJSON() {
		return cmdutil.PrintJSON(opts.IO.Out, issues)
	}

	w := tabwriter.NewWriter(opts.IO.Out, 0, 0, 2, ' ', 0)
	for _, i := range issues {
		labels := ""
		for j, l := range i.Labels {
			if j > 0 {
				labels += ", "
			}
			labels += l.Name
		}
		_, _ = fmt.Fprintf(w, "#%d\t%s\t%s\t%s\n", i.Number, i.Title, i.State, labels)
	}
	return w.Flush()
}
