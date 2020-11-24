package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/kbinani/screenshot"
)

func CaptureScreen() string {
	n := screenshot.NumActiveDisplays()
	files := ""
	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)

		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			panic(err)
		}
		fileName := fmt.Sprintf("%d_%dx%d.png", i, bounds.Dx(), bounds.Dy())
		file, _ := os.Create("./storage/" + fileName)
		defer file.Close()
		png.Encode(file, img)

		fmt.Printf("#%d : %v \"%s\"\n", i, bounds, fileName)
		// return filename
		files = fileName
	}
	return files
}
