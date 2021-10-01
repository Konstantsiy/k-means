package cluster

import (
	"github.com/sirupsen/logrus"
	"math"
	"math/rand"
	"time"
)

type Point struct {
	X float64
	Y float64
}

func (p Point) distance(p2 Point) float64 {
	return math.Sqrt(math.Pow(p.X-p2.X, 2) + math.Pow(p.Y-p2.Y, 2))
}

type Cluster struct {
	Center Point
	Points []Point
}

func (cluster *Cluster) repositionCenter() {
	var x, y float64
	var clusterCount = len(cluster.Points)

	for i := 0; i < clusterCount; i++ {
		x += cluster.Points[i].X
		y += cluster.Points[i].Y
	}
	cluster.Points = []Point{}
	cluster.Center = Point{x / float64(clusterCount), y / float64(clusterCount)}
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