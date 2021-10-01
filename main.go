package main

import (
	"fmt"
	"github.com/Konstantsiy/kmeans/binarization"
	"github.com/Konstantsiy/kmeans/blur"
	"github.com/Konstantsiy/kmeans/characteristic"
	"github.com/Konstantsiy/kmeans/cluster"
	"github.com/Konstantsiy/kmeans/util"
	"github.com/anthonynsimon/bild/imgio"
	"math"
	"os"
)

func main() {
	filename, level, k, _ := util.PrepareVars()
	curDir, _ := os.Getwd()
	path := curDir+"/images/"

	img, _ := imgio.Open(path+filename+".jpg")
	imgBin1 := binarization.BinarizeImageWithLevel(img, level)
	util.SavePNG(imgBin1, path, filename, "bin_1")
	
	imgBlur := blur.Gaussian(img, 3.3)
	imgBin2 := binarization.BinarizeImageWithLevel(imgBlur, level)
	util.SavePNG(imgBin2, path, filename, "bin_2")

	bm := binarization.GetBinMap(*imgBin2)
	objects, _ := binarization.FindObjectsRec(bm)
	fmt.Printf("objects count: %d\n", len(objects))

	for i, cors := range objects {
		if len(cors) == 50 {
			delete(objects, i)
		}
	}

	var obj_chars []characteristic.ObjectCharacteristic
	for k, v := range objects {
		s, p, c, e, o, w, h := characteristic.CalcCharacteristics(bm, v)
		if s == 1 || p == 1 {
			continue
		}
		fmt.Println("--------------------")
		fmt.Printf("object id: %d\n", k)
		fmt.Printf("geometric:\tsquare: %d \tperimeter: %d \tcompact: %.2f \telongation: %.2f " +
			"\torientation: %.2f\tmass center: (%.2f, %.2f)\n", s, p, c, e, o, w, h)

		obj_chars = append(obj_chars, characteristic.ObjectCharacteristic{
			ObjectID: k,
			Ch: characteristic.Characteristic{Square: s, Perimeter: p},
		})

		ph := characteristic.CalcPhotometrics(img, v)
		fmt.Printf("photometric:\taverage: (%.2f, %.2f, %.2f)\tgray average: %.2f\tdelta: " +
			"(%.2f, %.2f, %.2f)\tgrey disp: %.2f\n",
			ph.AverageRGB.R, ph.AverageRGB.G, ph.AverageRGB.B, ph.AverageGray,
			ph.DispRGB.R, ph.DispRGB.G, ph.DispRGB.B, ph.DispGray)
		fmt.Println("--------------------")
	}

	dataset := cluster.PrepareDataset(obj_chars)

	clusters := cluster.RunKMeans(dataset, k)

	for i, cl := range clusters { // отображение координат центров кластеров
		fmt.Printf("%d centered at (%.f, %.f)\n", i+1, cl.Center.X, cl.Center.Y)
	}

	objects_colors := make(map[byte]int)
	for _, ob_ch := range obj_chars {
		for cl_i, cl := range clusters {
			for _, point := range cl.Points {
				sq := int(math.Round(point.X))
				per := int(math.Round(point.Y))
				if ob_ch.Ch.Square == sq && ob_ch.Ch.Perimeter == per {
					objects_colors[ob_ch.ObjectID] = cl_i+1
					break
				}
			}
		}
	}

	for obj_k, cors := range objects {
		if color_i, ok := objects_colors[obj_k]; ok {
			for _, c := range cors {
				bm[c.H][c.W] = byte(color_i)
			}
		}
	}

	imgRes := binarization.BinMapToImage(bm, *imgBin2)
	util.SavePNG(imgRes, path, filename, "bin_3")
}