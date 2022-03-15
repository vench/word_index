package word_index

import "sort"

type Feature string

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

// Insert item save sorting.
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
	index map[Feature]*Items
}

func NewSearch() *Search {
	return &Search{
		index: make(map[Feature]*Items),
	}
}

func (s *Search) Find(feature ...Feature) []ItemID {
	results := make([][]ItemID, len(feature))

	for i := range feature {
		data, exists := s.index[feature[i]]
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
			m, ok := s.index[feature]
			if !ok {
				m = NewItems()
			}
			m.Insert(item)
			s.index[feature] = m
		}
	}
}
