package mosaic

import (
	"image/color"
	"math"
)

func QuantizeMosaic(mosaic LegoMosaic, dithering bool) (LegoMosaic, error) {
	for y := 0; y < mosaic.dy; y++ {
		for x := 0; x < mosaic.dx; x++ {
			oldColor := mosaic.at(x, y)
			newColor := findClosestColor(oldColor.toRgb(), mosaic.colorPalette)
			mosaic.set(x, y, newColor)

			// https://en.wikipedia.org/wiki/Floyd%E2%80%93Steinberg_dithering
			if dithering {
				quantErr := calcErr(newColor.toRgb(), oldColor.toRgb())

				if x < mosaic.dx - 1 {
					mosaic.setRgba(x + 1, y, calcQuantErrColor(mosaic.at(x + 1, y).toRgb(), quantErr, 7.0 / 16.0))
					if y < mosaic.dy - 1 {
						mosaic.setRgba(x + 1, y + 1, calcQuantErrColor(mosaic.at(x + 1, y + 1).toRgb(), quantErr, 1.0 / 16.0))
						if x > 0 {
							mosaic.setRgba(x - 1, y + 1, calcQuantErrColor(mosaic.at(x - 1, y + 1).toRgb(), quantErr, 3.0 / 16.0))
						}
					}
				}
				if y < mosaic.dy - 1 {
					mosaic.setRgba(x, y + 1, calcQuantErrColor(mosaic.at(x, y + 1).toRgb(), quantErr, 5.0 / 16.0))
				}
			}
		}
	}

	return mosaic, nil
}

type quantError struct {
	r float64
	g float64
	b float64
	a float64
}

func findClosestColor(c color.RGBA, palette Palette) PaletteItem {
	result, bestSum := 0, float64(1<<64-1)
	for index, p := range palette {
		paletteColor := p.toRgb()
		sum := compareColorsRedmean(c, paletteColor)

		if sum < bestSum {
			result, bestSum = index, sum
			if sum == 0 {
				break
			}
		}
	}

	return palette[result]
}

// https://en.wikipedia.org/wiki/Color_difference
func compareColorsRedmean(c1 color.RGBA, c2 color.RGBA) float64 {
	R1, G1, B1 := float64(c1.R), float64(c1.G), float64(c1.B)
	R2, G2, B2 := float64(c2.R), float64(c2.G), float64(c2.B)

	r := (R1 + R2) / 2

	diff := (2.0 + r / 256.0) * sqDiffFloat(R2, R1) + 4.0 * sqDiffFloat(G2, G1) + (2.0 + (255.0 - r) / 256.0) * sqDiffFloat(B2, B1)

	return math.Sqrt(diff)
}

func compareColorsSimple(c1 color.RGBA, c2 color.RGBA) float64 {
	R1, G1, B1 := float64(c1.R), float64(c1.G), float64(c1.B)
	R2, G2, B2 := float64(c2.R), float64(c2.G), float64(c2.B)
	sqDiff := sqDiffFloat(R2, R1) +
		sqDiffFloat(G2, G1) +
		sqDiffFloat(B2, B1)
	return math.Sqrt(sqDiff)
}

func calcErr(newColor color.RGBA, oldColor color.RGBA) quantError  {
	r := float64(oldColor.R) - float64(newColor.R)
	g := float64(oldColor.G) - float64(newColor.G)
	b := float64(oldColor.B) - float64(newColor.B)
	a := 0.0

	return quantError{
		r: r,
		g: g,
		b: b,
		a: a,
	}
}

func calcQuantErrColor(oldColor color.Color, err quantError, multiplier float64) color.RGBA  {
	oR, oG, oB, _ := oldColor.RGBA()
	return color.RGBA{
		R: uint8(float64(oR) + math.Round(err.r * multiplier)),
		G: uint8(float64(oG) + math.Round(err.g * multiplier)),
		B: uint8(float64(oB) + math.Round(err.b * multiplier)),
		A: 255,
	}
}

func sqDiffFloat(x, y float64) float64 {
	d := x - y
	return math.Pow(d, 2)
}
