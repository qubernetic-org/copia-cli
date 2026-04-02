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

var validJSONFields = []string{"full_name", "description", "private", "updated_at"}

// ListOptions holds all inputs for the repo list command.
type ListOptions struct {
	IO         *iostreams.IOStreams
	HTTPClient *http.Client
	Host       string
	Token      string
	Org        string
	Limit      int
	JSON       cmdutil.JSONFlags
}

type repoEntry struct {
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	UpdatedAt   string `json:"updated_at"`
}

// NewCmdList creates the `copia repo list` command.
func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	opts := &ListOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List repositories",
		Aliases: []string{"ls"},
		Example: `  copia repo list
  copia repo list --org my-org
  copia repo list --json fullName,description`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.IO = f.IOStreams
			host, token, err := f.ResolveAuth()
			if err != nil {
				return err
			}
			opts.Host = host
			opts.Token = token
			opts.HTTPClient = &http.Client{}
			return ListRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Org, "org", "o", "", "List repositories for an organization")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of repositories")
	cmdutil.AddJSONFlags(cmd, &opts.JSON, validJSONFields)

	return cmd
}

func ListRun(opts *ListOptions) error {
	endpoint := "/api/v1/user/repos"
	if opts.Org != "" {
		endpoint = fmt.Sprintf("/api/v1/orgs/%s/repos", opts.Org)
	}

	url := fmt.Sprintf("https://%s%s?limit=%d", opts.Host, endpoint, opts.Limit)
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

	var repos []repoEntry
	if err := json.Unmarshal(body, &repos); err != nil {
		return fmt.Errorf("parsing response: %w", err)
	}

	if opts.JSON.IsJSON() {
		return cmdutil.PrintJSON(opts.IO.Out, repos)
	}

	w := tabwriter.NewWriter(opts.IO.Out, 0, 0, 2, ' ', 0)
	for _, r := range repos {
		visibility := "public"
		if r.Private {
			visibility = "private"
		}
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", r.FullName, r.Description, visibility)
	}
	return w.Flush()
}
