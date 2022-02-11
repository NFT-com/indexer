package block_test

import (
	"testing"

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
			name:     "should return the correct block string",
			input:    &b,
			expected: "some_input_text",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := test.input.String()
			if output != test.expected {
				t.Errorf("test: %s failed due to %s is not the expected %s", test.name, output, test.expected)
				return
			}
		})
	}
}
