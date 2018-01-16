package matutil

import (
	"fmt"
)

type Mat [][]int

type Vec []int

func (m Mat) String() string {
	// TODO
	return ""
}

func (v Vec) String() string {
	// TODO
	s := fmt.Sprintf("%v", []int(v))
	//fmt.Println(s)
	return s
}

func Create(rows, cols int) Mat {
	m := make(Mat, rows)
	for i := 0; i < rows; i++ {
		m[i] = make(Vec, cols)
	}
	return m
}

func (mat Mat) Col(j int) Vec {
	r := len(mat)
	v := make(Vec, r)
	for i := 0; i < r; i++ {
		v[i] = mat[i][j]
	}
	return v
}

func (mat Mat) Dims() (int, int) {
	return len(mat), len(mat[0])
}

func (mat Mat) Print() {
	r, _ := mat.Dims()
	for i := 0; i < r; i++ {
		fmt.Println(mat[i])
	}
}
