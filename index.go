package word_index

import (
	"strings"
	"fmt"
	"regexp"
	"sort"
)

type Index interface {
	Find(string) bool
	Add(...string)
}

type index struct {
	data []string
}

//
func (i*index) Find(str string) bool  {
	words := strings.Split(strings.ToLower(str), ` `)
	for _,word := range words {
		expr := fmt.Sprintf(`(^|\s)%s(\s|\.|\,|$)`, strings.TrimSpace(word))
		r,_ := regexp.Compile(expr)
		for _, dataStr := range i.data {
			if len(r.FindStringSubmatchIndex(dataStr)) > 0 {
				return true
			}
		}
	}
	return false
}

//
func (i*index) Add(str ... string ) {
	for i :=0; i < len(str); i ++ {
		str[i] = strings.ToLower(str[i])
	}
	i.data = append(i.data, str...)
}

//
func NewBaseIndex()Index {
	return &index{data:[]string{}}
}


//
type indexCharsItem struct {
	words []string
}

//
func (i *indexCharsItem) Find(query string) bool {

	variants := make([]string, 0)
	if query[len(query)-1] == ')' {
		base := make([]rune, 0)
		start := false
		variant := make([]rune, 0)
		for _,r := range []rune(query){
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
			} else if start  {
				variant = append(variant, r)
			} else {
				base = append(base, r)
			}
		}
	}

	low := 0
	high := len(i.words) - 1

	for low <= high {
		median := (low + high) / 2
		if i.words[median][0] < query[0] {
			low = median + 1
		}else{
			high = median - 1
		}
	}

	for low < len(i.words) && i.words[low][0] == query[0] {
		if i.words[low] == query {
			return true
		} else if len(variants) > 0 {
			for _,variant := range variants {
				if i.words[low] == variant {
					return true
				}
			}
		}
		low ++
	}
	return  false
}

//
type indexChars struct {
	data []*indexCharsItem
}

//
func (i*indexChars) Add(str ... string ) {
	for _, s := range str {
		words := strings.Split(strings.ToLower(s), ` `)
		sort.Slice(words, func(i, j int) bool {
			if words[i] < words[j] {
				return true
			}
			return false
		})
		n := &indexCharsItem{words:words}
		i.data = append(i.data, n)
	}
}

//
func (i*indexChars) Find(str string) bool  {
	words := strings.Split(strings.ToLower(str), ` `)
	for _, word := range words {
		for _, dataStr := range i.data {
			if dataStr.Find(word) {
				return  true
			}
		}
	}
	return  false
}

//
func NewIndexChars()Index {
	return &indexChars{data:[]*indexCharsItem{}}
}