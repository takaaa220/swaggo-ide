package util

import (
	"testing"

	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/protocol"
)

func TestIsInComment(t *testing.T) {
	t.Parallel()

	type args struct {
		src string
		pos protocol.Position
	}
	tests := map[string]struct {
		args args
		want bool
	}{
		"line comment1": {
			args: args{
				src: "hello world\n// comment\nhello world",
				pos: protocol.Position{
					Line:      1,
					Character: 3,
				},
			},
			want: true,
		},
		"line comment2": {
			args: args{
				src: "hello world\n 	// comment\nhello world",
				pos: protocol.Position{
					Line:      1,
					Character: 8,
				},
			},
			want: true,
		},
		"not comment1": {
			args: args{
				src: "hello world\n 	// comment\nhello world",
				pos: protocol.Position{
					Line:      2,
					Character: 8,
				},
			},
			want: false,
		},
	}
	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := IsInComment(tt.args.src, tt.args.pos)
			if got != tt.want {
				t.Errorf("IsInComment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrimPrefixForComment(t *testing.T) {
	type args struct {
		line string
	}
	tests := map[string]struct {
		args             args
		want             string
		wantTrimmedCount int
	}{
		"return_text_trimmed_the_prefix_for_comment_1": {
			args: args{
				line: " 	// comment comment",
			},
			want:             "comment comment",
			wantTrimmedCount: 5,
		},
		"return_text_trimmed_the_prefix_for_comment_2": {
			args: args{
				line: "// comment comment",
			},
			want:             "comment comment",
			wantTrimmedCount: 3,
		},
		"return original text if it is not a comment": {
			args: args{
				line: "this is not a comment",
			},
			want:             "this is not a comment",
			wantTrimmedCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.args.line, func(t *testing.T) {
			got, trimmedCount := TrimPrefixForComment(tt.args.line)
			if got != tt.want {
				t.Errorf("TrimPrefixForComment().trimmed = %v, want %v", got, tt.want)
			}
			if trimmedCount != tt.wantTrimmedCount {
				t.Errorf("TrimPrefixForComment().trimmedCount = %v, want %v", trimmedCount, tt.wantTrimmedCount)
			}
		})
	}
}
