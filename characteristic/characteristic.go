package characteristic

import (
	"image"
	"kmeans-test/binarization"
	"math"
)

type Characteristic struct {
	Square int
	Perimeter int
}

type RGB struct {
	R float64
	G float64
	B float64
}

type Photometrics struct {
	AverageRGB RGB
	AverageGray float64
	DispRGB RGB
	DispGray float64
}

func isBoundary(bm [][]byte, h, w int) bool {
	if h == 0 || h == len(bm)-1 || w == 0 || w == len(bm[0])-1 {
		return true
	}
	return bm[h+1][w] == 0 || bm[h-1][w] == 0 || bm[h][w+1] == 0 || bm[h][w-1] == 0
}

func CalcPerim(bm [][]byte, coordinates binarization.Coordinates) int {
	n := 0
	for _, c := range coordinates {
		if isBoundary(bm, c.H, c.W) {
			n++
		}
	}

	return n
}

func moment(i, j int, wMean, hMean float64, coordinates binarization.Coordinates) float64 {
	var result float64
	for _, c := range coordinates {
		result += math.Pow(float64(c.W)-wMean, float64(i)) * math.Pow(float64(c.H)-hMean, float64(j))
	}
	return result
}

func CalcCharacteristics(bm [][]byte, coordinates binarization.Coordinates) (int, int, float64, float64, float64, float64, float64) {
	square := len(coordinates)
	perimeter := CalcPerim(bm, coordinates)
	compact := math.Pow(float64(perimeter), 2) / float64(square)

	sumW, sumH := 0, 0
	for _, c := range coordinates {
		sumH += c.H
		sumW += c.W
	}

	hMean := float64(sumH) / float64(square)
	wMean := float64(sumW) / float64(square)

	m02 := moment(0, 2, wMean, hMean, coordinates)
	m20 := moment(2, 0, wMean, hMean, coordinates)
	m11 := moment(1, 1, wMean, hMean, coordinates)

	nominator := m20 + m02 + math.Sqrt(math.Pow(m20-m02, 2)+4*math.Pow(m11, 2))
	denominator := m20 + m02 - math.Sqrt(math.Pow(m20-m02, 2)+4*math.Pow(m11, 2))

	elongation := nominator / denominator

	orientation := 0.5 * math.Atan((2 * m11) / (m20 - m02))

	return square, perimeter, compact, elongation, orientation, wMean, hMean
}

func GetMinMaxFloat64(arr []float64) (float64, float64) {
	min, max := arr[0], arr[0]
	for _, v := range arr {
		if v < min {
			min = v
		} else if v > max {
			max = v
		}
	}
	return min, max
}

func GetMinMaxUint32(arr []uint32) (uint32, uint32) {
	min, max := arr[0], arr[0]
	for _, v := range arr {
		if v < min {
			min = v
		} else if v > max {
			max = v
		}
	}
	return min, max
}

func CalcPhotometrics(img image.Image, coordinates binarization.Coordinates) Photometrics {
	var rList, gList, bList []uint32
	var grayList []float64
	var rSum, gSum, bSum uint32
	var graySum float64

	for _, c := range coordinates {
		r, g, b, _ := img.At(c.H, c.W).RGBA()
		rSum += r
		gSum += g
		bSum += b

		gray := 0.3 * float64(r) + 0.59 * float64(g) + 0.11 * float64(b)
		graySum += gray

		rList = append(rList, r)
		gList = append(gList, g)
		bList = append(bList, b)
		grayList = append(grayList, gray)
	}

	rAve, gAve, bAve, grayAve := float64(rSum) / float64(len(rList)), float64(gSum) / float64(len(gList)),
		float64(bSum) / float64(len(bList)), graySum / float64(len(grayList))

	rMin, rMax := GetMinMaxUint32(rList)
	gMin, gMax := GetMinMaxUint32(gList)
	bMin, bMax := GetMinMaxUint32(bList)
	grayMin, grayMax := GetMinMaxFloat64(grayList)
	rDisp, gDisp, bDisp, grayDisp := rMax - rMin, gMax - gMin, bMax - bMin, grayMax - grayMin

	return Photometrics{
		AverageRGB: RGB{R: rAve, G: gAve, B: bAve}, AverageGray: grayAve,
		DispRGB: RGB{R: float64(rDisp), G: float64(gDisp), B: float64(bDisp)}, DispGray: grayDisp,
	}
}