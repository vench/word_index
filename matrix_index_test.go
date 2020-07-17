package word_index

import (
	"testing"
)

func TestNewMatrixIndex(t *testing.T) {

	documents := []string{
		`abc zyz`,
		`test best aaa`,
		`anna vera zoom`,
		`aaa zet zzzz`,
	}

	index := NewMatrixIndex()
	err := index.Fit(documents...)
	if err != nil {
		t.Fatal(err)
	}

	for _, document := range documents {
		result := index.Find(document)
		if len(result) == 0 {
			t.Fatalf(``)
		}
	}

	result := index.Find(`aaa`)
	if len(result) != 2 {
		t.Fatalf(``)
	}
	if result[0] != 1 {
		t.Fatalf(``)
	}
	if result[1] != 3 {
		t.Fatalf(``)
	}

	result = index.Find(`vera`)
	if len(result) != 1 {
		t.Fatalf(``)
	}
	if result[0] != 2 {
		t.Fatalf(``)
	}

}

func TestNewMatrixIndex_MergeOrderArray(t *testing.T) {
	res := [][]int{
		{},
		{1,2,4,5},
		{2,3,6,7,8},
		{7,9,10},
		{1,1,1,1,1,1,1,2,2,2,2,2},
	}
	ok := []int {
		1,2,3,4,5,6,7,8,9,10,
	}
	r := MergeOrderedArray(res)
	if len(r) != len(ok) {
		t.Fatalf(`len(r) != len(ok)`)
	}
	for j := 0; j < len(ok); j ++ {
		if ok[j] != r[j] {
			t.Fatalf(``)
		}
	}

	//
	res = [][]int{
		{},
		{2,2,2,2,2,2},
		{1,1,1,1,1,1,1,2,2,2,2,2},
	}
	ok = []int {
		1,2,
	}
	r = MergeOrderedArray(res)
	if len(r) != len(ok) {
		t.Fatalf(`len(r) != len(ok)`)
	}
	for j := 0; j < len(ok); j ++ {
		if ok[j] != r[j] {
			t.Fatalf(``)
		}
	}
}