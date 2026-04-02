package repos

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

var validJSONFields = []string{"full_name", "description", "html_url"}

type SearchOptions struct {
	IO         *iostreams.IOStreams
	HTTPClient *http.Client
	Host       string
	Token      string
	Query      string
	Limit      int
	JSON       cmdutil.JSONFlags
}

type repoEntry struct {
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	HTMLURL     string `json:"html_url"`
}

type searchResponse struct {
	Data []repoEntry `json:"data"`
}

func NewCmdSearchRepos(f *cmdutil.Factory) *cobra.Command {
	opts := &SearchOptions{}

	cmd := &cobra.Command{
		Use:   "repos <query>",
		Short: "Search repositories",
		Example: `  copia search repos plc
  copia search repos "automation controller" --json fullName,description`,
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
			return SearchRun(opts)
		},
	}

	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of results")
	cmdutil.AddJSONFlags(cmd, &opts.JSON, validJSONFields)

	return cmd
}

func SearchRun(opts *SearchOptions) error {
	if err := cmdutil.ValidateLimit(opts.Limit); err != nil {
		return err
	}

	u := fmt.Sprintf("https://%s/api/v1/repos/search?q=%s&limit=%d",
		opts.Host, url.QueryEscape(opts.Query), opts.Limit)

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

	var result searchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if opts.JSON.IsJSON() {
		return cmdutil.PrintJSON(opts.IO.Out, result.Data)
	}

	w := tabwriter.NewWriter(opts.IO.Out, 0, 0, 2, ' ', 0)
	for _, r := range result.Data {
		_, _ = fmt.Fprintf(w, "%s\t%s\n", r.FullName, r.Description)
	}
	return w.Flush()
}
