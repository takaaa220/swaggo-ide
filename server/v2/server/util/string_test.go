package util

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTrimBraces(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name        string
		symbolPairs [][]rune
		input       string
		want        string
	}{
		{
			name: "return function to trim braces",
			symbolPairs: [][]rune{
				{'{', '}'},
				{'[', ']'},
			},
			input: "{hello}",
			want:  "hello",
		},
		{
			name: "return original string if it does not have braces",
			symbolPairs: [][]rune{
				{'{', '}'},
				{'[', ']'},
			},
			input: "hello",
			want:  "hello",
		},
		{
			name: "return original string if it does not have closing brace",
			symbolPairs: [][]rune{
				{'{', '}'},
				{'[', ']'},
			},
			input: "{hello",
			want:  "{hello",
		},
		{
			name: "return original string if it does not have opening brace",
			symbolPairs: [][]rune{
				{'{', '}'},
				{'[', ']'},
			},
			input: "hello}",
			want:  "hello}",
		},
		{
			name: "return original string if input length is less than 2",
			symbolPairs: [][]rune{
				{'{', '}'},
				{'[', ']'},
			},
			input: "{",
			want:  "{",
		},
		{
			name: "return original string if pair length is less than 2",
			symbolPairs: [][]rune{
				{'{'},
				{'}'},
			},
			input: "{hello}",
			want:  "{hello}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TrimBraces(tt.symbolPairs)(tt.input)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
