package mosaic

import (
	"golang.org/x/image/draw"
	"image"
	"image/png"
	"image/jpeg"
	"os"
)

func LoadImage(filePath string) (image.Image, error)  {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	return img, err
}

func SaveImage(filePath string, image image.Image) error {
	dstFile, err := os.Create(filePath)
	if err != nil {
		return err
	}

	err = png.Encode(dstFile, image)
	if err != nil {
		return err
	}

	defer dstFile.Close()
	return nil
}

func scaleImageTo(src image.Image, rect image.Rectangle, scale draw.Scaler) image.Image {
	dst := image.NewRGBA(rect)
	scale.Scale(dst, rect, src, src.Bounds(), draw.Over, nil)
	return dst
}

func ResizeImage(originalImg image.Image, size image.Point) (error, image.Image) {
	scaledImgRect := image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: size,
	}
	scaledImg := scaleImageTo(originalImg, scaledImgRect, draw.BiLinear)

	return nil, scaledImg
}
