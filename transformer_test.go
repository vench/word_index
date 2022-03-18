package word_index

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SequenceToFeature(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		in   string
		out  []Feature
	}{
		{
			name: "empty text",
			in:   "",
			out:  make([]Feature, 0),
		},
		{
			name: "latin text",
			in:   "Delivers patented phone-based verification and two-factor authentication using a time-based, one-time passcode sent over SMS",
			out:  NewFeatures("delivers", "patented", "phone-based", "verification", "and", "two-factor", "authentication", "using", "time-based", "one-time", "passcode", "sent", "over", "sms"),
		},
		{
			name: "ru text",
			in:   "Если цифровая клавиатура не работает на Mac?",
			out:  NewFeatures("если", "цифровая", "клавиатура", "не", "работает", "на", "mac"),
		},
	}

	for i := range tt {
		tc := tt[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := SequenceToFeature(tc.in)
			require.Equal(t, tc.out, r)
		})
	}
}
