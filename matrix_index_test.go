package word_index

import (
	"fmt"
	"testing"
)

func TestMatrixIndex(t *testing.T) {
	bi := NewMatrixIndex()
	tIndexPlainText(t, bi)

	bi = NewMatrixIndex()
	tIndexMathText(t, bi)
}

func TestNewMatrixIndex(t *testing.T) {

	documents := []string{
		`abc zyz`,
		`test best aaa`,
		`anna vera zoom`,
		`aaa zet zzzz`,
		`This site can’t be reached`,
	}

	index := NewMatrixIndex()
	err := index.Fit(documents...)
	if err != nil {
		t.Fatal(err)
	}

	for _, document := range documents {
		result := index.Query(document)
		if len(result) == 0 {
			t.Fatalf(``)
		}
	}

	result := index.Query(`aaa`)
	if len(result) != 2 {
		t.Fatalf(`%v`, result)
	}
	if result[0] != 1 {
		t.Fatalf(``)
	}
	if result[1] != 3 {
		t.Fatalf(``)
	}

	result = index.Query(`vera`)
	if len(result) != 1 {
		t.Fatalf(``)
	}
	if result[0] != 2 {
		t.Fatalf(``)
	}

	result = index.Query(`reac`)
	if len(result) != 0 {
		t.Fatalf(``)
	}
	result = index.Query(`reac*`)
	if len(result) != 1 {
		t.Fatalf(``)
	}
	if result[0] != 4 {
		t.Fatalf(``)
	}
}

func TestNewMatrixIndex_MergeOrderedArrayAnd(t *testing.T) {
	res := [][]int{
		{1, 2, 3, 4},
		{1, 3},
		{1, 2, 3, 4, 5, 6},
	}
	ret := MergeOrderedArrayAnd(res)
	if len(ret) != 2 {
		t.Fatalf(`len(ret) != 2, is %d`, len(ret))
	}
	if ret[0] != 1 {
		t.Fatalf(``)
	}
	if ret[1] != 3 {
		t.Fatalf(``)
	}

	//
	res = [][]int{
		{1, 2, 4, 7, 8, 11},
		{1, 3, 7, 10},
		{1, 2, 3, 4, 5, 6, 7, 10},
	}
	ret = MergeOrderedArrayAnd(res)
	if len(ret) != 2 {
		t.Fatalf(`len(ret) != 2, is %d`, len(ret))
	}
	if ret[0] != 1 {
		t.Fatalf(``)
	}
	if ret[1] != 7 {
		t.Fatalf(``)
	}

	//
	res = [][]int{
		{0, 2, 4, 7, 8, 11},
		{1, 3, 7, 10},
		{1, 2, 3, 4, 5, 6, 9, 10},
	}
	ret = MergeOrderedArrayAnd(res)
	if len(ret) != 0 {
		t.Fatalf(`len(ret) != 0, is %d`, len(ret))
	}

	//
	res = [][]int{
		{0, 2, 4},
		{2, 3, 7, 10},
		{1, 2, 3, 4, 5, 6, 9, 10},
	}
	ret = MergeOrderedArrayAnd(res)
	if len(ret) != 1 {
		t.Fatalf(`len(ret) != 0, is %d`, len(ret))
	}
	if ret[0] != 2 {
		t.Fatalf(``)
	}
}

func TestNewMatrixIndex_MergeOrderArray(t *testing.T) {
	res := [][]int{
		{},
		{1, 2, 4, 5},
		{2, 3, 6, 7, 8},
		{7, 9, 10},
		{1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2},
	}
	ok := []int{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
	}
	r := MergeOrderedArray(res)
	if len(r) != len(ok) {
		fmt.Println(r, ok)
		t.Fatalf(`len(r) != len(ok)`)
	}
	for j := 0; j < len(ok); j++ {
		if ok[j] != r[j] {
			t.Fatalf(``)
		}
	}

	//
	res = [][]int{
		{},
		{2, 2, 2, 2, 2, 2},
		{1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2},
	}
	ok = []int{
		1, 2,
	}
	r = MergeOrderedArray(res)
	if len(r) != len(ok) {
		t.Fatalf(`len(r) != len(ok)`)
	}
	for j := 0; j < len(ok); j++ {
		if ok[j] != r[j] {
			t.Fatalf(``)
		}
	}
}

func TestMatrixIndex_QueryAndOr(t *testing.T) {
	documents := []string{
		`abc zyz test site`,
		`test best aaa`,
		`anna vera zoom site`,
		`aaa zet zzzz`,
		`This site test can’t be reached`,
	}

	index := NewMatrixIndex()
	err := index.Fit(documents...)
	if err != nil {
		t.Fatal(err)
	}

	result := index.QueryAndOr(`test`, true)
	if len(result) != 3 {
		t.Fatalf(``)
	}
	if result[0] != 0 {
		t.Fatalf(``)
	}
	if result[1] != 1 {
		t.Fatalf(``)
	}
	if result[2] != 4 {
		t.Fatalf(``)
	}
	result = index.QueryAndOr(`site`, true)
	if len(result) != 3 {
		t.Fatalf(``)
	}
	if result[0] != 0 {
		t.Fatalf(``)
	}
	if result[1] != 2 {
		t.Fatalf(``)
	}
	if result[2] != 4 {
		t.Fatalf(``)
	}
	result = index.QueryAndOr(`test site`, true)
	if len(result) != 2 {
		t.Fatalf(``)
	}
	if result[0] != 0 {
		t.Fatalf(``)
	}
	if result[1] != 4 {
		t.Fatalf(``)
	}
	result = index.QueryAndOr(`test site`, false)
	if len(result) != 4 {
		t.Fatalf(``)
	}
	if result[0] != 0 {
		t.Fatalf(``)
	}
	if result[1] != 1 {
		t.Fatalf(``)
	}
	if result[2] != 2 {
		t.Fatalf(``)
	}
	if result[3] != 4 {
		t.Fatalf(``)
	}
}
