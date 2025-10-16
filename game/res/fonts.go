package res

import (
	"bytes"
	"io/fs"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const fontFile = "data/fonts/LessRoundBox.ttf"

// Fonts resource for access to UI fonts.
type Fonts struct {
	Default text.Face
	Title   text.Face
}

func NewFonts(fSys fs.FS) Fonts {
	data, err := fs.ReadFile(fSys, fontFile)
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(data)
	source, err := text.NewGoTextFaceSource(reader)
	if err != nil {
		log.Fatal("error loading font file: ", err)
	}
	defaultFace := text.GoTextFace{
		Source: source,
		Size:   22,
	}
	titleFace := text.GoTextFace{
		Source: source,
		Size:   48,
	}

	return Fonts{
		Default: &defaultFace,
		Title:   &titleFace,
	}
}

/*
func makeSize(tt *sfnt.Font, size int) (text.Face, error) {
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
*/
