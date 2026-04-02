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
	All        bool
	JSON       cmdutil.JSONFlags
}

type subject struct {
	Title string `json:"title"`
	Type  string `json:"type"`
}

type repoRef struct {
	FullName string `json:"full_name"`
}

type notification struct {
	ID         int64   `json:"id"`
	Subject    subject `json:"subject"`
	Repository repoRef `json:"repository"`
	Unread     bool    `json:"unread"`
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	opts := &ListOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List notifications",
		Long:    "List notifications for the authenticated user. By default, only unread notifications are shown.",
		Aliases: []string{"ls"},
		Example: "  $ copia-cli notification list",
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

	cmd.Flags().BoolVar(&opts.All, "all", false, "Show read and unread notifications")
	cmdutil.AddJSONFlags(cmd, &opts.JSON, []string{"id", "subject", "repository", "unread"})
	return cmd
}

func ListRun(opts *ListOptions) error {
	url := fmt.Sprintf("https://%s/api/v1/notifications?page=1", opts.Host)
	if opts.All {
		url += "&all=true"
	}
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

	var notifications []notification
	if err := json.Unmarshal(body, &notifications); err != nil {
		return err
	}

	if opts.JSON.IsJSON() {
		return cmdutil.PrintJSON(opts.IO.Out, notifications)
	}

	if len(notifications) == 0 {
		_, _ = fmt.Fprintln(opts.IO.Out, "No unread notifications")
		return nil
	}

	w := tabwriter.NewWriter(opts.IO.Out, 0, 0, 2, ' ', 0)
	for _, n := range notifications {
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", n.Repository.FullName, n.Subject.Type, n.Subject.Title)
	}
	return w.Flush()
}
