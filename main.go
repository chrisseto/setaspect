package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"regexp"
	"strconv"
)

var (
	outFlag    = flag.String("out", "-", "the path to output the adjusted image to. defaults to stdout")
	ratioFlag  = flag.String("ratio", "16:9", "the aspect ratio to pad to, in the format h:w. default 16:9")
	ratioRegex = regexp.MustCompile(`^([1-9]\d?):([1-9]\d?)$`)
)

func init() {
	flag.StringVar(outFlag, "o", "-", "a shortcut for -out")
	flag.StringVar(ratioFlag, "r", "16:9", "a shortcut for -ratio")
}

func padImage(i image.Image, width, height int) image.Image {
	x := i.Bounds().Dx()
	y := i.Bounds().Dy()
	offset := image.Pt(0, 0)

	haveAspect := float32(x) / float32(y)
	wantAspect := float32(width) / float32(height)

	if haveAspect < wantAspect {
		x = int(wantAspect * float32(y))
		offset.X = (x - i.Bounds().Dx()) / 2
	} else {
		y = int(float32(x) / (float32(width) / float32(height)))
		offset.Y = (y - i.Bounds().Dy()) / 2
	}

	padded := image.NewRGBA(image.Rect(0, 0, x, y))

	draw.Draw(padded, i.Bounds().Add(offset), i, image.Pt(0, 0), draw.Src)

	return padded
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

	if flag.ErrHelp != nil || len(flag.Args()) != 1 {
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

	i, err := png.Decode(file)
	if err != nil {
		exitf("failed to parse '%s' as png: %+v", flag.Arg(0), err)
	}

	padded := padImage(i, width, height)

	output := os.Stdout
	if *outFlag != "-" {
		var err error
		output, err = os.Create(*outFlag)
		if err != nil {
			exitf("failed to create output file '%s': %+v", *outFlag, err)
		}
		defer output.Close()
	}

	if err := png.Encode(output, padded); err != nil {
		exitf("failed to write png data to '%s': %+v", *outFlag, err)
	}
}
