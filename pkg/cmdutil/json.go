package cmdutil

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// JSONFlags holds --json flag state.
type JSONFlags struct {
	Fields []string
}

// AddJSONFlags adds --json flag to a command.
func AddJSONFlags(cmd *cobra.Command, jf *JSONFlags, validFields []string) {
	cmd.Flags().StringSliceVar(&jf.Fields, "json", nil,
		fmt.Sprintf("Output JSON with selected fields: %v", validFields))
}

// IsJSON returns true if --json was specified.
func (jf *JSONFlags) IsJSON() bool {
	return jf.Fields != nil
}

// PrintJSON writes v as indented JSON to w.
func PrintJSON(w io.Writer, v interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
