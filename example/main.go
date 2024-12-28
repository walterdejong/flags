/*
	example program for `flags` module
*/

package main

import (
	"fmt"
	"github.com/walterdejong/flags"
	"os"
)

type Options struct {
	Help     bool   `flags:"-h, --help"`
	Quiet    bool   `flags:"-q, --quiet             suppress output"`
	Verbose  int    `flags:"-v, --verbose           be more verbose (may be given multiple times)"`
	Num      int    `flags:"-n, --num=NUMBER        specify number"`
	Unsigned uint   `flags:"-u, --unsigned=NUMBER   specify number >= 0"`
	File     string `flags:"-f, --file=FILE         specify filename"`
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("usage: example [options] [args ...]")
		os.Exit(1)
	}

	opts := Options{}
	args, err := flags.Parse(os.Args, &opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	fmt.Printf("opts == %#v\n", opts)
	fmt.Printf("args == %#v\n", args)

	if opts.Help {
		fmt.Println("usage: example [options] [args ...]")
		flags.PrintHelp(&opts)
		os.Exit(1)
	}
}

// EOB
