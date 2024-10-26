package internal

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_splitter_split(t *testing.T) {
	t.Parallel()

	type splitterArgs struct {
		str           string
		maxSplitCount int
	}

	tests := []struct {
		args splitterArgs
		want []string
	}{
		{
			args: splitterArgs{
				str:           `@Param 			id   path      int  true  "Account ID"`,
				maxSplitCount: -1,
			},
			want: []string{"@Param", "id", "path", "int", "true", "Account ID"},
		},
		{
			args: splitterArgs{
				str:           `@Success      200  {object} model.Account `,
				maxSplitCount: -1,
			},
			want: []string{"@Success", "200", "{object}", "model.Account"},
		},
		{
			args: splitterArgs{
				str:           `@Router       /accounts/{id} [get]`,
				maxSplitCount: -1,
			},
			want: []string{"@Router", "/accounts/{id}", "[get]"},
		},
		{
			args: splitterArgs{
				str:           `@Summary hello world test`,
				maxSplitCount: 2,
			},
			want: []string{"@Summary", "hello world test"},
		},
		{
			args: splitterArgs{
				str:           `@Summary  hello`,
				maxSplitCount: 3,
			},
			want: []string{"@Summary", "hello"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.args.str, func(t *testing.T) {
			t.Parallel()

			s := newSplitter(tt.args.str, tt.args.maxSplitCount)

			got := s.split()
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Error("mismatch (-want, +got)", diff)
			}
		})
	}
}
