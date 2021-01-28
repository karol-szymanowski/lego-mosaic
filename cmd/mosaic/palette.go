package mosaic

import (
	"encoding/json"
	"image/color"
	"io/ioutil"
	"os"
)

type paletteRgb struct {
	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
	A uint8 `json:"a"`
}

type PaletteItem struct {
	Id        string     `json:"id"`
	Color     paletteRgb `json:"color"`
	ColorName string     `json:"name"`
}

type Palette []PaletteItem

type paletteItem interface {
	toRgb() color.RGBA
}

func (p PaletteItem) toRgb() color.RGBA {
	return color.RGBA{
		R: p.Color.R,
		G: p.Color.G,
		B: p.Color.B,
		A: p.Color.A,
	}
}

func LoadColorPalette(palettePath string) (Palette, error) {
	jsonPalette, err := os.Open(palettePath)
	if err != nil {
		return nil, err
	}
	defer jsonPalette.Close()

	byteValue, _ := ioutil.ReadAll(jsonPalette)

	var palette []PaletteItem
	err = json.Unmarshal(byteValue, &palette)
	if err != nil {
		return nil, err
	}

	return palette, nil
}
