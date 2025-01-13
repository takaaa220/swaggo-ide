package swag

import "testing"

func TestIsCommentLine(t *testing.T) {
	t.Parallel()

	type args struct {
		line string
	}
	tests := []struct {
		args args
		want bool
	}{
		{
			args: args{
				line: "// @Summary Show an account",
			},
			want: true,
		},
		{
			args: args{
				line: "package main",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.args.line, func(t *testing.T) {
			t.Parallel()
			if got := isCommentLine(tt.args.line); got != tt.want {
				t.Errorf("IsCommentLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
