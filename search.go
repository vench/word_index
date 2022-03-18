package word_index

import (
	"sort"
	"strings"
)

const anyForm = "*"

type FeatureID uint64

type Feature string

func (f Feature) AnyForm() bool {
	return f[len(f)-1:] == anyForm
}

func (f Feature) Form() Feature {
	return f[:len(f)-1]
}

func (f Feature) String() string {
	return string(f)
}

var (
	emptyItemID = make([]ItemID, 0)
)

func NewFeatures(feature ...Feature) []Feature {
	return feature
}

type ItemID int64

type Item struct {
	ID      ItemID
	Feature []Feature
}

func NewItem(id ItemID, feature ...Feature) *Item {
	return &Item{
		ID:      id,
		Feature: feature,
	}
}

type Items struct {
	items []*Item
}

func NewItems() *Items {
	return &Items{
		items: make([]*Item, 0),
	}
}

// Insert item save sorting by ID.
func (i *Items) Insert(item *Item) {
	index := sort.Search(len(i.items), func(j int) bool {
		return i.items[j].ID >= item.ID
	})

	if len(i.items) == index {
		i.items = append(i.items, item)
		return
	}

	i.items = append(i.items[:index+1], i.items[index:]...)
	i.items[index] = item
}

type Search struct {
	index       map[FeatureID]*Items
	featureDict map[Feature]FeatureID
	features    []Feature
}

func NewSearch() *Search {
	return &Search{
		index:       make(map[FeatureID]*Items),
		featureDict: make(map[Feature]FeatureID),
	}
}

func (s *Search) transformFeature(feature ...Feature) []FeatureID {
	features := make([]FeatureID, 0, len(feature))

	for i := range feature {
		if id, exists := s.featureDict[feature[i]]; exists {
			features = append(features, id)
			continue
		}

		if !feature[i].AnyForm() {
			continue
		}

		form := feature[i].Form()
		if len(form) == 0 {
			continue
		}

		index := sort.Search(len(s.features), func(j int) bool {
			return s.features[j] >= form
		})

		for ; index < len(s.features); index++ {
			f := s.features[index]
			if !strings.Contains(f.String(), form.String()) {
				break
			}

			id, exists := s.featureDict[f]
			if !exists {
				continue
			}
			features = append(features, id)
		}
	}

	return features
}

func (s *Search) Find(feature ...Feature) []ItemID {
	featureTransform := s.transformFeature(feature...)
	results := make([][]ItemID, len(featureTransform))

	for i := range featureTransform {
		data, exists := s.index[featureTransform[i]]
		if !exists {
			results[i] = emptyItemID
		} else {
			results[i] = make([]ItemID, len(data.items))
			for j := range data.items {
				results[i][j] = data.items[j].ID
			}
		}
	}

	return mergeOrderedArrayOr(results...)
}

func (s *Search) Add(items ...*Item) {
	for i := range items {
		item := items[i]
		for j := range item.Feature {
			feature := item.Feature[j]

			id, ok := s.featureDict[feature]
			if !ok {
				id = FeatureID(len(s.featureDict) + 1)
				s.featureDict[feature] = id
			}

			m, ok := s.index[id]
			if !ok {
				m = NewItems()
			}
			m.Insert(item)
			s.index[id] = m
		}
	}

	features := NewFeatures()
	for f := range s.featureDict {
		features = append(features, f)
	}

	sort.Slice(features, func(i, j int) bool {
		return features[i] < features[j]
	})

	s.features = features
}
