package word_index

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
)

const (
	tagAny          = `*`
	tagAnyRune rune = '*'
	emptyFind       = -1
)

//
type Index interface {
	Find(string) int
	Add(...string)
	DocumentAt(int) (string, bool)
}

type indexRegexp struct {
	data []string
}

//
func (i *indexRegexp) Find(str string) int {
	words := strings.Split(strings.ToLower(str), ` `)
	for _, word := range words {
		word := strings.ReplaceAll(word, tagAny, `.*`)
		expr := fmt.Sprintf(`(^|\s)%s(\s|\.|\,|$)`, strings.TrimSpace(word))
		r, _ := regexp.Compile(expr)
		for n, dataStr := range i.data {
			if len(r.FindStringSubmatchIndex(dataStr)) > 0 {
				return n
			}
		}
	}
	return emptyFind
}

//
func (i *indexRegexp) Add(str ...string) {
	for i := 0; i < len(str); i++ {
		str[i] = strings.ToLower(str[i])
	}
	i.data = append(i.data, str...)
}

//
func (i *indexRegexp) DocumentAt(index int) (string, bool) {
	if len(i.data) > index && index >= 0 {
		return i.data[index], true
	}
	return ``, false
}

//
func NewIndexRegexp() Index {
	return &indexRegexp{data: []string{}}
}

//
type indexBinItem struct {
	words    []string
	document string
}

//
func (i *indexBinItem) find(query string, variants []string) int {

	if len(query) == 0 {
		return emptyFind
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
			return low
		} else if query[len(query)-1:] == tagAny {
			for n, r := range []rune(query[len(query)-1:]) {
				if r == tagAnyRune {
					return low
				} else if len(i.words[low]) <= n || rune(i.words[low][n]) != r {
					break
				}
			}
		} else if len(variants) > 0 {
			for _, variant := range variants {
				if i.words[low] == variant {
					return low
				}
			}
		}
		low++
	}
	return emptyFind
}

//
type indexBin struct {
	data []*indexBinItem
}

//
func (i *indexBin) makeVariants(q string) (string, []string) {
	variants := make([]string, 0)

	if q[len(q)-1] == ')' {
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
func (i *indexBin) Add(str ...string) {
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

		n := &indexBinItem{words: words, document: s}
		i.data = append(i.data, n)
	}
}

//
func (i *indexBin) Find(str string) int {
	words := strings.Split(strings.ToLower(str), ` `)
	for _, word := range words {

		query, variants := i.makeVariants(word)
		for _, d := range i.data {
			if index := d.find(query, variants); index >= 0 {
				return index
			}
		}
	}
	return emptyFind
}

//
func (i *indexBin) DocumentAt(index int) (string, bool) {
	if len(i.data) > index && index >= 0 {
		return i.data[index].document, true
	}
	return ``, false
}

//
func NewIndexBin() Index {
	return &indexBin{data: []*indexBinItem{}}
}


//
type indexBinSync struct {
	indexBin
	mx sync.RWMutex
}

func (i *indexBinSync) Add(str ...string) {
	i.mx.Lock()
	defer i.mx.Unlock()
	i.indexBin.Add(str...)
}

func (i*indexBinSync) Find(str string) int {
	i.mx.RLock()
	defer i.mx.RUnlock()
	return i.indexBin.Find(str)
}

func (i*indexBinSync) DocumentAt(index int) (string, bool) {
	i.mx.RLock()
	defer i.mx.RUnlock()
	return i.indexBin.DocumentAt(index)
}

//
func NewIndexBinSync() Index {
	return &indexBinSync{indexBin:indexBin{data: []*indexBinItem{}}}
}


//
type indexInterpolationItem struct {
	words    []string
	document string
}

//
func (i *indexInterpolationItem) find(query string, variants []string) int {

	if len(query) == 0 {
		return emptyFind
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
			return mid
		}
	}

	if i.words[low][0] == query[0] {
		for n := low; n < len(i.words); n++ {
			if i.words[n] == query {
				return n
			} else if query[len(query)-1:] == tagAny {
				for j, r := range []rune(query[len(query)-1:]) {
					if r == tagAnyRune {
						return low
					} else if len(i.words[low]) <= n || rune(i.words[low][j]) != r {
						break
					}
				}
			} else if len(variants) > 0 {
				for _, variant := range variants {
					if i.words[n] == variant {
						return low
					}
				}
			}
		}
	}

	if i.words[high][0] == query[0] {
		for n := high; n < len(i.words); n++ {
			if i.words[n] == query {
				return n
			} else if query[len(query)-1:] == tagAny {
				for j, r := range []rune(query[len(query)-1:]) {
					if r == tagAnyRune {
						return low
					} else if len(i.words[low]) <= n || rune(i.words[low][j]) != r {
						break
					}
				}
			} else if len(variants) > 0 {
				for _, variant := range variants {
					if i.words[n] == variant {
						return low
					}
				}
			}
		}
	}

	return emptyFind
}

//
type indexInterpolation struct {
	data []*indexInterpolationItem
}

func (i *indexInterpolation) Find(str string) int {
	words := strings.Split(strings.ToLower(str), ` `)
	for _, word := range words {

		query, variants := i.makeVariants(word)
		for _, d := range i.data {
			if index := d.find(query, variants); index >= 0 {
				return index
			}
		}
	}
	return emptyFind
}

func (i *indexInterpolation) makeVariants(q string) (string, []string) {
	variants := make([]string, 0)

	if q[len(q)-1] == ')' {
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

func (i *indexInterpolation) Add(str ...string) {
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

		n := &indexInterpolationItem{words: words, document: s}
		i.data = append(i.data, n)
	}
}

func (i *indexInterpolation) DocumentAt(index int) (string, bool) {
	if len(i.data) > index && index >= 0 {
		return i.data[index].document, true
	}
	return ``, false
}

//
func NewIndexInterpolation() Index {
	return &indexInterpolation{data: []*indexInterpolationItem{}}
}
