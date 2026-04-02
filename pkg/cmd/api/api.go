package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/qubernetic/copia-cli/pkg/cmdutil"
	"github.com/qubernetic/copia-cli/pkg/iostreams"
)

// APIOptions holds all inputs for the api command.
type APIOptions struct {
	IO         *iostreams.IOStreams
	HTTPClient *http.Client
	Host       string
	Token      string
	Method     string
	Path       string
	Fields     []string
	Headers    []string
}

// NewCmdAPI creates the `copia api` command.
func NewCmdAPI(f *cmdutil.Factory) *cobra.Command {
	opts := &APIOptions{}

	cmd := &cobra.Command{
		Use:   "api <path>",
		Short: "Make an API request",
		Long:  "Make an authenticated HTTP request to the Copia REST API and print the response. The endpoint argument should be a path of a Gitea API v1 endpoint.",
		Example: `  # Get authenticated user
  $ copia-cli api /user

  # Create an issue
  $ copia-cli api -X POST /repos/my-org/my-repo/issues --field title="Bug report"

  # Delete a repo
  $ copia-cli api -X DELETE /repos/my-org/old-repo

  # Custom header
  $ copia-cli api /user --header "Accept: application/json"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Path = args[0]
			opts.IO = f.IOStreams

			host, token, err := f.ResolveAuth()
			if err != nil {
				return err
			}
			opts.Host = host
			opts.Token = token
			opts.HTTPClient = &http.Client{}
			return APIRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Method, "method", "X", "", "HTTP method (default: GET, or POST if --field is used)")
	cmd.Flags().StringSliceVarP(&opts.Fields, "field", "f", nil, "Add JSON body field (key=value)")
	cmd.Flags().StringSliceVarP(&opts.Headers, "header", "H", nil, "Add HTTP header (key: value)")

	return cmd
}

func APIRun(opts *APIOptions) error {
	if opts.Path == "" {
		return fmt.Errorf("path required")
	}

	// Normalize path
	path := opts.Path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if !strings.HasPrefix(path, "/api/v1") {
		path = "/api/v1" + path
	}

	// Determine method
	method := opts.Method
	if method == "" {
		if len(opts.Fields) > 0 {
			method = "POST"
		} else {
			method = "GET"
		}
	}

	// Build body from fields
	var body io.Reader
	if len(opts.Fields) > 0 {
		payload := map[string]string{}
		for _, f := range opts.Fields {
			parts := strings.SplitN(f, "=", 2)
			if len(parts) == 2 {
				payload[parts[0]] = parts[1]
			}
		}
		b, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		body = bytes.NewReader(b)
	}

	url := fmt.Sprintf("https://%s%s", opts.Host, path)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+opts.Token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Custom headers
	for _, h := range opts.Headers {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) == 2 {
			req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}

	resp, err := opts.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("connecting to %s: %w", opts.Host, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Pretty-print JSON if possible
	var prettyJSON bytes.Buffer
	if json.Indent(&prettyJSON, respBody, "", "  ") == nil {
		_, _ = fmt.Fprintln(opts.IO.Out, prettyJSON.String())
	} else {
		_, _ = fmt.Fprintln(opts.IO.Out, string(respBody))
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error (HTTP %d)", resp.StatusCode)
	}

	return nil
}
