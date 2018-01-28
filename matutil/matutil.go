package matutil

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Mat [][]int

type Vec []int

func (m Mat) String() string {
	s := strings.Replace(fmt.Sprintf("%v", [][]int(m)), " ", ",", -1)
	return s
}

func (v Vec) String() string {
	s := strings.Replace(fmt.Sprintf("%v", []int(v)), " ", ",", -1)
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

func (mat Mat) Cols(cols []int) Mat {
	c := len(cols)
	m := Create(len(mat), c)
	for i := 0; i < len(mat); i++ {
		for j := 0; j < c; j++ {
			m[i][j] = mat[i][cols[j]]
		}
	}
	return m
}

func (mat Mat) ColsT(cols []int, row int) Mat {
	c := len(cols)
	m := Create(row, c)
	for i := 0; i < row; i++ {
		for j := 0; j < c; j++ {
			m[i][j] = mat[i][cols[j]]
		}
	}
	return m
}

func (mat Mat) Colst(cols []int, row int) Vec {
	c := len(cols)
	v := make(Vec, c)
	for j := 0; j < c; j++ {
		v[j] = mat[row][cols[j]]
	}
	return v
}

func (mat Mat) Match(cols []int, vals Mat, r int) bool {
	for i, c := range(cols) {
		if mat[r][c] != vals[r][i] {
			return false
		}
	}
	return true
}

func (mat Mat) Dims() (int, int) {
	return len(mat), len(mat[0])
}

func StringToVec(s string) Vec {
	var v []int
	dec := json.NewDecoder(strings.NewReader(s))
	err := dec.Decode(&v)
	if err != nil {
		panic(err)
	}
	return v
}

func StringToMat(s string) Mat {
	var m [][]int
	dec := json.NewDecoder(strings.NewReader(s))
	err := dec.Decode(&m)
	if err != nil {
		panic(err)
	}
	return m
}

func (mat Mat) Print() {
	r, _ := mat.Dims()
	for i := 0; i < r; i++ {
		fmt.Println(mat[i])
	}
}
