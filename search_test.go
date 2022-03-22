package word_index

import (
	"os"
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

	require.Equal(t, NewFeatures("abc", "bar", "foo", "gaz", "xyz"), s.features)
	require.Equal(t, 5, len(s.index))

	result := s.Find("foo")
	require.Equal(t, 2, len(result))
	require.Equal(t, ItemID(4), result[0])
	require.Equal(t, ItemID(5), result[1])
}

func TestItems_Insert(t *testing.T) {
	t.Parallel()

	var i items
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

	require.True(t, len(i.Items) > 0)

	for j := 1; j < len(i.Items); j++ {
		require.Greater(t, i.Items[j].ID, i.Items[j-1].ID)
	}
}

func TestSearch_transformFeature(t *testing.T) {
	t.Parallel()

	s := NewSearch()
	s.Add(&Item{
		ID:      5,
		Feature: NewFeatures("foo", "gaz"),
	})

	ids := s.transformFeature("foo")
	require.Equal(t, []FeatureID{1}, ids)

	ids = s.transformFeature("foo*")
	require.Equal(t, []FeatureID{1}, ids)

	ids = s.transformFeature("*")
	require.Equal(t, []FeatureID{}, ids)

	ids = s.transformFeature("abc*")
	require.Equal(t, []FeatureID{}, ids)
}

func TestSearch_SaveLoad(t *testing.T) {
	t.Parallel()

	f, err := os.CreateTemp("/tmp", "search-test-*")
	require.NoError(t, err)
	defer f.Close()

	s := NewSearch()
	s.Add(&Item{
		ID:      5,
		Feature: NewFeatures("foo", "gaz"),
	},
		&Item{
			ID:      7,
			Feature: NewFeatures("test", "best", "bar", "foo"),
		})

	err = s.Save(f)
	require.NoError(t, err)

	fr, err := os.Open(f.Name())
	require.NoError(t, err)
	defer fr.Close()

	sLoad := NewSearch()
	err = sLoad.Load(fr)
	require.NoError(t, err)

	require.Equal(t, s, sLoad)
}
