package util

import (
	"encoding/json"
	"image"
	"image/png"
	"io/fs"
	"os"
)

func FromJson(file string, obj any) error {
	ff, err := os.Open(file)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(ff)
	decoder.DisallowUnknownFields()
	return decoder.Decode(obj)
}

func ToJson(file string, obj any) error {
	js, err := json.MarshalIndent(obj, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, js, 0666)
}

func FromJsonFs(f fs.FS, file string, obj any) error {
	ff, err := f.Open(file)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(ff)
	decoder.DisallowUnknownFields()
	return decoder.Decode(obj)
}

func ReadImage(p string) (image.Image, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	baseSprite, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return baseSprite, nil
}

func WriteImage(file string, img image.Image) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}
