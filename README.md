# k-means

The k-means algorithm is a data clustering algorithm. The purpose of the clustering task is to divide a set of objects into clusters (classes) based on geometric and photometric features: 

- area 
- perimeter
- compactness
- elongation
- average brightness
- average color 
- dispersions

# How it works

## 1. The original image

![original](https://github.com/Konstantsiy/k-means/blob/master/images/5.jpg)

## 2. Binarization

The binarization of the source image works by default
at level 200, as this is the most optimal option selected
during testing:

![binarization](https://github.com/Konstantsiy/k-means/blob/master/images/5_bin_1.png)

## 3. Gaussian Blur

The value of each pixel of the image is "averaged” with neighboring
pixels in accordance with the weight coefficients of the "floating
window", while on the first pass averaging occurs by
the "horizontal“ neighbors, and on the second — on the ”vertical":

![blur](https://github.com/Konstantsiy/k-means/blob/master/images/5_bin_2.png)

## 4. Object search (recursion) on the binary matrix

The search result is an associative array, where the key is the object id, and the value is an array of coordinates of points of this object on a binary matrix.

## 5. Calculation and of geometric and photometric features

Then, based on them, the dataset is filled in and
clustering is performed with a given number of clusters `k`.

## 6. Clustering (k-means) and applying the resulting binary matrix to the image

Cluster analysis is carried out on the basis of 4 geometric features of each object: area, volume, compactness and elongation.

In order to remove minor objects (single
bright pixels or some irregularities)
objects with a perimeter less than 50 do not participate in clustering:

![clustering](https://github.com/Konstantsiy/k-means/blob/master/images/5_bin_3.png)

# How to run
```bigquery
go run main.go <image_name.jpg> <binarization_level> <clusters_count>
```
