//go:build !wasm
// +build !wasm

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

var (
	outFlag    = flag.String("out", "-", "the path to output the adjusted image to")
	ratioFlag  = flag.String("ratio", "16:9", "the aspect ratio to pad to, in the format h:w")
	ratioRegex = regexp.MustCompile(`^([1-9]\d?):([1-9]\d?)$`)
)

func init() {
	flag.StringVar(outFlag, "o", "-", "a shortcut for -out")
	flag.StringVar(ratioFlag, "r", "16:9", "a shortcut for -ratio")
}

func isTTY() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

func exitf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	if *outFlag == "-" && isTTY() {
		exitf("refusing to output binary to a TTY")
	}

	match := ratioRegex.FindStringSubmatch(*ratioFlag)
	if match == nil {
		exitf("invalid ratio '%s'\n", *ratioFlag)
	}

	width, _ := strconv.Atoi(match[1])
	height, _ := strconv.Atoi(match[2])

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		exitf("failed to open '%s': %+v", flag.Arg(0), err)
	}

	defer file.Close()

	padded, err := SetAspect(file, width, height)
	if err != nil {
		exitf("%+v", err)
	}

	output := os.Stdout
	if *outFlag != "-" {
		var err error
		output, err = os.Create(*outFlag)
		if err != nil {
			exitf("failed to create output file '%s': %+v", *outFlag, err)
		}
		defer output.Close()
	}

	if _, err := io.Copy(output, padded); err != nil {
		exitf("failed to write png data to '%s': %+v", *outFlag, err)
	}
}
