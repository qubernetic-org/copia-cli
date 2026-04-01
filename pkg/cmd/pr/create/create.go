package create

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
)

// CreateOptions holds all inputs for the pr create command.
type CreateOptions struct {
	IO         *iostreams.IOStreams
	HTTPClient *http.Client
	Host       string
	Token      string
	Owner      string
	Repo       string
	Title      string
	Body       string
	Base       string
	Head       string
}

type createRequest struct {
	Title string `json:"title"`
	Body  string `json:"body,omitempty"`
	Base  string `json:"base"`
	Head  string `json:"head"`
}

type createResponse struct {
	Number  int64  `json:"number"`
	Title   string `json:"title"`
	HTMLURL string `json:"html_url"`
}

// NewCmdCreate creates the `copia pr create` command.
func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	opts := &CreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a pull request",
		Example: `  copia pr create --title "feat: add wrapper" --base main --head feature/wrapper
  copia pr create --title "fix: timeout" --base develop --head fix/timeout --body "Fixes #12"`,
		RunE: func(cmd *cobra.Command, args []string) error {
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
			return createRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Title, "title", "t", "", "PR title (required)")
	cmd.Flags().StringVarP(&opts.Body, "body", "b", "", "PR body")
	cmd.Flags().StringVar(&opts.Base, "base", "main", "Base branch")
	cmd.Flags().StringVarP(&opts.Head, "head", "H", "", "Head branch")

	return cmd
}

func createRun(opts *CreateOptions) error {
	if opts.Title == "" {
		return fmt.Errorf("title required")
	}

	payload := createRequest{
		Title: opts.Title,
		Body:  opts.Body,
		Base:  opts.Base,
		Head:  opts.Head,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://%s/api/v1/repos/%s/%s/pulls", opts.Host, opts.Owner, opts.Repo)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+opts.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := opts.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("connecting to %s: %w", opts.Host, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create PR (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result createResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(opts.IO.Out, "Created PR #%d: %s\n%s\n", result.Number, result.Title, result.HTMLURL)
	return nil
}
