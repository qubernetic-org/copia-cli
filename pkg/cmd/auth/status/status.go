package status

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/qubernetic/copia-cli/internal/config"
	"github.com/qubernetic/copia-cli/pkg/cmdutil"
	"github.com/qubernetic/copia-cli/pkg/iostreams"
)

// StatusOptions holds all inputs for the status command.
type StatusOptions struct {
	IO         *iostreams.IOStreams
	ConfigPath string
	HTTPClient *http.Client
}

// NewCmdStatus creates the `copia auth status` command.
func NewCmdStatus(f *cmdutil.Factory) *cobra.Command {
	opts := &StatusOptions{}

	cmd := &cobra.Command{
		Use:     "status",
		Short:   "View authentication status",
		Example: "  copia auth status",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.IO = f.IOStreams
			opts.ConfigPath = config.DefaultPath()
			opts.HTTPClient = &http.Client{}
			return StatusRun(opts)
		},
	}

	return cmd
}

func StatusRun(opts *StatusOptions) error {
	cfg, err := config.Load(opts.ConfigPath)
	if err != nil {
		return err
	}

	if len(cfg.Hosts) == 0 {
		return fmt.Errorf("not logged in to any Copia instance. Run 'copia auth login'")
	}

	for host, hc := range cfg.Hosts {
		tokenStatus := "Token valid"
		if opts.HTTPClient != nil {
			url := fmt.Sprintf("https://%s/api/v1/user", host)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				tokenStatus = fmt.Sprintf("Error: %v", err)
			} else {
				req.Header.Set("Authorization", "token "+hc.Token)
				resp, err := opts.HTTPClient.Do(req)
				if err != nil {
					tokenStatus = fmt.Sprintf("Error: %v", err)
				} else {
					defer func() { _ = resp.Body.Close() }()
					if resp.StatusCode != http.StatusOK {
						tokenStatus = fmt.Sprintf("Token invalid (HTTP %d)", resp.StatusCode)
					}
				}
			}
		}

		_, _ = fmt.Fprintf(opts.IO.Out, "%s\n  User: %s\n  %s\n", host, hc.User, tokenStatus)
	}

	return nil
}
