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

var validJSONFields = []string{"number", "title", "state", "author", "base", "head", "updated_at"}

// ListOptions holds all inputs for the pr list command.
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

type branchRef struct {
	Label string `json:"label"`
}

type userRef struct {
	Login string `json:"login"`
}

type prEntry struct {
	Number    int64     `json:"number"`
	Title     string    `json:"title"`
	State     string    `json:"state"`
	User      userRef   `json:"user"`
	Base      branchRef `json:"base"`
	Head      branchRef `json:"head"`
	UpdatedAt string    `json:"updated_at"`
}

// NewCmdList creates the `copia pr list` command.
func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	opts := &ListOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List pull requests",
		Aliases: []string{"ls"},
		Example: `  copia pr list
  copia pr list --state closed
  copia pr list --json number,title,state`,
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
			return listRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.State, "state", "s", "open", "Filter by state: {open|closed|all}")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of pull requests")
	cmdutil.AddJSONFlags(cmd, &opts.JSON, validJSONFields)

	return cmd
}

func listRun(opts *ListOptions) error {
	url := fmt.Sprintf("https://%s/api/v1/repos/%s/%s/pulls?state=%s&limit=%d",
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

	var prs []prEntry
	if err := json.Unmarshal(body, &prs); err != nil {
		return err
	}

	if opts.JSON.IsJSON() {
		return cmdutil.PrintJSON(opts.IO.Out, prs)
	}

	w := tabwriter.NewWriter(opts.IO.Out, 0, 0, 2, ' ', 0)
	for _, pr := range prs {
		_, _ = fmt.Fprintf(w, "#%d\t%s\t%s\t%s\t%s <- %s\n",
			pr.Number, pr.Title, pr.State, pr.User.Login, pr.Base.Label, pr.Head.Label)
	}
	return w.Flush()
}
