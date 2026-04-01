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

type ListOptions struct {
	IO         *iostreams.IOStreams
	HTTPClient *http.Client
	Host       string
	Token      string
	JSON       cmdutil.JSONFlags
}

type orgEntry struct {
	Username    string `json:"username"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	opts := &ListOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List your organizations",
		Aliases: []string{"ls"},
		Example: `  copia org list
  copia org list --json username,full_name`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.IO = f.IOStreams
			host, token, err := f.ResolveAuth()
			if err != nil {
				return err
			}
			opts.Host = host
			opts.Token = token
			opts.HTTPClient = &http.Client{}
			return listRun(opts)
		},
	}

	cmdutil.AddJSONFlags(cmd, &opts.JSON, []string{"username", "full_name", "description"})
	return cmd
}

func listRun(opts *ListOptions) error {
	url := fmt.Sprintf("https://%s/api/v1/user/orgs", opts.Host)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+opts.Token)

	resp, err := opts.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("connecting to %s: %w", opts.Host, err)
	}
	_ = resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error (HTTP %d)", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var orgs []orgEntry
	if err := json.Unmarshal(body, &orgs); err != nil {
		return err
	}

	if opts.JSON.IsJSON() {
		return cmdutil.PrintJSON(opts.IO.Out, orgs)
	}

	w := tabwriter.NewWriter(opts.IO.Out, 0, 0, 2, ' ', 0)
	for _, o := range orgs {
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", o.Username, o.FullName, o.Description)
	}
	return w.Flush()
}
