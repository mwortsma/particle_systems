package matutil

import (
	"encoding/json"
	"fmt"
	"strings"
	"math"
	"strconv"
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

func (mat Mat) AllColsT(row int) Mat {
	_,c := mat.Dims()
	m := Create(row, c)
	for i := 0; i < row; i++ {
		for j := 0; j < c; j++ {
			m[i][j] = mat[i][j]
		}
	}
	return m
}

func (mat Mat) ColsRange(cols []int, start int, stop int) Mat {
	c := len(cols)
	m := Create(stop-start, c)
	for i := start; i < stop; i++ {
		for j := 0; j < c; j++ {
			m[i-start][j] = mat[i][cols[j]]
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

func (mat Mat) Match(cols []int, vals Mat, t int, tau int) bool {

	if tau > 0 && t < len(mat) - tau {
		return true
	}
	var a int
	if tau <= 0 || t-tau < 0{
		a = t
	} else {
		a = t - (len(mat) - tau)
	}


	for i, c := range cols {
		if mat[t][c] != vals[a][i] {
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

func BinaryStrings(n int) []Vec {
	s := make([]Vec, int(math.Pow(2.0, float64(n))))
	for i := 0; i < int(math.Pow(2.0, float64(n))); i++ {
		b := []byte(string(strconv.FormatInt(int64(i), 2)))
		a := make(Vec, n)
		for j := 0; j < len(b); j++ {
			a[j] = int(b[len(b)-j-1]-48)
		}
		s[i] = a
	}
	return s

}

func BinaryMats(r,c int) []Mat {
	l := int(math.Pow(2.0, float64(r*c)))
	fmt.Println(l)
	s := make([]Mat, l)
	strings := BinaryStrings(r*c)
	for k, str := range strings {
		s[k] = Create(r,c)
		for i := 0; i < r; i ++ {
			for j := 0; j < c; j++ {
				s[k][i][j] = str[c*i + j]
			}
		}
	}
	return s
}

func Concat(m1 Mat, m2 Mat) Mat {
	rows, cols := m1.Dims()
	_, cols2 := m2.Dims()

	new_mat := Create(rows, cols + cols2)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			new_mat[i][j] = m1[i][j] 
		}
	}
	for i := 0; i < rows; i++ {
		for j := 0; j < cols2; j++ {
			new_mat[i][j+cols] = m2[i][j] 
		}
	}
	return new_mat
}
