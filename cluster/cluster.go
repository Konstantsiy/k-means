package cluster

import (
	"github.com/Konstantsiy/kmeans/characteristic"
	"github.com/sirupsen/logrus"
	"math"
	"math/rand"
	"time"
)

type Point struct {
	Square     float64
	Perimeter  float64
	Compact    float64
	Elongation float64
}

func (p Point) distance(p2 Point) float64 {
	return math.Sqrt(
		math.Pow(p.Square-p2.Square, 2) +
			math.Pow(p.Perimeter-p2.Perimeter, 2) +
			math.Pow(p.Compact-p2.Compact, 2) +
			math.Pow(p.Elongation-p2.Elongation, 2))
}

type Cluster struct {
	Center Point
	Points []Point
}

func (cluster *Cluster) repositionCenter() {
	var squaresSum, perimetersSum, compactsSum, elongationsSum float64
	var count = len(cluster.Points)

	for i := 0; i < count; i++ {
		squaresSum += cluster.Points[i].Square
		perimetersSum += cluster.Points[i].Perimeter
		compactsSum += cluster.Points[i].Compact
		elongationsSum += cluster.Points[i].Elongation
	}

	cluster.Points = []Point{}
	cluster.Center = Point{
		squaresSum / float64(count),
		perimetersSum / float64(count),
		compactsSum / float64(count),
		elongationsSum / float64(count),
	}
}

func initClusters(dataset []Point, k int) []Cluster {
	rand.Seed(time.Now().UnixNano())
	var clusters []Cluster

	for i := 0; i < k; i++ {
		center := dataset[rand.Intn(len(dataset))]
		clusters = append(clusters, Cluster{Center: center, Points: []Point{}})
	}

	return clusters
}

func RunKMeans(dataset []Point, k int) []Cluster {
	pointsClusterIndex := make([]int, len(dataset))
	clusters := initClusters(dataset, k) // инициализация рандомных кластеров из датасета
	maxIter, counter := 30, 0 // максимум интраций, счетчик инераций

	for hasChanged := true; hasChanged; { // ожидаем, пока изменения не прекратятся
		hasChanged = false
		for i := 0; i < len(dataset); i++ {
			var minDist float64 // считаем расстояния от точки до ближайшего кластера
			var updatedClusterIndex int // индекс изменившегося кластера
			for j := 0; j < len(clusters); j++ {
				tmpDist := dataset[i].distance(clusters[j].Center)
				if minDist == 0 || tmpDist < minDist {
					minDist = tmpDist
					updatedClusterIndex = j
				}
			}
			// добавляем точку из датасета к ближайшему кластеру
			clusters[updatedClusterIndex].Points = append(clusters[updatedClusterIndex].Points, dataset[i])
			// помечаем изменившийся кластер
			if pointsClusterIndex[i] != updatedClusterIndex {
				pointsClusterIndex[i] = updatedClusterIndex
				hasChanged = true
			}
		}
		// если были изменения, то пересчитываем центры масс кластеров
		if hasChanged {
			for i := 0; i < len(clusters); i++ {
				clusters[i].repositionCenter()
			}
		}
		counter++
		// проверка на превышение числа итераций
		if counter >= maxIter {
			logrus.Info("exceeded the maximum number of iterations")
			return clusters
		}
	}
	return clusters
}

func PrepareDataset(objectsChars []characteristic.ObjectCharacteristic) []Point {
	var dataset []Point
	for _, o_ch := range objectsChars {
		dataset = append(dataset, Point{
			float64(o_ch.Ch.Square),
			float64(o_ch.Ch.Perimeter),
			o_ch.Ch.Compact,
			o_ch.Ch.Elongation,
		})
	}
	return dataset
}