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
	return m.QueryAndOr(query, false)
}

func (m *MatrixIndex) QueryAndOr(query string, useAnd bool) []int {
	words := strings.Split(strings.ToLower(query), ` `)
	high := len(m.items) - 1
	results := make([][]int, len(words))
	for i, word := range words {
		q, variants := makeVariants(word)
		results[i] = m.findBin(q, variants, 0, high)
	}
	if useAnd {
		return MergeOrderedArrayAnd(results)
	}
	return MergeOrderedArray(results)
}

func (m *MatrixIndex) findBin(word string, variants []string, low, high int) []int {
	w := strings.TrimSpace(word)
	if len(w) < 2 {
		return []int{}
	}
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
		if m.items[median].word < w {
			low = median + 1
		} else {
			high = median - 1
		}
	}

	results := make([][]int, 0)
	for low < len(m.items) && m.compareWord(m.items[low].word, word, variants) {
		/*if len(result) == 0 {
				result = m.items[low].index
			} else {
				result = MergeOrderedArray([][]int{result, m.items[low].index})
		} */
		results = append(results, m.items[low].index)

		low++
	}

	return MergeOrderedArray(results)
}

func (m *MatrixIndex) compareWord(word, query string, variants []string) bool {
	if word == query {
		return true
	}
	if query[len(query)-1:] == tagAny {
		if word == query[:len(query)-1] {
			return true
		}
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
	maxValue := 0

	for j := 0; j < len(a); j++ {
		if len(a[j]) == 0 {
			a = append(a[:j], a[j+1:]...)
			continue
		}
		if len(a[j]) > maxLen {
			maxLen = len(a[j])
		}
		if maxValue < a[j][len(a[j])-1] {
			maxValue = a[j][len(a[j])-1]
		}
	}
	offsets := make([]int, len(a))
	maxValue++
	b := make([]int, 0, maxLen)
	lastIndex := -1
	minValue := maxValue
	for true {

		minIndexResult := -1
		for j := 0; j < len(a); j++ {
			if len(a[j]) > offsets[j] {
				if a[j][offsets[j]] < minValue {
					minValue = a[j][offsets[j]]
					minIndexResult = j
				}
			} else {
				a = append(a[:j], a[j+1:]...)
				offsets = append(offsets[:j], offsets[j+1:]...)
				j--
			}
		}
		if minIndexResult == -1 {
			break
		}
		if lastIndex < minValue {
			b = append(b, minValue)
			lastIndex = minValue
		}
		minValue = maxValue
		//a[minIndexResult] = a[minIndexResult][1:]
		offsets[minIndexResult]++
	}
	return b
}

func MergeOrderedArrayAnd(a [][]int) []int {
	b := make([]int, 0)
	minIndex := 0
	for i := 1; i < len(a); i++ {
		if len(a[minIndex]) > len(a[i]) {
			minIndex = i
		}
	}
	offsets := make([]int, len(a))
	for i, v := range a[minIndex] {
		_ = i
		has := true
		for j := 0; j < len(a); j++ {
			if j == minIndex {
				continue
			}
			for ; offsets[j] < len(a[j]); offsets[j]++ {
				if a[j][offsets[j]] > v {
					has = false
					break
				}
				if has = a[j][offsets[j]] == v; has {
					break
				}
			}
			if !has {
				break
			}
		}
		if has {
			b = append(b, v)
		}
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
