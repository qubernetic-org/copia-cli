package cmdutil

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONFlags_IsJSON(t *testing.T) {
	jf := JSONFlags{}
	assert.False(t, jf.IsJSON())

	jf.Fields = []string{"name", "title"}
	assert.True(t, jf.IsJSON())
}

func TestPrintJSON(t *testing.T) {
	var buf bytes.Buffer
	data := map[string]string{"name": "test", "value": "123"}

	err := PrintJSON(&buf, data)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), `"name": "test"`)
	assert.Contains(t, buf.String(), `"value": "123"`)
}
