package word_index

import (
	"fmt"
	"math"
	"sort"
	"testing"
)

func TestVector_DistCos(t *testing.T) {
	tests := []struct {
		V1   []float64
		V2   []float64
		Dist float64
	}{
		{V1: []float64{}, V2: []float64{}, Dist: 0},
		{V1: []float64{1}, V2: []float64{1}, Dist: 1},
		{V1: []float64{1, 2, 3}, V2: []float64{1, 2, 3}, Dist: 1},
		{V1: []float64{1, 1, 0}, V2: []float64{1, 1, 1}, Dist: 2 / math.Sqrt(2*3)},
		{V1: []float64{1, 1, 1}, V2: []float64{1, 1, 1}, Dist: 3 / math.Sqrt(3*3)},       // 3/sqrt(9)=1
		{V1: []float64{1, 1, 1}, V2: []float64{10, 10, 10}, Dist: 30 / math.Sqrt(300*3)}, // 30/sqrt(900)=1
	}
	for i, test := range tests {
		d := distCos(test.V1, test.V2)
		if fmt.Sprintf(`%.4f`, test.Dist) != fmt.Sprintf(`%.4f`, d) {
			t.Fatalf(`N: %d, not equal dist: %.4f != %.4f`, i, test.Dist, d)
		}
	}

}

func TestVector_DistMonteCarlo(t *testing.T) {
	// TODO
}

func TestVector_DistEuclidean2(t *testing.T) {
	v1 := NewVector(1, []float64{0,0,0,0,0,0,0,0,0,0}, nil)
	v2 := NewVector(2, []float64{10,10,10,10,10,10,10,10,10,10}, nil)
	dist := v1.DistEuclidean(v2)
	fmt.Println(dist)
}

func TestNewIndexVector_SearchNeighborhood(t *testing.T) {
	iv, err := NewIndexVector()
	if err != nil {
		t.Fatalf(`%s`, err.Error())
	}
	// dim 1
	err = iv.Fit([]*Vector{
		{Id: 1, V: []float64{1}},
		{Id: 2, V: []float64{1}},
		{Id: 3, V: []float64{2}},
		{Id: 4, V: []float64{101}},
		{Id: 5, V: []float64{101}},
	})
	if err != nil {
		t.Fatalf(`%s`, err.Error())
	}
	list, err := iv.SearchNeighborhood([]float64{1}, []float64{0})
	if err != nil {
		t.Fatalf(`%s`, err.Error())
	}
	if len(list) != 2 {
		t.Fatalf(`len not equals 2, %d`, len(list))
	}

	list, err = iv.SearchNeighborhood([]float64{1}, []float64{math.Sqrt(2)})
	if err != nil {
		t.Fatalf(`%s`, err.Error())
	}
	if len(list) != 3 {
		t.Fatalf(`len not equals 3, %d`, len(list))
	}

	// dim 2
	err = iv.Fit([]*Vector{
		{Id: 1, V: []float64{1, 1}},
		{Id: 2, V: []float64{1, 2}},
		{Id: 3, V: []float64{2, 2}},
		{Id: 4, V: []float64{101, 100}},
		{Id: 5, V: []float64{101, 102}},
	})
	if err != nil {
		t.Fatalf(`%s`, err.Error())
	}
	list, err = iv.SearchNeighborhood([]float64{1, 1}, []float64{})
	if err != nil {
		t.Fatalf(`%s`, err.Error())
	}
	if len(list) != 1 {
		t.Fatalf(`len not equals 1, %d`, len(list))
	}
	list, err = iv.SearchNeighborhood([]float64{1, 1}, []float64{2, 2})
	if err != nil {
		t.Fatalf(`%s`, err.Error())
	}
	if len(list) != 3 {
		t.Fatalf(`len not equals 3, %d`, len(list))
	}
}

