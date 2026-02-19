// Command genicon generates assets/icon.ico — a 32x32 "VB" icon for the system tray.
//
// Run: go run tools/genicon/main.go
package main

import (
	"encoding/binary"
	"image"
	"image/color"
	"os"
)

func main() {
	const size = 32
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	bg := color.RGBA{R: 41, G: 128, B: 185, A: 255} // #2980B9
	fg := color.RGBA{R: 255, G: 255, B: 255, A: 255} // white

	// Fill background.
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			img.Set(x, y, bg)
		}
	}

	// Draw "V" (columns 3-13, rows 6-25).
	drawLetter(img, letterV, 3, 6, fg)
	// Draw "B" (columns 18-28, rows 6-25).
	drawLetter(img, letterB, 18, 6, fg)

	os.MkdirAll("assets", 0755)
	writeICO("assets/icon.ico", img)
}

// Each letter is defined as a 11x20 bitmap (string rows, '#' = pixel).
var letterV = []string{
	"##.......##",
	"##.......##",
	"##.......##",
	".##.....##.",
	".##.....##.",
	".##.....##.",
	"..##...##..",
	"..##...##..",
	"..##...##..",
	"...##.##...",
	"...##.##...",
	"...##.##...",
	"....###....",
	"....###....",
	"....###....",
	".....#.....",
	".....#.....",
	".....#.....",
	".....#.....",
	".....#.....",
}

var letterB = []string{
	"########...",
	"########...",
	"##.....##..",
	"##.....##..",
	"##.....##..",
	"########...",
	"########...",
	"########...",
	"##......##.",
	"##......##.",
	"##......##.",
	"##......##.",
	"##......##.",
	"##......##.",
	"##......##.",
	"##......##.",
	"##.....##..",
	"##.....##..",
	"########...",
	"########...",
}

func drawLetter(img *image.RGBA, letter []string, ox, oy int, c color.Color) {
	for row, line := range letter {
		for col, ch := range line {
			if ch == '#' {
				img.Set(ox+col, oy+row, c)
			}
		}
	}
}

func writeICO(path string, img *image.RGBA) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	pixelSize := w * h * 4
	andMaskRowSize := ((w + 31) / 32) * 4
	andMaskSize := andMaskRowSize * h
	imageDataSize := 40 + pixelSize + andMaskSize

	// ICONDIR
	binary.Write(f, binary.LittleEndian, uint16(0)) // reserved
	binary.Write(f, binary.LittleEndian, uint16(1)) // type = ICO
	binary.Write(f, binary.LittleEndian, uint16(1)) // count = 1

	// ICONDIRENTRY
	bw := byte(w)
	if w >= 256 {
		bw = 0
	}
	bh := byte(h)
	if h >= 256 {
		bh = 0
	}
	f.Write([]byte{bw, bh, 0, 0})                                    // width, height, palette, reserved
	binary.Write(f, binary.LittleEndian, uint16(1))                   // color planes
	binary.Write(f, binary.LittleEndian, uint16(32))                  // bits per pixel
	binary.Write(f, binary.LittleEndian, uint32(imageDataSize))       // image data size
	binary.Write(f, binary.LittleEndian, uint32(6+16))                // offset to image data

	// BITMAPINFOHEADER
	binary.Write(f, binary.LittleEndian, uint32(40))   // header size
	binary.Write(f, binary.LittleEndian, int32(w))     // width
	binary.Write(f, binary.LittleEndian, int32(h*2))   // height (doubled for ICO)
	binary.Write(f, binary.LittleEndian, uint16(1))    // planes
	binary.Write(f, binary.LittleEndian, uint16(32))   // bit count
	binary.Write(f, binary.LittleEndian, uint32(0))    // compression
	binary.Write(f, binary.LittleEndian, uint32(0))    // image size (can be 0)
	binary.Write(f, binary.LittleEndian, int32(0))     // x ppm
	binary.Write(f, binary.LittleEndian, int32(0))     // y ppm
	binary.Write(f, binary.LittleEndian, uint32(0))    // colors used
	binary.Write(f, binary.LittleEndian, uint32(0))    // colors important

	// Pixel data — bottom-up BGRA.
	for y := h - 1; y >= 0; y-- {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			f.Write([]byte{byte(b >> 8), byte(g >> 8), byte(r >> 8), byte(a >> 8)})
		}
	}

	// AND mask — all zeros (fully opaque).
	andRow := make([]byte, andMaskRowSize)
	for y := 0; y < h; y++ {
		f.Write(andRow)
	}
}
