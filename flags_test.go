/*
	flags_test.go	WJ124
*/

package flags

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testOptions struct {
	Help     bool   `flags:"-h, --help"`
	Quiet    bool   `flags:"-q, --quiet             suppress output"`
	Verbose  int    `flags:"-v, --verbose           be more verbose (may be given multiple times)"`
	Num      int    `flags:"-n, --num=NUMBER        specify number"`
	Unsigned uint   `flags:"-u, --unsigned=NUMBER   specify number >= 0"`
	File     string `flags:"-f, --file=FILE         specify filename"`
}

func TestParse(t *testing.T) {
	opts := testOptions{}
	argv := []string{"prog", "-vvv", "--num=-1", "-u", "42", "--file", "hello", "foo", "bar", "-q"}

	args, err := Parse(argv, &opts)
	if err != nil {
		t.Errorf("Parse() error: %q", err)
	}

	assert.Equal(t, opts.Help, false)
	assert.Equal(t, opts.Quiet, true)
	assert.Equal(t, opts.Verbose, 3)
	assert.Equal(t, opts.Num, -1)
	assert.Equal(t, opts.Unsigned, uint(0x2a))
	assert.Equal(t, opts.File, "hello")
	assert.Equal(t, args, []string{"foo", "bar"});
}

func TestParseEmpty(t *testing.T) {
	opts := testOptions{}
	argv := []string{"prog"}

	args, err := Parse(argv, &opts)
	if err != nil {
		t.Errorf("Parse() error: %q", err)
	}

	assert.Equal(t, opts.Help, false)
	assert.Equal(t, opts.Quiet, false)
	assert.Equal(t, opts.Verbose, 0)
	assert.Equal(t, opts.Num, 0)
	assert.Equal(t, opts.Unsigned, uint(0))
	assert.Equal(t, opts.File, "")
	assert.Equal(t, args, []string{});
}

func TestParseUnknownOption(t *testing.T) {
	// passing an unknown option is an error

	opts := testOptions{}
	argv := []string{"prog", "-p"}

	_, err := Parse(argv, &opts)
	if err == nil {
		t.Fail()
	}
}

func TestParseUnknownLongOption(t *testing.T) {
	// passing an unknown long option is an error

	opts := testOptions{}
	argv := []string{"prog", "--foo"}

	_, err := Parse(argv, &opts)
	if err == nil {
		t.Fail()
	}
}

func TestParseRepeatBool(t *testing.T) {
	// repeating a boolean option is not an error

	opts := testOptions{}
	argv := []string{"prog", "-qqqq"}

	_, err := Parse(argv, &opts)
	if err != nil {
		t.Fail()
	}
}

func TestParseRepeatOptArg(t *testing.T) {
	// repeating an optarg is an error

	opts := testOptions{}
	argv := []string{"prog", "-n", "1", "--num=2"}

	_, err := Parse(argv, &opts)
	if err == nil {
		t.Fail()
	}
}

// EOB
