package iostreams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTest_ReturnsWorkingStreams(t *testing.T) {
	ios, stdin, stdout, stderr := Test()

	stdin.WriteString("input data")
	assert.NotNil(t, ios)
	assert.Equal(t, "input data", stdin.String())
	assert.Equal(t, "", stdout.String())
	assert.Equal(t, "", stderr.String())
	assert.False(t, ios.IsStdoutTTY())
	assert.False(t, ios.IsStdinTTY())
}

func TestIOStreams_SetAndGetTTY(t *testing.T) {
	ios, _, _, _ := Test()

	ios.SetStdoutTTY(true)
	assert.True(t, ios.IsStdoutTTY())

	ios.SetStdoutTTY(false)
	assert.False(t, ios.IsStdoutTTY())

	ios.SetStdinTTY(true)
	assert.True(t, ios.IsStdinTTY())

	ios.SetStderrTTY(true)
	assert.True(t, ios.IsStderrTTY())
}

func TestTest_WriteToOut(t *testing.T) {
	ios, _, stdout, _ := Test()

	_, err := ios.Out.Write([]byte("hello"))
	assert.NoError(t, err)
	assert.Equal(t, "hello", stdout.String())
}

func TestTest_WriteToErrOut(t *testing.T) {
	ios, _, _, stderr := Test()

	_, err := ios.ErrOut.Write([]byte("error msg"))
	assert.NoError(t, err)
	assert.Equal(t, "error msg", stderr.String())
}
