package word_index

import (
	"sort"
	"strings"
	"sync"
)

const (
	tagAny     = `*`
	tagAnyRune = '*'
	emptyFind  = -1
)

//
type Index interface {
	Find(string) int
	FindOff(string, int) int
	FindAll(string) []int
	FindAt(int, string) bool
	Add(...string)
	DocumentAt(int) (string, bool)
}

//
type variant struct {
	query    string
	variants []string
}

//
type indexItem struct {
	words    []string
	document string
}

func (i *indexItem) findInterpolation(query string, variants []string) bool {

	if len(query) == 0 {
		return false
	}

	var (
		mid  int
		low  = 0
		high = len(i.words) - 1
	)

	for i.words[low][0] < query[0] && i.words[high][0] > query[0] {
		mid = low + (int(query[0]-i.words[low][0])*(high-low))/int(i.words[high][0]-i.words[low][0])

		if i.words[mid] < query {
			low = mid + 1
		} else if i.words[mid] > query {
			high = mid - 1
		} else {
			return true
		}
	}

	if i.words[low][0] == query[0] {
		for n := low; n < len(i.words); n++ {
			if i.words[n] == query {
				return true
			} else if query[len(query)-1:] == tagAny {
				for n := 0; n < len(query); n++ {
					r := query[n]
					if r == tagAnyRune {
						return true
					} else if len(i.words[low]) <= n || i.words[low][n] != r {
						break
					}
				}
			} else if len(variants) > 0 {
				for _, variant := range variants {
					if i.words[n] == variant {
						return true
					}
				}
			}
		}
	}

	if i.words[high][0] == query[0] {
		for n := high; n < len(i.words); n++ {
			if i.words[n] == query {
				return true
			} else if query[len(query)-1:] == tagAny {
				for j, r := range []rune(query) {
					if r == tagAnyRune {
						return true
					} else if len(i.words[low]) <= n || rune(i.words[low][j]) != r {
						break
					}
				}
			} else if len(variants) > 0 {
				for _, variant := range variants {
					if i.words[n] == variant {
						return true
					}
				}
			}
		}
	}

	return false
}

//
func (i *indexItem) findBin(query string, variants []string) bool {

	if len(query) == 0 {
		return false
	}

	low := 0
	high := len(i.words) - 1

	for low <= high {
		median := (low + high) / 2
		if i.words[median][0] < query[0] {
			low = median + 1
		} else {
			high = median - 1
		}
	}

	for low < len(i.words) && i.words[low][0] == query[0] {
		if i.words[low] == query {
			return true
		} else if query[len(query)-1:] == tagAny {
			for n := 0; n < len(query); n++ {
				r := query[n]
				if r == tagAnyRune {
					return true
				} else if len(i.words[low]) <= n || i.words[low][n] != r {
					break
				}
			}
		} else if len(variants) > 0 {
			for _, variant := range variants {
				if i.words[low] == variant {
					return true
				}
			}
		}
		low++
	}
	return false
}

//
type indexWord struct {
	data      []*indexItem
	binSearch bool
}

func (i *indexWord) FindAll(str string) []int {
	words := strings.Split(strings.ToLower(str), ` `)
	variants := make([]*variant, len(words))
	for n, word := range words {
		q, v := i.makeVariants(word)
		vr := &variant{query: q, variants: v}
		variants[n] = vr
	}

	result := make([]int, 0)
	var offset = 0
	for true {
		i := i.findOff(variants, offset)
		if i == emptyFind {
			break
		}
		result = append(result, i)
		offset = i + 1
	}
	return result
}

func (i *indexWord) FindOff(str string, offset int) int {
	words := strings.Split(strings.ToLower(str), ` `)
	variants := make([]*variant, len(words))
	for n, word := range words {
		q, v := i.makeVariants(word)
		vr := &variant{query: q, variants: v}
		variants[n] = vr
	}
	return i.findOff(variants, offset)
}

