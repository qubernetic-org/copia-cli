package issues

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/qubernetic/copia-cli/pkg/cmdutil"
	"github.com/qubernetic/copia-cli/pkg/iostreams"
)

// SearchOptions holds all inputs for the search issues command.
type SearchOptions struct {
	IO         *iostreams.IOStreams
	HTTPClient *http.Client
	Host       string
	Token      string
	Owner      string
	Repo       string
	Query      string
	State      string
	Limit      int
	JSON       cmdutil.JSONFlags
}

type issueEntry struct {
	Number int64  `json:"number"`
	Title  string `json:"title"`
	State  string `json:"state"`
}

// NewCmdSearchIssues creates the `copia search issues` command.
func NewCmdSearchIssues(f *cmdutil.Factory) *cobra.Command {
	opts := &SearchOptions{}

	cmd := &cobra.Command{
		Use:   "issues <query>",
		Short: "Search issues in a repository",
		Long:  "Search issues within the current repository. Requires repo context (git remote or owner/repo argument).",
		Example: `  copia search issues "sensor timeout"
  copia search issues bug --state closed`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Query = args[0]
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
			return searchRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.State, "state", "s", "", "Filter by state: {open|closed}")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of results")
	cmdutil.AddJSONFlags(cmd, &opts.JSON, []string{"number", "title", "state"})

	return cmd
}

func searchRun(opts *SearchOptions) error {
	u := fmt.Sprintf("https://%s/api/v1/repos/%s/%s/issues?q=%s&limit=%d&type=issues",
		opts.Host, opts.Owner, opts.Repo, url.QueryEscape(opts.Query), opts.Limit)
	if opts.State != "" {
		u += "&state=" + url.QueryEscape(opts.State)
	}

	req, err := http.NewRequest("GET", u, nil)
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
		return err
	}

	if opts.JSON.IsJSON() {
		return cmdutil.PrintJSON(opts.IO.Out, issues)
	}

	w := tabwriter.NewWriter(opts.IO.Out, 0, 0, 2, ' ', 0)
	for _, i := range issues {
		_, _ = fmt.Fprintf(w, "#%d\t%s\t%s\n", i.Number, i.Title, i.State)
	}
	return w.Flush()
}
