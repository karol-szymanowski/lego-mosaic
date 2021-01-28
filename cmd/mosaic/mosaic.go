package mosaic

import (
	"image"
	"image/color"
)

type mosaic [][]PaletteItem
type LegoParts map[string] int

type LegoMosaic struct {
	dx           int
	dy           int
	mosaic       mosaic
	LegoParts    LegoParts
	colorPalette Palette
}

func (l LegoMosaic) set(x int, y int, c PaletteItem) {
	l.LegoParts[c.ColorName]++
	l.mosaic[y][x] = c
}

func (l LegoMosaic) setRgba(x int, y int, c color.RGBA) {
	pItem := PaletteItem{
		Id:        "custom rgb",
		Color:     paletteRgb{
			R: c.R,
			G: c.G,
			B: c.B,
			A: c.A,
		},
		ColorName: "custom rgb",
	}
	l.mosaic[y][x] = pItem
}

func (l LegoMosaic) at(x int, y int) PaletteItem {
	return l.mosaic[y][x]
}


func (l LegoMosaic) LegoPartsCount() int {
	return l.dx * l.dy
}

func (l LegoMosaic) ToImage() image.Image {
	imgRect := image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: l.dx, Y: l.dy},
	}
	img := image.NewRGBA(imgRect)

	for y, items := range l.mosaic {
		for x, c := range items {
			img.Set(x, y, c.toRgb())
		}
	}
	return img
}

func NewLegoMosaic(img image.Image, palette Palette) LegoMosaic  {
	rect := img.Bounds()
	dx, dy := rect.Dx(), rect.Dy()

	mosaic := make([][]PaletteItem, dy)
	for y := range mosaic {
		mosaic[y] = make([]PaletteItem, dx)
		for x := range mosaic[y] {
			R, G, B, A := img.At(x, y).RGBA()
			mosaic[y][x] = PaletteItem{
				Id:        "unknown",
				Color:     paletteRgb{
					R: uint8(R),
					G: uint8(G),
					B: uint8(B),
					A: uint8(A),
				},
				ColorName: "unknown",
			}
		}
	}

	legoParts := make(map[string]int)
	for _, c := range palette {
		legoParts[c.ColorName] = 0
	}

	legoMosaic := LegoMosaic{
		dx:           dx,
		dy:           dy,
		mosaic:       mosaic,
		LegoParts:    legoParts,
		colorPalette: palette,
	}

	return legoMosaic
}



