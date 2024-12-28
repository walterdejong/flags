/*
	flags_test.go	WJ124

	Copyright (c) 2024 Walter de Jong <walter@heiho.net>

	Permission is hereby granted, free of charge, to any person obtaining a copy of
	this software and associated documentation files (the "Software"), to deal in
	the Software without restriction, including without limitation the rights to
	use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
	of the Software, and to permit persons to whom the Software is furnished to do
	so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all
	copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
	SOFTWARE.
*/

package flags_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/walterdejong/flags"
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

	args, err := flags.Parse(argv, &opts)
	if err != nil {
		t.Errorf("Parse() error: %q", err)
	}

	assert.Equal(t, opts.Help, false)
	assert.Equal(t, opts.Quiet, true)
	assert.Equal(t, opts.Verbose, 3)
	assert.Equal(t, opts.Num, -1)
	assert.Equal(t, opts.Unsigned, uint(0x2a))
	assert.Equal(t, opts.File, "hello")
	assert.Equal(t, args, []string{"foo", "bar"})
}

func TestParseEmpty(t *testing.T) {
	opts := testOptions{}
	argv := []string{"prog"}

	args, err := flags.Parse(argv, &opts)
	if err != nil {
		t.Errorf("Parse() error: %q", err)
	}

	assert.Equal(t, opts.Help, false)
	assert.Equal(t, opts.Quiet, false)
	assert.Equal(t, opts.Verbose, 0)
	assert.Equal(t, opts.Num, 0)
	assert.Equal(t, opts.Unsigned, uint(0))
	assert.Equal(t, opts.File, "")
	assert.Equal(t, args, []string{})
}

func TestParseUnknownOption(t *testing.T) {
	// passing an unknown option is an error

	opts := testOptions{}
	argv := []string{"prog", "-p"}

	_, err := flags.Parse(argv, &opts)
	if err == nil {
		t.Fail()
	}
}

func TestParseUnknownLongOption(t *testing.T) {
	// passing an unknown long option is an error

	opts := testOptions{}
	argv := []string{"prog", "--foo"}

	_, err := flags.Parse(argv, &opts)
	if err == nil {
		t.Fail()
	}
}

func TestParseRepeatBool(t *testing.T) {
	// repeating a boolean option is not an error

	opts := testOptions{}
	argv := []string{"prog", "-qqqq"}

	_, err := flags.Parse(argv, &opts)
	if err != nil {
		t.Fail()
	}
}

func TestParseRepeatOptArg(t *testing.T) {
	// repeating an optarg is an error

	opts := testOptions{}
	argv := []string{"prog", "-n", "1", "--num=2"}

	_, err := flags.Parse(argv, &opts)
	if err == nil {
		t.Fail()
	}
}

// EOB
