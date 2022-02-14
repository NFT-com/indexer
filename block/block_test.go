package block_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/NFT-com/indexer/block"
)

func TestBlock_String(t *testing.T) {
	b := block.Block("some_input_text")
	tests := []struct {
		name     string
		input    *block.Block
		expected string
	}{
		{
			name:     "returns an empty string on nil block",
			input:    nil,
			expected: "",
		},
		{
			name:     "returns the correct block string",
			input:    &b,
			expected: "some_input_text",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			output := test.input.String()
			assert.Equal(t, test.expected, output)
		})
	}
}
