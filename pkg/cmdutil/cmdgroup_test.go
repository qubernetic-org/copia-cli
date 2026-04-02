package cmdutil

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestAddGroup(t *testing.T) {
	parent := &cobra.Command{Use: "issue"}
	create := &cobra.Command{Use: "create"}
	list := &cobra.Command{Use: "list"}
	close := &cobra.Command{Use: "close"}
	view := &cobra.Command{Use: "view"}

	AddGroup(parent, "General commands", create, list)
	AddGroup(parent, "Targeted commands", close, view)

	assert.Len(t, parent.Groups(), 2)
	assert.Equal(t, "General commands", parent.Groups()[0].Title)
	assert.Equal(t, "Targeted commands", parent.Groups()[1].Title)

	assert.Equal(t, "General commands", create.GroupID)
	assert.Equal(t, "General commands", list.GroupID)
	assert.Equal(t, "Targeted commands", close.GroupID)
	assert.Equal(t, "Targeted commands", view.GroupID)

	assert.Len(t, parent.Commands(), 4)
}
