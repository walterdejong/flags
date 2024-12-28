flags
=====

`flags` is a better command-line options parser for Go.

## Introduction

Go is great for programming command-line utilities. The standard `flag`
module, for parsing command-line arguments, is rather meh however.
This is my take at providing a better command-line parser for Go programs.

Command-line argument parsers fall in two categories:

* "getopt" style
* "builder pattern" style

The `getopt` style parsers require more coding, but they are very flexible
because any special behaviors can (and have to) be programmed in.

The `builder pattern` style parsers are mini-frameworks that require you
to learn and adhere to the framework. They are easy for trivial options
parsing, and not-so-great when you want special behavior.

This `flags` module is a `getopt` style parser with a twist: it uses Golang's
unique abilities of struct tagging and reflection to bind options, allowing
you to write a lot less code than usually with getopt-style parsers.
In a sense, the "builder" is builtin to the `flags` module.

## Usage

Step one: define the command-line options in a `struct`. This struct type
may be named anything you like; this example uses typename `Options`.
The struct fields are tagged with a special `flags` tag, that defines both
the behavior and the long help message of the command-line option.

```go
    import "github.com/walterdejong/flags"
    
    type Options struct {
        Help     bool   `flags:"-h, --help"`
        Quiet    bool   `flags:"-q, --quiet             suppress output"`
        Verbose  int    `flags:"-v, --verbose           be more verbose (may be given multiple times)"`
        Num      int    `flags:"-n, --num=NUMBER        specify number"`
        Unsigned uint   `flags:"-u, --unsigned=NUMBER   specify number >= 0"`
        File     string `flags:"-f, --file=FILE         specify filename"`
    }
```

Notes for this particular example:
* the `help` option has no long help message; therefore this option will be
  omitted when printing the long help information
* the `quiet` option is a switch that can be true or false
* the `verbose` option is an integer counter
* the other options require an argument

Step two: parse the command-line arguments!

The `flags.Parse()` function takes a command-line, fills in the given options
structure, and returns any (left-over) arguments.

```go
    opts := Options{}
    args, err := flags.Parse(os.Args, &opts)
    if err != nil {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(2)
    }
```

Step three: program logic; handling options and arguments.

```go
    fmt.Printf("opts == %#v\n", opts)
    fmt.Printf("args == %#v\n", args)
    
    if opts.Help {
        fmt.Println("usage: example [options] [args ...]")
        flags.PrintHelp(&opts)
        os.Exit(1)
    }
```

The `flags.PrintHelp()` function prints long help information.
Note that it does not print any short usage information, program description,
or copyright line, and this is entirely by design.

Finally, we can run the example program:

```
    $ ./example -vvv --num=10 -f hello foo bar -q
    opts == main.Options{Help:false, Quiet:true, Verbose:3, Num:10, Unsigned:0x0, File:"hello"}
    args == []string{"foo", "bar"}
```

The parser is a GNU-style parser that (unlike traditional getopt) allows
options to occur after arguments.

## Copyright and License

Copyright 2024 by Walter de Jong <walter@heiho.net>

This software is provided under terms described in the MIT license.
