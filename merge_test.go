package word_index

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testArrayInput1 = [][]ItemID{
		{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		},
		{
			2, 4, 6, 8, 10, 12, 14,
		},
		{
			1, 3, 5, 7, 9, 11, 13,
		},
		{},
	}
	testArrayInput2 = [][]ItemID{
		{
			100, 200, 300, 1000, 1001,
		},
		{
			100, 301, 400, 777, 1001, 9000,
		},
		{
			100, 400, 500, 600, 777, 900, 1001, 1003,
		},
	}
)

func Test_mergeOrderedArray(t *testing.T) {
	t.Parallel()

	result := mergeOrderedArrayOr(testArrayInput1...)
	tt := []ItemID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	require.Equal(t, tt, result)

	result = mergeOrderedArrayOr(testArrayInput2...)
	tt = []ItemID{100, 200, 300, 301, 400, 500, 600, 777, 900, 1000, 1001, 1003, 9000}
	require.Equal(t, tt, result)
}

func Test_mergeOrderedArrayAnd(t *testing.T) {
	t.Parallel()

	result := mergeOrderedArrayAnd(testArrayInput1...)
	tt := make([]ItemID, 0)
	require.Equal(t, tt, result)

	result = mergeOrderedArrayAnd(testArrayInput2...)
	tt = []ItemID{100, 1001}
	require.Equal(t, tt, result)
}
