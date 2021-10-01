package main

import (
	"kmeans-test/binarization"
	"kmeans-test/util"
	"os"
)

func main() {
	filename, level, _, _ := util.PrepareVars()
	curDir, _ := os.Getwd()
	path := curDir+"/images/"

	img, _ := imgio.Open(path+filename+".jpg")
	imgBin1 := binarization.BinarizeImageWithLevel(img, level)
	util.SavePNG(imgBin1, path, filename, "bin_1")
}