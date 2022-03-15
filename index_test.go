package word_index

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearch_Find(t *testing.T) {
	t.Parallel()

	s := NewSearch()
	s.Add(&Item{
		ID:      3,
		Feature: NewFeatures("abc", "xyz"),
	})
	s.Add(&Item{
		ID:      2,
		Feature: NewFeatures("xyz"),
	})
	s.Add(&Item{
		ID:      1,
		Feature: NewFeatures("abc"),
	})
	s.Add(&Item{
		ID:      4,
		Feature: NewFeatures("foo", "bar"),
	})
	s.Add(&Item{
		ID:      5,
		Feature: NewFeatures("foo", "gaz"),
	})

	require.Equal(t, 5, len(s.index))

	result := s.Find("foo")
	require.Equal(t, 2, len(result))
	require.Equal(t, ItemID(4), result[0])
	require.Equal(t, ItemID(5), result[1])
}

func TestItems_Insert(t *testing.T) {
	t.Parallel()

	var i Items
	i.Insert(&Item{
		ID:      1,
		Feature: nil,
	})
	i.Insert(&Item{
		ID:      5,
		Feature: nil,
	})
	i.Insert(&Item{
		ID:      4,
		Feature: nil,
	})

	i.Insert(&Item{
		ID:      3,
		Feature: nil,
	})
	i.Insert(&Item{
		ID:      2,
		Feature: nil,
	})

	require.True(t, len(i.items) > 0)

	for j := 1; j < len(i.items); j++ {
		require.Greater(t, i.items[j].ID, i.items[j-1].ID)
	}
}
