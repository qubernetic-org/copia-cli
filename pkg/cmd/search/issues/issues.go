package issues

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
)

type SearchOptions struct {
	IO         *iostreams.IOStreams
	HTTPClient *http.Client
	Host       string
	Token      string
	Query      string
	State      string
	Limit      int
	JSON       cmdutil.JSONFlags
}

type repoRef struct {
	FullName string `json:"full_name"`
}

type issueEntry struct {
	Number     int64   `json:"number"`
	Title      string  `json:"title"`
	State      string  `json:"state"`
	Repository repoRef `json:"repository"`
}

func NewCmdSearchIssues(f *cmdutil.Factory) *cobra.Command {
	opts := &SearchOptions{}

	cmd := &cobra.Command{
		Use:   "issues <query>",
		Short: "Search issues",
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
			opts.HTTPClient = &http.Client{}
			return searchRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.State, "state", "s", "", "Filter by state: {open|closed}")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of results")
	cmdutil.AddJSONFlags(cmd, &opts.JSON, []string{"number", "title", "state", "repository"})

	return cmd
}

func searchRun(opts *SearchOptions) error {
	u := fmt.Sprintf("https://%s/api/v1/repos/search?q=%s&limit=%d&type=issues",
		opts.Host, url.QueryEscape(opts.Query), opts.Limit)
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
		_, _ = fmt.Fprintf(w, "%s#%d\t%s\t%s\n", i.Repository.FullName, i.Number, i.Title, i.State)
	}
	return w.Flush()
}
