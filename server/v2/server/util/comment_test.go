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
