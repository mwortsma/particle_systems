package matutil

import "fmt"

type Mat [][]int

type Vec []int

func Create(rows, cols int) Mat {
	m := make(Mat, rows)
	for i := 0; i < rows; i++ {
		m[i] = make(Vec, cols)
	}
	return m
}

func (mat Mat) Dims() (int, int) {
	return len(mat), len(mat[0])
}

func PrintMat(mat Mat) {
	r, _ := mat.Dims()
	for i := 0; i < r; i++ {
		fmt.Println(mat[i])
	}
}
