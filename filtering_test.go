package word_index

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func newTestSearch(t *testing.T) *Search {
	t.Helper()

	s := NewSearch()
	s.Add(
		NewItem(1001, "march", "student", "assassinated", "leading", "politician", "main", "architect", "genocide"),
		NewItem(2001, "politician", "architect", "genocide", "million", "genocide", "including", "most"),
		NewItem(3001, "family", "joined", "about", "berlin", "clandestine", "assassination", "trial"),
		NewItem(4001, "seeking", "campaign", "revenge", "trial", "held", "strategy", "defense"),
		NewItem(5001, "was", "june", "genocide"),
	)
	return s
}

func Test_Filtering(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name  string
		input Filtering
		out   []ItemID
	}{
		{
			name:  "filter in genocide",
			input: NewFilterIn("genocide"),
			out:   []ItemID{1001, 2001, 5001},
		},
		{
			name:  "filter in march",
			input: NewFilterIn("march"),
			out:   []ItemID{1001},
		},
		{
			name:  "filter in march,was",
			input: NewFilterIn("march", "was"),
			out:   []ItemID{1001, 5001},
		},
		{
			name: "filter march and leading",
			input: NewAndOperator(
				NewFilterIn("march"),
				NewFilterIn("leading"),
			),
			out: []ItemID{1001},
		},
		{
			name: "filter march or trial",
			input: NewOrOperator(
				NewFilterIn("march"),
				NewFilterIn("trial"),
			),
			out: []ItemID{1001, 3001, 4001},
		},
		{
			name: "filter march or (trial and held)",
			input: NewOrOperator(
				NewFilterIn("march"),
				NewAndOperator(
					NewFilterIn("trial"),
					NewFilterIn("held"),
				),
			),
			out: []ItemID{1001, 4001},
		},
		{
			name: "filter trial and (assassination or held)",
			input: NewAndOperator(
				NewFilterIn("trial"),
				NewOrOperator(
					NewFilterIn("assassination"),
					NewFilterIn("held"),
				),
			),
			out: []ItemID{3001, 4001},
		},
	}

	s := newTestSearch(t)

	for i := range tt {
		tc := tt[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := tc.input.Filter(s)
			require.Equal(t, tc.out, result)
		})
	}
}
