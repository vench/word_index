package word_index

import (
	"sort"
	"strings"
)

type MatrixIndex struct {
	items     []*matrixIndexItem
	documents []string
}

func (m *MatrixIndex) Find(query string) int {
	result := m.Query(query)
	if len(result) > 0 {
		return result[0]
	}
	return emptyFind
}

func (m *MatrixIndex) FindOff(query string, low int) int {
	if m.FindAt(low, query) {
		return low
	}
	return emptyFind
}

func (m *MatrixIndex) FindAll(query string) []int {
	return m.Query(query)
}

func (m *MatrixIndex) FindAt(index int, query string) bool {
	result := m.Query(query)
	if len(result) == 0 {
		return false
	}

	low, high := 0, len(result)-1
	for low <= high {
		median := (low + high) / 2
		if result[median] < index {
			low = median + 1
		} else {
			high = median - 1
		}
	}
	return result[low] == index
}

func (m *MatrixIndex) Add(documents ...string) {
	documents = append(m.documents, documents...)
	m.Fit(documents...)
}

func (m *MatrixIndex) DocumentAt(index int) (string, bool) {
	if len(m.documents) > index {
		return m.documents[index], true
	}
	return "", false
}

func (m *MatrixIndex) Fit(documents ...string) error {

	mWords := make(map[string]map[int]struct{})
	for index, document := range documents {
		words := strings.Split(strings.ToLower(document), ` `)
		for _, word := range words {
			w, ok := mWords[word]
			if !ok {
				w = make(map[int]struct{})
			}
			w[index] = struct{}{}
			mWords[word] = w
		}
	}

	items := make([]*matrixIndexItem, len(mWords))
	i := 0
	for word, index := range mWords {

		item := &matrixIndexItem{word: word, index: make([]int, len(index))}
		j := 0
		for inx, _ := range index {
			item.index[j] = inx
			j++
		}

		sort.Slice(item.index, func(i, j int) bool {
			return item.index[i] < item.index[j]
		})

		items[i] = item
		i++
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].word < items[j].word
	})

	m.items = items
	m.documents = documents
	return nil
}

func (m *MatrixIndex) Query(query string) []int {
	words := strings.Split(strings.ToLower(query), ` `)
	high := len(m.items) - 1
	results := make([][]int, len(words))
	for i, word := range words {
		q, variants := makeVariants(word)
		results[i] = m.findBin(q, variants, 0, high)
	}

	return MergeOrderedArray(results)
}

func (m *MatrixIndex) findBin(word string, variants []string, low, high int) []int {
	w := word
	if w[len(w)-1] == tagAnyRune {
		w = w[:len(w)-1]
	} else if w[len(w)-1] == ')' {
		for i := len(w) - 1; i >= 0; i-- {
			if w[i] == '(' {
				w = w[:i]
				break
			}
		}
	}
	for low <= high {
		median := (low + high) / 2
		if m.items[median].word[0] < w[0] {
			low = median + 1
		} else {
			high = median - 1
		}
	}

	result := make([]int, 0) //  len(m.items[low].word) >= len(w) && m.items[low].word[len(w)-1] == word[len(w)-1]
	for low < len(m.items) && m.items[low].word[0] == word[0] {
		if m.compareWord(m.items[low].word, word, variants) {
			result = append(result, m.items[low].index...)
		}
		low++
	}

	return result
}

func (m *MatrixIndex) compareWord(word, query string, variants []string) bool {
	if word == query {
		return true
	}
	if query[len(query)-1:] == tagAny {
		for n := 0; n < len(query); n++ {
			r := query[n]
			if r == tagAnyRune {
				return true
			} else if len(word) <= n || word[n] != r {
				break
			}
		}
	}
	if len(variants) > 0 {
		for _, variant := range variants {
			if word == variant {
				return true
			}
		}
	}
	return false
}

func MergeOrderedArray(a [][]int) []int {
	maxLen := 0
	maxIndex := 0
	for j := 0; j < len(a); j++ {
		if len(a[j]) == 0 {
			a = append(a[:j], a[j+1:]...)
			continue
		}
		if len(a[j]) > maxLen {
			maxLen = len(a[j])
		}
		if maxIndex < a[j][len(a[j])-1] {
			maxIndex = a[j][len(a[j])-1]
		}
	}
	maxIndex++
	b := make([]int, 0, maxLen)
	lastIndex := -1
	minIndex := maxIndex
	for true {

		minIndexResult := -1
		for j := 0; j < len(a); j++ {
			if len(a[j]) > 0 {
				if a[j][0] < minIndex {
					minIndex = a[j][0]
					minIndexResult = j
				}
			} else {
				a = append(a[:j], a[j+1:]...)
				j--
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
	word  string
	index []int
}

func NewMatrixIndex() *MatrixIndex {
	return &MatrixIndex{}
}
