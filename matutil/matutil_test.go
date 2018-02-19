package matutil

import (
	//"fmt"
	"github.com/mwortsma/particle_systems/matutil"
	"testing"
)

func TestMat(t *testing.T) {
	r, c := 4, 5
	m := matutil.Create(r, c)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			m[i][j] = i + j
		}
	}
	t.Log(m.String())
	t.Log(m.Col(0).String())
	t.Log(m.Col(1).String())
	t.Log(m.Cols([]int{2, 0}).String())
	t.Log(matutil.StringToMat(m.ColsT([]int{2, 0}, 2).String()))
	t.Log(matutil.StringToVec(m.Col(1).String()).String())
}


func TestQ(t *testing.T) {
	t.Log(matutil.QStrings(2,4))
}
