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

var validJSONFields = []string{"tag_name", "name", "draft", "prerelease", "published_at"}

type ListOptions struct {
	IO         *iostreams.IOStreams
	HTTPClient *http.Client
	Host       string
	Token      string
	Owner      string
	Repo       string
	Limit      int
	JSON       cmdutil.JSONFlags
}

type releaseEntry struct {
	ID          int64  `json:"id"`
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	Draft       bool   `json:"draft"`
	Prerelease  bool   `json:"prerelease"`
	PublishedAt string `json:"published_at"`
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	opts := &ListOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List releases",
		Long:    "List releases in a Copia repository, ordered by creation date.",
		Aliases: []string{"ls"},
		Example: `  $ copia-cli release list
  $ copia-cli release list --json tagName,name`,
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

	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of releases")
	cmdutil.AddJSONFlags(cmd, &opts.JSON, validJSONFields)

	return cmd
}

func ListRun(opts *ListOptions) error {
	if err := cmdutil.ValidateLimit(opts.Limit); err != nil {
		return err
	}

	url := fmt.Sprintf("https://%s/api/v1/repos/%s/%s/releases?limit=%d",
		opts.Host, opts.Owner, opts.Repo, opts.Limit)

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

	var releases []releaseEntry
	if err := json.Unmarshal(body, &releases); err != nil {
		return fmt.Errorf("parsing response: %w", err)
	}

	if opts.JSON.IsJSON() {
		return cmdutil.PrintJSON(opts.IO.Out, releases)
	}

	w := tabwriter.NewWriter(opts.IO.Out, 0, 0, 2, ' ', 0)
	for _, r := range releases {
		status := "release"
		if r.Draft {
			status = "draft"
		} else if r.Prerelease {
			status = "pre-release"
		}
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", r.TagName, r.Name, status)
	}
	return w.Flush()
}
