package blur

import (
	"github.com/Konstantsiy/kmeans/util"
	"image"
	"image/color"
	"math"
)

func Gaussian(img image.Image, sigma float64) *image.RGBA {
	cp := util.AsRGBA(img)

	rscl := make([]uint8, cp.Rect.Dx()*cp.Rect.Dy())
	gscl := make([]uint8, cp.Rect.Dx()*cp.Rect.Dy())
	bscl := make([]uint8, cp.Rect.Dx()*cp.Rect.Dy())

	i := 0
	for y := 0; y < cp.Rect.Dy(); y++ {
		for x := 0; x < cp.Rect.Dx(); x++ {
			r, g, b, _ := cp.At(x, y).RGBA()
			rscl[i] = uint8(r)
			gscl[i] = uint8(g)
			bscl[i] = uint8(b)
			i++
		}
	}

	rtcl := make([]uint8, len(rscl))
	gtcl := make([]uint8, len(rscl))
	btcl := make([]uint8, len(rscl))
	gaussBlur_4(rscl, rtcl, cp.Rect.Dx(), cp.Rect.Dy(), int(sigma))
	gaussBlur_4(gscl, gtcl, cp.Rect.Dx(), cp.Rect.Dy(), int(sigma))
	gaussBlur_4(bscl, btcl, cp.Rect.Dx(), cp.Rect.Dy(), int(sigma))

	i = 0
	for y := 0; y < cp.Rect.Dy(); y++ {
		for x := 0; x < cp.Rect.Dx(); x++ {
			cp.Set(x, y, color.RGBA{rtcl[i], gtcl[i], btcl[i], 255})
			i++
		}
	}

	return cp
}

func boxesForGauss(sigma float64, n int) []int {
	wIdeal := math.Sqrt((12 * sigma * sigma / float64(n)) + 1)
	wl := int(math.Floor(wIdeal))
	if wl%2 == 0 {
		wl--
	}

	mIdeal := (12*sigma*sigma - float64(n*wl*wl+4*n*wl+3*n)) / float64(-4*wl-4)
	// Round to the nearest number
	m := int(mIdeal + 0.5)

	sizes := make([]int, n)
	for i := 0; i < n; i++ {
		if i < int(m) {
			sizes[i] = wl
		} else {
			sizes[i] = wl + 2
		}
	}

	return sizes
}

func gaussBlur_4(scl, tcl []uint8, w, h, r int) {
	bxs := boxesForGauss(float64(r), 3)
	boxBlur_4(scl, tcl, w, h, (bxs[0]-1)/2)
	boxBlur_4(tcl, scl, w, h, (bxs[1]-1)/2)
	boxBlur_4(scl, tcl, w, h, (bxs[2]-1)/2)
}

func boxBlur_4(scl, tcl []uint8, w, h, r int) {
	for i := 0; i < len(scl); i++ {
		tcl[i] = scl[i]
	}
	boxBlurH_4(tcl, scl, w, h, r) // Horizontal blur
	boxBlurT_4(scl, tcl, w, h, r) // Total blur
}

func boxBlurH_4(scl, tcl []uint8, w, h, r int) {
	var iarr float64 = 1 / float64(r+r+1)
	for i := 0; i < h; i++ {
		ti := i * w
		li := ti
		ri := ti + r

		fv := int(scl[ti])
		lv := int(scl[ti+w-1])

		val := (r + 1) * fv

		for j := 0; j < r; j++ {
			val += int(scl[ti+j])
		}

		for j := 0; j <= r; j++ {
			ri++
			val += int(scl[ri-1]) - fv
			ti++
			tcl[ti-1] = uint8(float64(val)*iarr + 0.5)
		}

		for j := r + 1; j < w-r; j++ {
			ri++
			li++
			val += int(scl[ri-1]) - int(scl[li-1])
			ti++
			tcl[ti-1] = uint8(float64(val)*iarr + 0.5)
		}

		for j := w - r; j < w; j++ {
			li++
			val += lv - int(scl[li-1])
			ti++
			tcl[ti-1] = uint8(float64(val)*iarr + 0.5)
		}
	}
}

func boxBlurT_4(scl, tcl []uint8, w, h, r int) {
	var iarr float64 = 1 / float64(r+r+1)
	for i := 0; i < w; i++ {
		ti := i
		li := ti
		ri := ti + r*w

		fv := int(scl[ti])
		lv := int(scl[ti+w*(h-1)])
		val := (r + 1) * fv

		for j := 0; j < r; j++ {
			val += int(scl[ti+j*w])
		}

		for j := 0; j <= r; j++ {
			val += int(scl[ri]) - fv
			tcl[ti] = uint8(float64(val)*iarr + 0.5)
			ri += w
			ti += w
		}

		for j := r + 1; j < h-r; j++ {
			val += int(scl[ri]) - int(scl[li])
			tcl[ti] = uint8(float64(val)*iarr + 0.5)
			li += w
			ri += w
			ti += w
		}

		for j := h - r; j < h; j++ {
			val += lv - int(scl[li])
			tcl[ti] = uint8(float64(val)*iarr + 0.5)
			li += w
			ti += w
		}
	}
}



