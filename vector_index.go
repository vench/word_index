package word_index

import (
	"fmt"
	"math"
	"sort"
)

type Vector struct {
	Id   uint32
	V    []float64
	Data interface{}
}

func (v *Vector) DistCos(a *Vector) float64 {
	return distCos(a.V, v.V)
}

func (v *Vector) DistMonteCarlo(a *Vector) float64 {
	return distMonteCarlo(a.V, v.V)
}

func (v *Vector) DistEuclidean(a *Vector) float64 {
	return distEuclidean(a.V, v.V)
}

func NewEmptyVector(id uint32, size int) *Vector {
	return &Vector{
		Id: id,
		V:  make([]float64, size),
	}
}

func NewVector(id uint32, v []float64, data interface{}) *Vector {
	return &Vector{
		Id:   id,
		V:    v,
		Data: data,
	}
}

type indexVectorItem struct {
	i         *Vector
	z         uint64
	neighbors []*indexVectorItem
}

type IndexVector struct {
	itemsMap           map[uint32]*indexVectorItem
	itemsOrderZ        []*indexVectorItem
	neighborsThreshold float64
}

func (iv *IndexVector) Fit(list []*Vector) error {
	items := make([]*indexVectorItem, len(list))
	itemsMap := make(map[uint32]*indexVectorItem)
	for i, v := range list {
		item := &indexVectorItem{
			i: v,
			z: ZOrderCurveFloat64(v.V),
		}
		items[i] = item
		itemsMap[item.i.Id] = item
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].z < items[j].z
	})


	if iv.neighborsThreshold != 0 {
		// update neighbors O(N^2)
		for i, v := range itemsMap {
			v.neighbors = make([]*indexVectorItem, 0)
			for j, v1 := range itemsMap {
				if i == j {
					continue
				}
				// TODO set sist type
				if v.i.DistEuclidean(v1.i) <= iv.neighborsThreshold {
					v.neighbors = append(v.neighbors, v1)
				}
			}
		}
	}


	iv.itemsMap = itemsMap
	iv.itemsOrderZ = items

	return nil
}

func (iv *IndexVector) SearchNeighborhood(v []float64, neighborhood []float64) ([]*Vector, error) {
	zSearch := ZOrderCurveFloat64(v)
	zNeighborhood := ZOrderCurveFloat64(neighborhood)
	zSearchLow := uint64(0)
	if zSearch > zNeighborhood {
		zSearchLow = zSearch - zNeighborhood
	}
	zSearchHigh := zSearch + zNeighborhood
	low := 0
	high := len(iv.itemsOrderZ) - 1
	for low <= high {
		median := (low + high) / 2
		if iv.itemsOrderZ[median].z < zSearchLow {
			low = median + 1
		} else {
			high = median - 1
		}
	}
	result := make([]*Vector, 0)
	for low < len(iv.itemsOrderZ) && iv.itemsOrderZ[low].z <= zSearchHigh {
		//fmt.Println(iv.itemsOrderZ[low].i.Id)
		result = append(result, iv.itemsOrderZ[low].i)
		low++
	}
	return result, nil
}

func (iv *IndexVector) Search(v []float64) ([]*Vector, error) {
	zSearch := ZOrderCurveFloat64(v)
	low := 0
	high := len(iv.itemsOrderZ) - 1
	for low <= high {
		median := (low + high) / 2
		if iv.itemsOrderZ[median].z < zSearch {
			low = median + 1
		} else {
			high = median - 1
		}
	}
	result := make([]*Vector, 0)
	for low < len(iv.itemsOrderZ) && iv.itemsOrderZ[low].z <= zSearch {
		//fmt.Println(iv.itemsOrderZ[low].i.Id)
		result = append(result, iv.itemsOrderZ[low].i)
		low++
	}
	return result, nil
}

func NewIndexVector() (*IndexVector, error) {
	return &IndexVector{}, nil
}

func ZOrderCurveFloat64(vec []float64) uint64 {
	v := make([]uint64, len(vec))
	for i, x := range vec {
		v[i] = zOrderCurveFloat64ToUint64(x)
	}
	return ZOrderCurve(v)
}

func zOrderCurveFloat64ToUint64(x float64) uint64 {
	return uint64(x * 1000000)
}

func ZOrderCurve(vec []uint64) uint64 {
	B := []uint64{0x00000000FFFFFFFF, 0x0000FFFF0000FFFF, 0x00FF00FF00FF00FF, 0x0F0F0F0F0F0F0F0F, 0x3333333333333333, 0x5555555555555555}
	S := []uint64{32, 16, 8, 4, 2, 1}

	for i := 0; i < len(S); i++ {
		for j := 0; j < len(vec); j++ {
			vec[j] = (vec[j] | (vec[j] << S[i])) & B[i]
		}
	}
	r := uint64(0)
	for i, v := range vec {
		r |= v << i
	}
	return r
}

func distCos(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}
	as, bs, ab := float64(0), float64(0), float64(0)
	for i := 0; i < len(a); i++ {
		as += a[i] * a[i]
		bs += b[i] * b[i]
		ab += a[i] * b[i]
	}

	if as == 0 || bs == 0 {
		return 0.0
	}
	return ab / (math.Sqrt(as) * math.Sqrt(bs))
}

func distEuclidean(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}
	s := float64(0)
	for i := 0; i < len(a); i++ {
		s += math.Pow(a[i]-b[i], 2)
	}
	return math.Sqrt(s)
}

func distMonteCarlo(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}
	s := float64(0)
	return s
}