func (i *indexWord) findOff(variants []*variant, offset int) int {


	for index := offset; index < len(i.data); index++ {
		d := i.data[index]

		for _, v := range variants {
			if i.binSearch {
				if ok := d.findBin(v.query, v.variants); ok {
					return index
				}
			} else {
				if ok := d.findInterpolation(v.query, v.variants); ok {
					return index
				}
			}
		}
	}

	return emptyFind
}

//
func (i *indexWord) makeVariants(q string) (string, []string) {
	variants := make([]string, 0)

	if len(q) > 0 && q[len(q)-1] == ')' {
		base := make([]rune, 0)
		start := false
		variant := make([]rune, 0)
		for _, r := range []rune(q) {
			if r == tagAnyRune {
				q = string(append(base, r))
				variants = make([]string, 0)
				break
			}
			if r == ')' {
				variants = append(variants, string(variant))
				break
			} else if r == '(' {
				start = true
				variant = append(variant, base...)
				variants = append(variants, string(variant))
			} else if start && r == '|' {
				variants = append(variants, string(variant))
				variant = make([]rune, 0)
				variant = append(variant, base...)
			} else if start {
				variant = append(variant, r)
			} else {
				base = append(base, r)
			}
		}
	}
	return q, variants
}

//
func (i *indexWord) Add(str ...string) {
	for _, s := range str {
		words := strings.Split(strings.ToLower(s), ` `)

		k := 0
		for k < len(words) {
			if len(words[k]) == 0 {
				words = append(words[:k], words[k+1:]...)
			} else {
				k++
			}
		}

		sort.Slice(words, func(i, j int) bool {
			if words[i] < words[j] {
				return true
			}
			return false
		})

		n := indexItem{words: words, document: s}
		i.data = append(i.data, &n)
	}
}

//
func (i *indexWord) Find(str string) int {
	return i.FindOff(str, 0)
}

//
func (i *indexWord) DocumentAt(index int) (string, bool) {
	if len(i.data) > index && index >= 0 {
		return i.data[index].document, true
	}
	return ``, false
}

//
func (i *indexWord) FindAt(index int, str string) bool {
	if index < 0 || len(i.data) < index {
		return false
	}
	words := strings.Split(strings.ToLower(str), ` `)
	for _, word := range words {
		query, variants := i.makeVariants(word)
		if i.binSearch {
			if ok := i.data[index].findBin(query, variants); ok {
				return true
			}
		} else {
			if ok := i.data[index].findInterpolation(query, variants); ok {
				return true
			}
		}

	}
	return false
}

//
func NewIndex() Index {
	return &indexWord{data: make([]*indexItem, 0), binSearch: true}
}

//
type indexWordSync struct {
	indexWord
	mx sync.RWMutex
}

func (i *indexWordSync) Add(str ...string) {
	i.mx.Lock()
	i.indexWord.Add(str...)
	i.mx.Unlock()
}

func (i *indexWordSync) Find(str string) int {
	i.mx.RLock()
	defer i.mx.RUnlock()
	return i.indexWord.Find(str)
}

func (i *indexWordSync) FindOff(str string, offset int) int {
	i.mx.RLock()
	defer i.mx.RUnlock()
	return i.indexWord.FindOff(str, offset)
}

func (i *indexWordSync) DocumentAt(index int) (string, bool) {
	i.mx.RLock()
	defer i.mx.RUnlock()
	return i.indexWord.DocumentAt(index)
}

func (i *indexWordSync) FindAt(index int, str string) bool {
	i.mx.RLock()
	defer i.mx.RUnlock()
	return i.indexWord.FindAt(index, str)
}

func (i *indexWordSync) FindAll(str string) []int {
	i.mx.RLock()
	defer i.mx.RUnlock()
	return i.indexWord.FindAll(str)
}

//
func NewIndexSync() Index {
	return &indexWordSync{indexWord: indexWord{data: make([]*indexItem, 0), binSearch: true}}
}
