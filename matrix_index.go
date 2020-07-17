package word_index

import (
	"sort"
	"strings"
)

type MatrixIndex struct {
	items []*matrixIndexItem
}

func (m *MatrixIndex) Fit(documents ...string) error  {

	mWords := make(map[string]map[int]struct{})
	for index,document := range documents {
		words := strings.Split(strings.ToLower(document), ` `)
		for _, word := range words {
			w,ok := mWords[word]
			if !ok {
				w = make(map[int]struct{})
			}
			w[index] = struct{}{}
			mWords[word] = w
		}
	}

	items := make([]*matrixIndexItem, len(mWords))
	i := 0
	for word,index := range mWords {

		item := &matrixIndexItem{word: word, index: make([]int, len(index)) }
		j := 0
		for inx,_ := range index {
			item.index[j] = inx
			j ++
		}
		items[i] = item
		i ++
	}
	
	sort.Slice(items, func(i, j int) bool {
		return items[i].word < items[j].word
	})

	m.items = items
	return nil
}

func (m*MatrixIndex) Find(query string) []int {

	words := strings.Split(strings.ToLower(query), ` `)
	high :=  len(m.items) - 1
	results := make([][]int, len(words))
	for i,word := range words {
		results[i] = m.findBin(word, 0, high)
	}

	return MergeOrderedArray(results)
}

func (m*MatrixIndex) findBin(word string, low, high int) []int {
	for low <= high {
		median := (low + high) / 2
		if m.items[median].word < word {
			low = median + 1
		} else {
			high = median - 1
		}
	}

	result := make([]int, 0)
	for low <  len(m.items) && m.items[low].word == word { // test comapere word
		result = append(result, m.items[low].index...)
		low++
	}

	return result
}

func MergeOrderedArray(a [][]int) []int {
	maxLen := 0
	maxIndex := 0
	for j := 0; j < len(a); j ++ {
		if len(a[j]) == 0 {
			a = append(a[:j], a[j+1:]...)
			continue
		}
		if len(a[j]) > maxLen {
			maxLen = len(a[j])
		}
		if maxIndex < a[j][len( a[j])-1] {
			maxIndex = a[j][len( a[j])-1]
		}
	}
	maxIndex ++
	b := make([]int, 0, maxLen)
	lastIndex := -1
	minIndex := maxIndex
	for true {

		minIndexResult := -1
		for j := 0; j < len(a); j ++ {
			if len(a[j]) > 0 {
				if a[j][0] < minIndex {
					minIndex = a[j][0]
					minIndexResult = j
				}
			} else {
				a = append(a[:j], a[j+1:]...)
				j --
			}
		}
		if minIndexResult == -1 {
			break
		}
		if lastIndex < minIndex {
			b = append(b, minIndex)
			lastIndex = minIndex
		}
		minIndex = maxIndex
		a[minIndexResult] = a[minIndexResult][1:]
	}
	return b
}

type matrixIndexItem struct {
	word string
	index []int
}

func NewMatrixIndex() *MatrixIndex {
	return &MatrixIndex{}
}