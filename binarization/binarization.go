package binarization

import (
	"github.com/Konstantsiy/kmeans/util"
	"image"
	"image/color"
)

type Coordinate struct {
	H int
	W int
}

type Coordinates []Coordinate

var colors = [][3]byte{
	{255, 0, 0},
	{0, 255, 0},
	{0, 0, 255},
	{255, 255, 0},
	{0, 255, 255},
	{255, 0, 255},
	{100, 0, 0},
	{0, 100, 0},
	{0, 0, 100},
	{100, 100, 0},
	{0, 100, 100},
	{100, 0, 100},
	{175, 0, 0},
	{0, 175, 0},
	{0, 0, 175},
	{175, 175, 0},
	{175, 0, 175},
	{0, 175, 175},
}

func BinMapToImage(bm [][]byte, img image.Gray) image.Image {
	src := util.AsRGBA(&img)
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			if bm[y][x] != 0 {
				colorID := bm[y][x]
				if colorID >= 18 {
					continue
				}
				c := colors[colorID]
				pos := y * src.Stride + x * 4

				src.Pix[pos+0] = c[0]
				src.Pix[pos+1] = c[1]
				src.Pix[pos+2] = c[2]
			}
		}
	}
	return src
}

func GetBinMap(img image.Gray) [][]byte {
	bounds := img.Bounds()
	binMap := make([][]byte, bounds.Dy())

	for y := 0; y < bounds.Dy(); y++ {
		binMap[y] = make([]byte, bounds.Dx())
		for x := 0; x < bounds.Dx(); x++ {
			pos := y * img.Stride + x

			if img.Pix[pos] == 0xFF {
				binMap[y][x] = 1
			} else if img.Pix[pos] == 0x00 {
				binMap[y][x] = 0
			}
		}
	}

	return binMap
}

func fill(bm [][]byte, x, y int, c byte, objects map[byte]Coordinates) {
	if bm[x][y] == 1 {
		bm[x][y] = c
		if _, ok := objects[c]; !ok {
			objects[c] = Coordinates{}
		}
		objects[c] = append(objects[c], Coordinate{H: x, W: y})
		if x > 0 {
			fill(bm, x - 1, y, c, objects)
		}
		if x < len(bm) - 1 {
			fill(bm, x + 1, y, c, objects)
		}
		if y > 0 {
			fill(bm, x, y - 1, c, objects)
		}
		if y < len(bm[0]) - 1 {
			fill(bm, x, y + 1, c, objects)
		}
	}
}

func FindObjectsRec(bm [][]byte) (map[byte]Coordinates, [][]byte) {
	objects := make(map[byte]Coordinates)
	var c byte = 1
	for i := 0; i < len(bm); i++ {
		for j := 0; j < len(bm[0]); j++ {
			c++
			fill(bm, i, j, c, objects)
		}
	}
	return objects, bm
}

func BinarizeImageWithLevel(img image.Image, level uint8) *image.Gray {
	src := util.AsRGBA(img)
	bounds := src.Bounds()

	dst := image.NewGray(bounds)

	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			srcPos := y*src.Stride + x*4
			dstPos := y*dst.Stride + x

			c := src.Pix[srcPos : srcPos+4]
			r := util.Rank(color.RGBA{c[0], c[1], c[2], c[3]})

			// transparent pixel is always white
			if c[0] == 0 && c[1] == 0 && c[2] == 0 && c[3] == 0 {
				dst.Pix[dstPos] = 0xFF
				continue
			}

			if uint8(r) >= level {
				dst.Pix[dstPos] = 0xFF
			} else {
				dst.Pix[dstPos] = 0x00
			}
		}
	}

	return dst
}
