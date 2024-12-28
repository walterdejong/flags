package flags

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegexShort(t *testing.T) {
	m := regexTag.FindStringSubmatch("-q")
	assert.Equal(t, m[regexTag.SubexpIndex("short")], "-q")
}

func TestRegexLong(t *testing.T) {
	m := regexTag.FindStringSubmatch("--quiet")
	assert.Equal(t, m[regexTag.SubexpIndex("long")], "--quiet")
}

func TestRegexShortAndLong(t *testing.T) {
	m := regexTag.FindStringSubmatch("-q, --quiet")
	assert.Equal(t, m[regexTag.SubexpIndex("short")], "-q")
	assert.Equal(t, m[regexTag.SubexpIndex("long")], "--quiet")
}

func TestRegexShortArg(t *testing.T) {
	m := regexTag.FindStringSubmatch("-n=NUM")
	assert.Equal(t, m[regexTag.SubexpIndex("short")], "-n")
	assert.Equal(t, m[regexTag.SubexpIndex("arg1")], "NUM")
}

func TestRegexLongArg(t *testing.T) {
	m := regexTag.FindStringSubmatch("--num=NUM")
	assert.Equal(t, m[regexTag.SubexpIndex("long")], "--num")
	assert.Equal(t, m[regexTag.SubexpIndex("arg2")], "NUM")
}

func TestRegexShortAndLongArg(t *testing.T) {
	m := regexTag.FindStringSubmatch("-n, --num=NUM")
	assert.Equal(t, m[regexTag.SubexpIndex("short")], "-n")
	assert.Equal(t, m[regexTag.SubexpIndex("arg1")], "")
	assert.Equal(t, m[regexTag.SubexpIndex("long")], "--num")
	assert.Equal(t, m[regexTag.SubexpIndex("arg2")], "NUM")
}

func TestRegexShortWithHelp(t *testing.T) {
	m := regexTag.FindStringSubmatch("-q help message")
	assert.Equal(t, m[regexTag.SubexpIndex("short")], "-q")
	assert.Equal(t, m[regexTag.SubexpIndex("help")], "help message")
}

func TestRegexLongWithHelp(t *testing.T) {
	m := regexTag.FindStringSubmatch("--quiet help message")
	assert.Equal(t, m[regexTag.SubexpIndex("long")], "--quiet")
	assert.Equal(t, m[regexTag.SubexpIndex("help")], "help message")
}

func TestRegexShortAndLongWithHelp(t *testing.T) {
	m := regexTag.FindStringSubmatch("-q, --quiet help message")
	assert.Equal(t, m[regexTag.SubexpIndex("short")], "-q")
	assert.Equal(t, m[regexTag.SubexpIndex("long")], "--quiet")
	assert.Equal(t, m[regexTag.SubexpIndex("help")], "help message")
}

func TestRegexShortArgWithHelp(t *testing.T) {
	m := regexTag.FindStringSubmatch("-n=NUM help message")
	assert.Equal(t, m[regexTag.SubexpIndex("short")], "-n")
	assert.Equal(t, m[regexTag.SubexpIndex("arg1")], "NUM")
	assert.Equal(t, m[regexTag.SubexpIndex("help")], "help message")
}

func TestRegexLongArgWithHelp(t *testing.T) {
	m := regexTag.FindStringSubmatch("--num=NUM help message")
	assert.Equal(t, m[regexTag.SubexpIndex("long")], "--num")
	assert.Equal(t, m[regexTag.SubexpIndex("arg2")], "NUM")
	assert.Equal(t, m[regexTag.SubexpIndex("help")], "help message")
}

func TestRegexShortAndLongArgWithHelp(t *testing.T) {
	m := regexTag.FindStringSubmatch("-n, --num=NUM help message")
	assert.Equal(t, m[regexTag.SubexpIndex("short")], "-n")
	assert.Equal(t, m[regexTag.SubexpIndex("arg1")], "")
	assert.Equal(t, m[regexTag.SubexpIndex("long")], "--num")
	assert.Equal(t, m[regexTag.SubexpIndex("arg2")], "NUM")
	assert.Equal(t, m[regexTag.SubexpIndex("help")], "help message")
}

func TestRegexError(t *testing.T) {
	m := regexTag.FindStringSubmatch("blurp")
	assert.Equal(t, len(m), 0)
}

// EOB
