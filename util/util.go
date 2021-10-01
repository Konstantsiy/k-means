package util

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"strconv"
)

func PrepareVars() (string, uint8, int, error) {
	args := os.Args

	filename := args[1]
	level, _ := strconv.ParseUint(args[2], 10, 8)
	k, _ := strconv.Atoi(args[3])

	return filename, uint8(level), k, nil
}

func AsRGBA(src image.Image) *image.RGBA {
	bounds := src.Bounds()
	res := image.NewRGBA(bounds)
	draw.Draw(res, bounds, src, bounds.Min, draw.Src)
	return res
}

func Rank(c color.RGBA) float64 {
	return float64(c.R)*0.3 + float64(c.G)*0.6 + float64(c.B)*0.1
}

func SaveJPG(img image.Image, path, filename, postfix string) error {
	filename += ".jpg"
	f, err := os.Create(path+ filename[0:len(filename)-4]+"_"+postfix+filename[len(filename)-4:])
	if err != nil {
		return err
	}

	err = jpeg.Encode(f, img, &jpeg.Options{Quality: 99})
	if err != nil {
		return err
	}
	return nil
}

func SavePNG(img image.Image, path, filename, postfix string) error {
	filename += ".png"
	f, err := os.Create(path+ filename[0:len(filename)-4]+"_"+postfix+filename[len(filename)-4:])
	if err != nil {
		return err
	}

	var enc png.Encoder
	enc.CompressionLevel = 90
	err = enc.Encode(f, img)
	if err != nil {
		return err
	}
	return nil
}
