package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
)

func AsDataURL(data []byte) string {
	return fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(data))
}

func SetAspect(image io.Reader, width, height int) (io.Reader, error) {
	i, err := png.Decode(image)
	if err != nil {
		return nil, fmt.Errorf("failed to parse as png: %w", err)
	}

	padded := padImage(i, width, height)

	var buf bytes.Buffer

	if err := png.Encode(&buf, padded); err != nil {
		return nil, fmt.Errorf("failed to encode as png: %w", err)
	}

	return &buf, nil
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
