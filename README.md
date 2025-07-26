# ASCII Image & GIF Converter

This Go package converts images and GIFs into colored ASCII art that can be displayed directly in the terminal. It supports both single images and animated GIFs, resizing them while preserving aspect ratio and applying ANSI 24-bit  color codes (TrueColor) for a vivid result.

## Features

- **Image to ASCII Conversion**: Converts any image to ASCII art with color.
- **GIF to ASCII Conversion**: Supports animated GIFs, generating ASCII frames with correct delays.
- **Automatic Resizing**:  
  - Set **both** `width` and `height` explicitly, **or**  
  - Provide **only width** or **only height** (set the other to `0`) and the aspect ratio is automatically preserved.
- **Terminal Animation**: Plays ASCII GIFs directly in the terminal with proper frame timing.
- **Custom Character Set**: Uses a defined character set for luminance-based ASCII mapping (Rec. 601).

---

## Installation

```bash
go get github.com/mattparadis/asciiConverter
```

---

## Usage

### Convert a Single Image

```go
package main

import (
    "fmt"
    "log"
    "asciiconverter"
)

func main() {
    // Example: specify only width (height is auto-calculated)
    img, err := asciiconverter.GetAsciiImage("path/to/image.jpg", 100, 0)
    if err != nil {
        log.Fatal(err)
    }
    asciiconverter.PrintImg(img)
}
```

### Convert and Play a GIF

```go
package main

import (
    "log"
    "asciiconverter"
)

func main() {
    // Example: specify only height (width is auto-calculated)
    frames, err := asciiconverter.GetAsciiGif("path/to/animation.gif", 0, 40)
    if err != nil {
        log.Fatal(err)
    }
    asciiconverter.PrintGif(frames, 3) // Play 3 loops
}
```

---

## API Overview

### **GetAsciiImage(path string, width, height int) ([]string, error)**
- Converts a single image to ASCII lines.
- If `width` or `height` is `0`, aspect ratio is preserved automatically.

### **GetAsciiGif(path string, width, height int) ([]\*AsciiGif, error)**
- Converts a GIF into a slice of ASCII frames, each with delay information.

### **PrintImg(img []string)**
- Prints a single ASCII image.

### **PrintGif(frames []\*AsciiGif, loop int)**
- Plays an ASCII GIF animation in the terminal, looping `loop` times.

---

## Requirements

- **Go 1.18+**
- [disintegration/imaging](https://github.com/disintegration/imaging) for image processing.

## License

This project is licensed under the [MIT License](LICENSE).
