package colorfmt_test

import (
	"bytes"
	"testing"

	"github.com/blazejsewera/go-test-proxy/colorfmt"
	"github.com/blazejsewera/go-test-proxy/test/assert"
)

func TestNoColorOutput(t *testing.T) {
	t.Run("creating a new colorfmt does not write anything to stdout and stderr", func(t *testing.T) {
		stdout := new(bytes.Buffer)
		stderr := new(bytes.Buffer)

		colorfmt.New(false, stdout, stderr)

		assert.Empty(t, stdout.Bytes())
		assert.Empty(t, stderr.Bytes())
	})

	t.Run("print functions write unstyled strings to stdout and stderr", func(t *testing.T) {
		stdout := new(bytes.Buffer)
		stderr := new(bytes.Buffer)
		cfmt := colorfmt.New(false, stdout, stderr)

		cfmt.Cprint(colorfmt.Bold, colorfmt.Red, "test")

		assert.Equal(t, "test", stdout.String())

		stdout.Reset()
		stderr.Reset()
		cfmt.Cprintf(colorfmt.Bold, colorfmt.Red, "test %d", 1)

		assert.Equal(t, "test 1", stdout.String())

		stdout.Reset()
		stderr.Reset()
		cfmt.Cerrprintf(colorfmt.Bold, colorfmt.Red, "test %d", 2)

		assert.Equal(t, "test 2", stderr.String())
	})
}

func TestColorOutput(t *testing.T) {
	t.Run("creating a new colorfmt writes base style to stdout and stderr", func(t *testing.T) {
		stdout := new(bytes.Buffer)
		stderr := new(bytes.Buffer)

		colorfmt.New(true, stdout, stderr)

		assert.EqualQuoted(t, "\x1b[0;37m", stdout.String())
		assert.EqualQuoted(t, "\x1b[0;37m", stderr.String())
	})

	t.Run("print functions write styled strings to stdout and stderr", func(t *testing.T) {
		stdout := new(bytes.Buffer)
		stderr := new(bytes.Buffer)
		cfmt := colorfmt.New(true, stdout, stderr)
		stdout.Reset()
		stderr.Reset()

		cfmt.Cprint(colorfmt.Bold, colorfmt.Red, "test")

		assert.EqualQuoted(t, "\x1b[1;31mtest\x1b[0m", stdout.String())

		stdout.Reset()
		stderr.Reset()
		cfmt.Cprintf(colorfmt.Bold, colorfmt.Red, "test %d", 1)

		assert.EqualQuoted(t, "\x1b[1;31mtest 1\x1b[0m", stdout.String())

		stdout.Reset()
		stderr.Reset()
		cfmt.Cerrprintf(colorfmt.Bold, colorfmt.Red, "test %d", 2)

		assert.EqualQuoted(t, "\x1b[1;31mtest 2\x1b[0m", stderr.String())
	})
}
