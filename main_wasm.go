//go:build wasm

package main

import (
	"bytes"
	"io"
	"syscall/js"
)

func getElementById(id string) js.Value {
	return js.Global().Get("document").Call("getElementById", id)
}

func adapter(this js.Value, args []js.Value) interface{} {
	file := args[0].Get("target").Get("files").Index(0)

	reader := js.Global().Get("FileReader").New()
	reader.Set("onload", onFileLoad)

	reader.Call("readAsArrayBuffer", file)

	return nil
}

var onFileLoad = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	js.Global().Set("my_result", this.Get("result"))
	buffer := js.Global().Get("Uint8Array").New(this.Get("result"))

	imageBytes := make([]byte, int(buffer.Get("byteLength").Float()))
	js.CopyBytesToGo(imageBytes, buffer)

	padded, err := SetAspect(bytes.NewReader(imageBytes), 16, 9)
	if err != nil {
		// Easiest way to handle errors :s
		panic(err)
	}

	// Ignore the error, we're reading from a in memory buffer.
	paddedBytes, _ := io.ReadAll(padded)

	output := getElementById("output")
	output.Set("src", AsDataURL(paddedBytes))

	return nil
})

func main() {
	getElementById("input").Call("addEventListener", "change", js.FuncOf(adapter))

	// Block forever, not sure why though?
	c := make(chan struct{}, 0)
	<-c
}
