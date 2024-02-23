package res

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const fontFile = "./assets/fonts/JupiteroidRegular.ttf"
const fontSize = 24

type Fonts struct {
	Default font.Face
}

func NewFonts() Fonts {
	content, err := os.ReadFile(fontFile)
	if err != nil {
		log.Fatal("error loading font file: ", err)
	}

	tt, err := opentype.Parse(content)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	fontFace, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(fontSize),
		DPI:     dpi,
		Hinting: font.HintingFull, // Use quantization to save glyph cache images.
	})
	if err != nil {
		log.Fatal(err)
	}
	// Adjust the line height.
	fontFace = text.FaceWithLineHeight(fontFace, 30)

	return Fonts{
		Default: fontFace,
	}
}