func TestIndexVector_Search(t *testing.T) {

	iv, err := NewIndexVector()
	if err != nil {
		t.Fatalf(`%s`, err.Error())
	}
	// dim 1
	err = iv.Fit([]*Vector{
		{Id: 1, V: []float64{1}},
		{Id: 2, V: []float64{1}},
		{Id: 3, V: []float64{2}},
		{Id: 4, V: []float64{101}},
		{Id: 5, V: []float64{101}},
	})
	if err != nil {
		t.Fatalf(`%s`, err.Error())
	}
	list, err := iv.Search([]float64{1})
	if err != nil {
		t.Fatalf(`%s`, err.Error())
	}
	if len(list) != 2 {
		t.Fatalf(`len not equals 2, %d`, len(list))
	}
	if list[0].Id != 1 {
		t.Fatalf(`id not equals 1, %d`, list[0].Id)
	}
	if list[1].Id != 2 {
		t.Fatalf(`id not equals 2, %d`, list[1].Id)
	}
	list, err = iv.Search([]float64{2})
	if err != nil {
		t.Fatalf(`%s`, err.Error())
	}
	if len(list) != 1 {
		t.Fatalf(`len not equals 1, %d`, len(list))
	}
	if list[0].Id != 3 {
		t.Fatalf(`id not equals 3, %d`, list[0].Id)
	}
	list, err = iv.Search([]float64{3})
	if err != nil {
		t.Fatalf(`%s`, err.Error())
	}
	if len(list) != 0 {
		t.Fatalf(`query is not empty`)
	}
}

func TestVector_DistEuclidean(t *testing.T) {
	tests := []struct {
		V1   []float64
		V2   []float64
		Dist float64
	}{
		{V1: []float64{}, V2: []float64{}, Dist: 0},
		{V1: []float64{1}, V2: []float64{1}, Dist: 0},
		{V1: []float64{1, 2, 3}, V2: []float64{1, 2, 3}, Dist: 0},
		{V1: []float64{1}, V2: []float64{2}, Dist: 1},
		{V1: []float64{1, 1}, V2: []float64{2, 2}, Dist: math.Sqrt(2)},
		{V1: []float64{1, 1}, V2: []float64{2, 8}, Dist: math.Sqrt(1 + 49)},
		{V1: []float64{1, 1, 0}, V2: []float64{1, 1, 1}, Dist: 1},
		{V1: []float64{1, 1, 1}, V2: []float64{1, 1, 1}, Dist: 0},
		{V1: []float64{1, 1, 1}, V2: []float64{10, 10, 10}, Dist: 15.5885},
	}
	for i, test := range tests {
		d := distEuclidean(test.V1, test.V2)
		//d2 := distCos(test.V1, test.V2)
		if fmt.Sprintf(`%.4f`, test.Dist) != fmt.Sprintf(`%.4f`, d) {
			t.Fatalf(`N: %d, not equal dist: %.4f != %.4f`, i, test.Dist, d)
		}
		//fmt.Println(i, d,d2)
	}

}

