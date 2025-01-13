package swag

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFindSwagComments(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		src  string
		want []SwagCommentsRange
	}{
		"return empty": {
			src: `package main

			// hello world
			func main() {
			}`,
			want: nil,
		},
		"return swag comments": {
			src: `package main

			// hello world
			// @Summary Show an account
			// @Description get string by ID
			func main() {
			}

			//   @Summary hello world
			//   @Description hello world
			// hello world
			func main2() {
			}

			func main3() {
			}
			`,
			want: []SwagCommentsRange{
				{
					Start: 2,
					End:   4,
				},
				{
					Start: 8,
					End:   10,
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := FindSwagComments(tt.src)
			if cmp.Diff(tt.want, got) != "" {
				t.Errorf("FindSwagComments() mismatch (-want +got):\n%s", cmp.Diff(tt.want, got))
			}
		})
	}
}
