package asciiconverter

import (
	"bufio"
	"fmt"
	"image"
	"image/gif"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
)

const charset = " .'`^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"

var runes = []rune(charset)

type AsciiGif struct {
	Lines []string
	Delay time.Duration
}

// openImg loads an image from the specified path and returns it as an image.Image.
func openImg(path string) (image.Image, error) {
	path = os.ExpandEnv(path)
	img, err := imaging.Open(path)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// openGif loads a GIF file from the specified path and returns all frames and their delays.
func openGif(path string) ([]*image.Paletted, []int, error) {
	path = os.ExpandEnv(path)
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	g, err := gif.DecodeAll(f)
	if err != nil {
		return nil, nil, err
	}
	return g.Image, g.Delay, nil
}

// resizeImg resizes the image to the specified width and height while maintaining aspect ratio if one dimension is zero.
func resizeImg(img image.Image, width, height int) (image.Image, error) {
	if width <= 0 && height <= 0 {
		return img, nil
	} else if width <= 0 {
		width = int(float64(height) * float64(img.Bounds().Dx()) / float64(img.Bounds().Dy()))
	} else if height <= 0 {
		height = int(float64(width) * float64(img.Bounds().Dy()) / float64(img.Bounds().Dx()))
	}
	resizedImg := imaging.Resize(img, width, height, imaging.Lanczos)
	if resizedImg == nil {
		return nil, os.ErrInvalid
	}
	return resizedImg, nil
}

// luminanceFromRGBA calculates the luminance of a pixel based on ITU-R BT.601 standard.
func luminanceFromRGBA(r, g, b int) float64 {
	// Standard ITU-R BT.601 (Rec. 601) formula
	return 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
}

// pickChar selects a character from the runes slice based on the grayscale value.
func pickChar(gray uint8, runes []rune) rune {
	idx := int(gray) * (len(runes) - 1) / 255
	return runes[idx]
}

// writeColoredRune writes a colored character (ANSI 24-bit color) to a string builder.
func writeColoredRune(sb *strings.Builder, r, g, b uint8, ch rune) {
	sb.WriteString("\x1b[48;2;0;0;0m\x1b[38;2;")
	sb.WriteString(strconv.Itoa(int(r)))
	sb.WriteByte(';')
	sb.WriteString(strconv.Itoa(int(g)))
	sb.WriteByte(';')
	sb.WriteString(strconv.Itoa(int(b)))
	sb.WriteByte('m')
	sb.WriteRune(ch)
	sb.WriteRune(ch)
}

// convertImageToAscii converts an image to a slice of strings, where each string is a colored ASCII line.
func convertImageToAscii(img image.Image) ([]string, error) {
	if img == nil {
		return nil, os.ErrInvalid
	}
	bounds := img.Bounds()
	asciiImg := []string{}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		strBuild := strings.Builder{}

		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			p := img.At(x, y)
			r, g, b, _ := p.RGBA()
			gray := uint8(luminanceFromRGBA(int(r>>8), int(g>>8), int(b>>8)))
			char := pickChar(gray, runes)
			writeColoredRune(&strBuild, uint8(r>>8), uint8(g>>8), uint8(b>>8), char)

		}
		strBuild.WriteString("\x1b[0m\n")
		asciiImg = append(asciiImg, strBuild.String())
	}
	return asciiImg, nil
}

// GetAsciiImage converts a single image from the specified path into ASCII art.
// The output is a slice of strings representing the ASCII image.
// width or height can be set to 0 to automatically preserve aspect ratio.
func GetAsciiImage(path string, width, height int) ([]string, error) {
	img, err := openImg(path)
	if err != nil {
		return nil, err
	}
	resizedImg, err := resizeImg(img, width, height)
	if err != nil {
		return nil, err
	}
	asciiImg, err := convertImageToAscii(resizedImg)
	if err != nil {
		return nil, err
	}
	return asciiImg, nil
}

// GetAsciiGif converts a GIF from the specified path into a slice of AsciiGif frames.
// Each frame contains ASCII art lines and its display delay.
// width or height can be set to 0 to automatically preserve aspect ratio.
func GetAsciiGif(path string, width, height int) ([]*AsciiGif, error) {
	imgs, delays, err := openGif(path)
	if err != nil {
		return nil, err
	}
	var asciiGifs []*AsciiGif
	for i, img := range imgs {
		resizedImg, err := resizeImg(img, width, height)
		if err != nil {
			return nil, err
		}
		asciiImg, err := convertImageToAscii(resizedImg)
		if err != nil {
			return nil, err
		}
		asciiGifs = append(asciiGifs, &AsciiGif{
			Lines: asciiImg,
			Delay: time.Duration(delays[i]) * time.Millisecond,
		})
	}
	return asciiGifs, nil
}

// PrintImg prints a single ASCII image to the terminal.
func PrintImg(img []string) {
	for _, line := range img {
		fmt.Print(line)
	}
}

// PrintGif plays an ASCII GIF animation in the terminal.
// The parameter 'loop' specifies the number of times the animation should repeat.
func PrintGif(frames []*AsciiGif, loop int) {
	const (
		clearScreen = "\033[2J"
		cursorHome  = "\033[H"
		hideCursor  = "\033[?25l"
		showCursor  = "\033[?25h"
		resetColor  = "\033[0m"
	)

	fmt.Print(hideCursor, clearScreen, cursorHome)
	defer fmt.Print(showCursor, resetColor)

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	minDelay := 33 * time.Millisecond // minimum frame delay to avoid flickering

	for {
		for _, f := range frames {
			fmt.Fprint(w, cursorHome)
			for _, line := range f.Lines {
				w.WriteString(line)
			}
			w.Flush()

			d := 10 * f.Delay
			if d <= 0 {
				d = minDelay
			}
			time.Sleep(d)
		}
		loop--
		if loop <= 0 {
			break
		}
	}
}