func TestZOrderCurveFloat64_locale(t *testing.T) {
	type pointZ struct {
		n     int
		point []float64
		z     uint64
	}
	// 64 dim
	points := []pointZ{
		{n: 2, point: []float64{0.11, 0.11, 0.12, 0.13, 0.11, 0.11, 0.12, 0.13, 0.11, 0.11, 0.12, 0.13, 0.11, 0.11, 0.12, 0.13, 0.11, 0.11, 0.12, 0.13, 0.11, 0.11, 0.12, 0.13, 0.11, 0.11, 0.12, 0.13, 0.11, 0.11, 0.12, 0.13, 0.11, 0.11, 0.12, 0.13, 0.11, 0.11, 0.12, 0.13, 0.11, 0.11, 0.12, 0.13, 0.11, 0.11, 0.12, 0.13, 0.11, 0.11, 0.12, 0.13, 0.11, 0.11, 0.12, 0.13, 0.11, 0.11, 0.12, 0.13, 0.11, 0.11, 0.12, 0.13}},
		{n: 4, point: []float64{0.000101, 0.000101, 0.0003, 0.000101, 0.000101, 0.000101, 0.0003, 0.000101, 0.000101, 0.000101, 0.0003, 0.000101, 0.000101, 0.000101, 0.0003, 0.000101, 0.000101, 0.000101, 0.0003, 0.000101, 0.000101, 0.000101, 0.0003, 0.000101, 0.000101, 0.000101, 0.0003, 0.000101, 0.000101, 0.000101, 0.0003, 0.000101, 0.000101, 0.000101, 0.0003, 0.000101, 0.000101, 0.000101, 0.0003, 0.000101, 0.000101, 0.000101, 0.0003, 0.000101, 0.000101, 0.000101, 0.0003, 0.000101, 0.000101, 0.000101, 0.0003, 0.000101, 0.000101, 0.000101, 0.0003, 0.000101, 0.000101, 0.000101, 0.0003, 0.000101, 0.000101, 0.000101, 0.0003, 0.000101}},
		{n: 1, point: []float64{0.10, 10, 0.10, 0.12, 0.10, 10, 0.10, 0.12, 0.10, 10, 0.10, 0.12, 0.10, 10, 0.10, 0.12, 0.10, 10, 0.10, 0.12, 0.10, 10, 0.10, 0.12, 0.10, 10, 0.10, 0.12, 0.10, 10, 0.10, 0.12, 0.10, 10, 0.10, 0.12, 0.10, 10, 0.10, 0.12, 0.10, 10, 0.10, 0.12, 0.10, 10, 0.10, 0.12, 0.10, 10, 0.10, 0.12, 0.10, 10, 0.10, 0.12, 0.10, 10, 0.10, 0.12, 0.10, 10, 0.10, 0.12}},
		{n: 3, point: []float64{0.000100, 0.000100, 0.0003, 0.000100, 0.000100, 0.000100, 0.0003, 0.000100, 0.000100, 0.000100, 0.0003, 0.000100, 0.000100, 0.000100, 0.0003, 0.000100, 0.000100, 0.000100, 0.0003, 0.000100, 0.000100, 0.000100, 0.0003, 0.000100, 0.000100, 0.000100, 0.0003, 0.000100, 0.000100, 0.000100, 0.0003, 0.000100, 0.000100, 0.000100, 0.0003, 0.000100, 0.000100, 0.000100, 0.0003, 0.000100, 0.000100, 0.000100, 0.0003, 0.000100, 0.000100, 0.000100, 0.0003, 0.000100, 0.000100, 0.000100, 0.0003, 0.000100, 0.000100, 0.000100, 0.0003, 0.000100, 0.000100, 0.000100, 0.0003, 0.000100, 0.000100, 0.000100, 0.0003, 0.000100}},
	}

	for i, point := range points {
		points[i].z = ZOrderCurveFloat64(point.point)
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i].z < points[j].z
	})

	for i := 0; i < len(points); i++ {
		if points[i].n != i+1 {
			t.Fatalf(`N_%d != i_%d`, points[i].n, i+1)
		}
	}

}

func TestZOrderCurveFloat64(t *testing.T) {

	type itemZ struct {
		index int
		z     uint64
	}
	items1 := make([]itemZ, 0)
	items2 := make([]itemZ, 0)

	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			v := []uint64{
				uint64(x),
				uint64(y),
			}
			z1 := ZOrderCurve(v)
			items1 = append(items1, itemZ{
				index: x*8 + y,
				z:     z1,
			})

			vf := []float64{
				float64(x),
				float64(y),
			}
			z2 := ZOrderCurveFloat64(vf)
			items2 = append(items2, itemZ{
				index: x*8 + y,
				z:     z2,
			})
			//fmt.Println(z1, ` `, z2, ` i `, x*8+y, x, y)
		}
	}

	sort.Slice(items1, func(i, j int) bool {
		return items1[i].z < items1[j].z
	})
	sort.Slice(items2, func(i, j int) bool {
		return items2[i].z < items2[j].z
	})

	i := 0
	if items1[i].index != items2[i].index {
		t.Errorf(`index1(%d) != index2(%d), n: %d`, items1[i].index, items2[i].index, 0)
	}
	i = len(items1) - 1
	if items1[i].index != items2[i].index {
		t.Errorf(`index1(%d) != index2(%d), n: %d`, items1[i].index, items2[i].index, 0)
	}
}
