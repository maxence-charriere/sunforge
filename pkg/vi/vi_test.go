package vi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnslugify(t *testing.T) {
	tests := []struct {
		slug     string
		expected string
	}{
		{
			slug:     "warrior",
			expected: "Warrior",
		},
		{
			slug:     "demon-hunter",
			expected: "Demon Hunter",
		},
		{
			slug:     "demon_hunter",
			expected: "Demon\nHunter",
		},
	}

	for _, test := range tests {
		t.Run(test.slug, func(t *testing.T) {
			res := Unslugify(test.slug)
			require.Equal(t, test.expected, res)
		})
	}
}
