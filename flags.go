/*
	flags.go	WJ124

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

// Package flags provides a command-line arguments parser.
package flags

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type optionDef struct {
	shortOpt   string
	shortArg   string
	longOpt    string
	longArg    string
	requireArg bool
	help       string
	fieldIndex int
	kind       reflect.Kind
}

// this matches "-a=ARG1, --long-opt=ARG2    help text"
// or variants thereof
var regexTag = regexp.MustCompile(`^(?:` +
	`(?:(?P<short>-[a-zA-Z0-9-])(?:=(?P<arg1>[^,\s]+))?)?` +
	`(?:\s*,\s*)?` +
	`)?` +
	`(?:(?P<long>--[a-zA-Z0-9][a-zA-Z0-9-]+)(?:=(?P<arg2>[^,\s]+))?)?` +
	`(?:[,\s]+(?P<help>[^-\s].+)?)?$`)

// Parse the given command-line arguments.
// Returns any remaining arguments, and err value.
// Parameter `argv` is the command line, like `os.Args`.
// Parameter `taggedStruct` must be a pointer to struct that contains
// fields that are tagged with "flags".
// For example:
//
//	type Options struct {
//	    Help     bool   `flags:"-h, --help"`
//	    Quiet    bool   `flags:"-q, --quiet             suppress output"`
//	    Verbose  int    `flags:"-v, --verbose           be more verbose (may be given multiple times)"`
//	    Num      int    `flags:"-n, --num=NUMBER        specify number"`
//	    Unsigned uint   `flags:"-u, --unsigned=NUMBER   specify number >= 0"`
//	    File     string `flags:"-f, --file=FILE         specify filename"`
//	}
func Parse(argv []string, taggedStructP any) ([]string, error) {
	// make a list of defined options from the tagged struct
	definedOptions := defineOptions(taggedStructP)

	// optionsMap maps option string "--foo" to index into definedOptions
	optionsMap := makeOptionsMap(definedOptions)

	// argAlready[opt] is true when argument (already) has been set
	argAlready := make(map[string]bool)

	s := reflect.ValueOf(taggedStructP).Elem()

	args := make([]string, 0)
	onlyArgsRemain := false
	expectOptArg := false
	expectOptIdx := 0

	for _, arg := range argv[1:] {
		if expectOptArg {
			optarg := arg
			optDef := &definedOptions[expectOptIdx]
			opt := optDef.longOpt
			if opt == "" {
				opt = optDef.shortOpt
			}

			combined := combineOptions(optDef.shortOpt, optDef.longOpt)
			if argAlready[combined] {
				return args, errors.New(fmt.Sprintf("option %s was passed multiple times", combined))
			}
			err := setOptArg(s.Field(optDef.fieldIndex), optarg)
			if err != nil {
				return args, errors.New(fmt.Sprintf("option %s: invalid value %q", opt, optarg))
			}
			argAlready[combined] = true

			expectOptArg = false
			continue
		}

		if onlyArgsRemain {
			args = append(args, arg)
			continue
		}
		if arg == "--" {
			onlyArgsRemain = true
			continue
		}

		if arg == "" {
			args = append(args, arg)
			continue
		}

		if arg[0] == '-' {
			// it is an option (either short or long)

			arr := strings.SplitN(arg, "=", 2)
			if len(arr) == 2 {
				opt := arr[0]
				optarg := arr[1]

				idx, ok := optionsMap[opt]
				if !ok {
					return args, errors.New(fmt.Sprintf("invalid option: %q", opt))
				}
				optDef := &definedOptions[idx]
				if !optDef.requireArg {
					return args, errors.New(fmt.Sprintf("option %s does not take an argument", opt))
				}
				combined := combineOptions(optDef.shortOpt, optDef.longOpt)
				if argAlready[combined] {
					return args, errors.New(fmt.Sprintf("option %s was passed multiple times", opt))
				}
				err := setOptArg(s.Field(optDef.fieldIndex), optarg)
				if err != nil {
					return args, errors.New(fmt.Sprintf("option %s: invalid value %q", opt, optarg))
				}
				argAlready[combined] = true
				continue
			}

			// 'bare' option without argument
			if len(arg) > 2 && arg[1] != '-' {
				// multiple short options given; split them out
				for k, c := range arg[1:] {
					opt := "-" + string(c)
					idx, ok := optionsMap[opt]
					if !ok {
						return args, errors.New(fmt.Sprintf("invalid option: %q (part of %q)", opt, arg))
					}
					optDef := &definedOptions[idx]
					if optDef.requireArg {
						// this is fine if it is the last character in arg
						// else consider it a syntax error
						if k >= len(arg[1:])-1 {
							expectOptArg = true
							expectOptIdx = idx
							break
						}
						return args, errors.New(fmt.Sprintf("option %s requires an argument", opt))
					}
					setOpt(s.Field(optDef.fieldIndex))
				}
				continue
			}

			opt := arg
			idx, ok := optionsMap[opt]
			if !ok {
				return args, errors.New(fmt.Sprintf("invalid option: %q", opt))
			}
			optDef := &definedOptions[idx]
			if optDef.requireArg {
				expectOptArg = true
				expectOptIdx = idx
				continue
			}
			setOpt(s.Field(optDef.fieldIndex))
		} else {
			// it is an argument
			args = append(args, arg)
		}
	}

	if expectOptArg {
		optDef := &definedOptions[expectOptIdx]
		opt := optDef.longOpt
		if opt == "" {
			opt = optDef.shortOpt
		}
		optarg := optDef.longArg
		if optarg == "" {
			optarg = optDef.shortArg
		}
		return args, errors.New(fmt.Sprintf("option %s requires an argument %s", opt, optarg))
	}
	return args, nil
}

func defineOptions(taggedStructP any) []optionDef {
	// Analyzes the tagged struct and converts it into a list of optionDef structs

	definedOptions := make([]optionDef, 0)

	s := reflect.ValueOf(taggedStructP).Elem()
	typeOfT := s.Type()
	for i := 0; i < typeOfT.NumField(); i++ {
		field := typeOfT.Field(i)
		tag := field.Tag.Get("flags")
		if tag == "" {
			continue
		}
		kind := field.Type.Kind()
		// typecheck the field
		switch kind {
		case reflect.Bool, reflect.Int, reflect.Uint, reflect.String:
			// pass
		default:
			panic(fmt.Sprintf("unable to handle type %v (not implemented)", kind))
		}

		// parse the tag; get the option names, and
		// find out whether the option takes an argument
		// we make an optionDef structure from the tag

		matches := regexTag.FindStringSubmatch(tag)
		if len(matches) == 0 {
			panic(fmt.Sprintf("syntax error in tag flags:%q", tag))
		}
		regexNames := regexTag.SubexpNames()
		m := make(map[string]string)
		for n, svalue := range matches {
			m[regexNames[n]] = svalue
		}
		opt := optionDef{
			shortOpt:   m["short"],
			shortArg:   m["arg1"],
			longOpt:    m["long"],
			longArg:    m["arg2"],
			requireArg: false,
			help:       m["help"],
			fieldIndex: i,
			kind:       kind,
		}
		opt.requireArg = opt.shortArg != "" || opt.longArg != ""

		if opt.shortOpt == "" && opt.longOpt == "" {
			panic(fmt.Sprintf("syntax error in tag flags:%q", tag))
		}

		definedOptions = append(definedOptions, opt)
	}
	return definedOptions
}

func makeOptionsMap(definedOptions []optionDef) map[string]int {
	// Returns map that maps option string "--foo" to index into definedOptions

	optionsMap := make(map[string]int)
	for i := range definedOptions {
		x := &definedOptions[i]
		if x.shortOpt != "" {
			optionsMap[x.shortOpt] = i
		}
		if x.longOpt != "" {
			optionsMap[x.longOpt] = i
		}
	}
	return optionsMap
}

func setOpt(field reflect.Value) {
	switch field.Kind() {
	case reflect.Bool:
		field.SetBool(true)
	case reflect.Int:
		field.SetInt(field.Int() + 1)
	case reflect.Uint:
		field.SetUint(field.Uint() + 1)
	default:
		panic(fmt.Sprintf("setOpt(): unsupported type %v", field.Type()))
	}
}

func setOptArg(field reflect.Value, optarg string) error {
	switch field.Kind() {
	case reflect.Bool:
		b, err := strconv.ParseBool(optarg)
		if err != nil {
			return err
		}
		field.SetBool(b)
	case reflect.Int:
		i, err := strconv.ParseInt(optarg, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(i)
	case reflect.Uint:
		i, err := strconv.ParseUint(optarg, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(i)
	case reflect.String:
		field.SetString(optarg)
	default:
		panic(fmt.Sprintf("setOptArg(): unsupported type %v", field.Type()))
	}
	return nil
}

func combineOptions(short string, long string) string {
	if short != "" && long != "" {
		return short + ", " + long
	}
	if short != "" {
		return short
	}
	return long
}

// PrintHelp prints long usage information.
// Note: this is a helper function that only prints the options;
// it does not produce a short usage line, program description,
// copyright line, and so on.
func PrintHelp(taggedStructP any) {
	// make a list of defined options from the tagged struct
	definedOptions := defineOptions(taggedStructP)

	const ColumnWidth = 28

	var opt, optarg string

	for i := range definedOptions {
		optDef := &definedOptions[i]

		if optDef.help == "" {
			// no help text: do not show
			continue
		}

		opt = combineOptions(optDef.shortOpt, optDef.longOpt)

		optarg = ""
		if optDef.shortArg != "" {
			optarg = optDef.shortArg
		}
		if optDef.longArg != "" {
			optarg = optDef.longArg
		}
		if optarg != "" {
			opt += "=" + optarg
		}

		if len(opt) >= ColumnWidth {
			// print on two lines
			fmt.Printf("  %s\n", opt)
			opt = " "
		}
		fmt.Printf("  %-*s  %s\n", ColumnWidth, opt, optDef.help)
	}
}

// EOB
