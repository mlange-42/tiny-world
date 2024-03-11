package res

import (
	"io/fs"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

const fontFile = "data/fonts/LessRoundBox.ttf"

// Fonts resource for access to UI fonts.
type Fonts struct {
	Default font.Face
	Title   font.Face
}

func NewFonts(fSys fs.FS) Fonts {
	content, err := fs.ReadFile(fSys, fontFile)
	if err != nil {
		log.Fatal("error loading font file: ", err)
	}
	tt, err := opentype.Parse(content)
	if err != nil {
		log.Fatal(err)
	}

	defaultFace, err := makeSize(tt, 22)
	if err != nil {
		log.Fatal(err)
	}
	titleFace, err := makeSize(tt, 32)
	if err != nil {
		log.Fatal(err)
	}

	return Fonts{
		Default: defaultFace,
		Title:   titleFace,
	}
}

func makeSize(tt *sfnt.Font, size int) (font.Face, error) {
	const dpi = 72
	fontFace, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}
	fontFace = text.FaceWithLineHeight(fontFace, float64(size))
	return fontFace, nil
}
