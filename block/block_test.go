package block_test

import (
	"testing"

	"github.com/NFT-com/indexer/block"
)

func TestBlock_String(t *testing.T) {
	b := block.Block("some_input_text")
	tts := []struct {
		name     string
		input    *block.Block
		expected string
	}{
		{
			name:     "should return a empty string on nil block",
			input:    nil,
			expected: "",
		},
		{
			name:     "should return the correct block string",
			input:    &b,
			expected: "some_input_text",
		},
	}
	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.input.String()
			if output != tt.expected {
				t.Errorf("test: %s failed due to %s is not the expected %s", tt.name, output, tt.expected)
				return
			}
		})
	}
}
