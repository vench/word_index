package word_index

import (
	"strings"
	"unicode"
)

func NewDocumentSearch(document ...string) *Search {
	items := make([]*Item, len(document))

	for i := range document {
		items[i] = &Item{
			ID:      ItemID(i + 1),
			Feature: SequenceToFeature(document[i]),
		}
	}

	s := NewSearch()
	s.Add(items...)

	return s
}

func SequenceToFeature(sequence string) []Feature {
	s := strings.FieldsFunc(sequence, func(r rune) bool {
		if unicode.IsSpace(r) {
			return true
		}

		switch r {
		case ',', '.', '!', '?', ':', ';':
			return true
		}

		return false
	})

	transform := func(s string) string {
		s = strings.ToLower(s)
		s = strings.TrimSpace(s)
		return s
	}

	validate := func(s string) bool {
		return len(s) > 1
	}

	f := make([]Feature, 0, len(s))
	for i := range s {
		w := transform(s[i])
		if !validate(w) {
			continue
		}

		f = append(f, Feature(w))
	}

	return f
}
