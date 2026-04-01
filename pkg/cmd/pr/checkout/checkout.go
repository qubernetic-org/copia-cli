package checkout

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
)

type CheckoutOptions struct {
	IO     *iostreams.IOStreams
	Number int64
}

func NewCmdCheckout(f *cmdutil.Factory) *cobra.Command {
	opts := &CheckoutOptions{}

	cmd := &cobra.Command{
		Use:   "checkout <number>",
		Short: "Check out a pull request locally",
		Example: `  copia pr checkout 7`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			num, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid PR number: %s", args[0])
			}
			opts.Number = num
			opts.IO = f.IOStreams
			return checkoutRun(opts)
		},
	}

	return cmd
}

func checkoutRun(opts *CheckoutOptions) error {
	refSpec := buildRefSpec(opts.Number)
	branchName := fmt.Sprintf("pr-%d", opts.Number)

	fetchCmd := exec.Command("git", "fetch", "origin", refSpec)
	fetchCmd.Stdout = opts.IO.Out
	fetchCmd.Stderr = opts.IO.ErrOut
	fetchCmd.Stdin = os.Stdin

	if err := fetchCmd.Run(); err != nil {
		return fmt.Errorf("fetching PR #%d: %w", opts.Number, err)
	}

	checkoutCmd := exec.Command("git", "checkout", branchName)
	checkoutCmd.Stdout = opts.IO.Out
	checkoutCmd.Stderr = opts.IO.ErrOut
	checkoutCmd.Stdin = os.Stdin

	if err := checkoutCmd.Run(); err != nil {
		return fmt.Errorf("checking out %s: %w", branchName, err)
	}

	_, _ = fmt.Fprintf(opts.IO.Out, "Checked out PR #%d as %s\n", opts.Number, branchName)
	return nil
}

func buildRefSpec(number int64) string {
	return fmt.Sprintf("pull/%d/head:pr-%d", number, number)
}
