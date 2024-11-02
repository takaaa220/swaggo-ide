package swag

import (
	"fmt"
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
		args                  splitterArgs
		wantFirstSplitElement splitElement
		wantRest              []splitElement
	}{
		{
			args: splitterArgs{
				str:           `@Param 			id   path      int  true  "Account ID"`,
				maxSplitCount: -1,
			},
			wantFirstSplitElement: splitElement{
				Text:  "@Param",
				Start: 0,
				End:   6,
			},
			wantRest: []splitElement{
				{
					Text:  "id",
					Start: 10,
					End:   12,
				},
				{
					Text:  "path",
					Start: 15,
					End:   19,
				},
				{
					Text:  "int",
					Start: 25,
					End:   28,
				},
				{
					Text:  "true",
					Start: 30,
					End:   34,
				},
				{
					Text:  "\"Account ID\"",
					Start: 36,
					End:   48,
				},
			},
		},
		{
			args: splitterArgs{
				str:           ` // @Param 			id   path      int  true  "Account ID"`,
				maxSplitCount: -1,
			},
			wantFirstSplitElement: splitElement{
				Text:  "@Param",
				Start: 4,
				End:   10,
			},
			wantRest: []splitElement{
				{
					Text:  "id",
					Start: 14,
					End:   16,
				},
				{
					Text:  "path",
					Start: 19,
					End:   23,
				},
				{
					Text:  "int",
					Start: 29,
					End:   32,
				},
				{
					Text:  "true",
					Start: 34,
					End:   38,
				},
				{
					Text:  "\"Account ID\"",
					Start: 40,
					End:   52,
				},
			},
		},
		{
			args: splitterArgs{
				str:           `@Success      200  {object} model.Account `,
				maxSplitCount: -1,
			},
			wantFirstSplitElement: splitElement{
				Text:  "@Success",
				Start: 0,
				End:   8,
			},
			wantRest: []splitElement{
				{
					Text:  "200",
					Start: 14,
					End:   17,
				},
				{
					Text:  "{object}",
					Start: 19,
					End:   27,
				},
				{
					Text:  "model.Account",
					Start: 28,
					End:   41,
				},
			},
		},
		{
			args: splitterArgs{
				str:           `@Router       /accounts/{id} [get]`,
				maxSplitCount: -1,
			},
			wantFirstSplitElement: splitElement{
				Text:  "@Router",
				Start: 0,
				End:   7,
			},
			wantRest: []splitElement{
				{
					Text:  "/accounts/{id}",
					Start: 14,
					End:   28,
				},
				{
					Text:  "[get]",
					Start: 29,
					End:   34,
				},
			},
		},
		{
			args: splitterArgs{
				str:           `@Accept`,
				maxSplitCount: 1,
			},
			wantFirstSplitElement: splitElement{
				Text:  "@Accept",
				Start: 0,
				End:   7,
			},
			wantRest: []splitElement{},
		},
		{
			args: splitterArgs{
				str:           `@Summary  hello world test`,
				maxSplitCount: 1,
			},
			wantFirstSplitElement: splitElement{
				Text:  "@Summary",
				Start: 0,
				End:   8,
			},
			wantRest: []splitElement{
				{
					Text:  "hello world test",
					Start: 10,
					End:   26,
				},
			},
		},
		{
			args: splitterArgs{
				str:           `@Summary  hello`,
				maxSplitCount: 3,
			},
			wantFirstSplitElement: splitElement{
				Text:  "@Summary",
				Start: 0,
				End:   8,
			},
			wantRest: []splitElement{
				{
					Text:  "hello",
					Start: 10,
					End:   15,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%s_%d", tt.args.str, tt.args.maxSplitCount), func(t *testing.T) {
			t.Parallel()

			firstSplitElement, splitRest := split(tt.args.str)
			rest := []splitElement{}
			for splitElement := range splitRest(tt.args.maxSplitCount) {
				rest = append(rest, splitElement)
			}

			if diff := cmp.Diff(tt.wantFirstSplitElement, firstSplitElement); diff != "" {
				t.Errorf("first split element mismatch (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tt.wantRest, rest); diff != "" {
				t.Errorf("rest mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// P sss  bbb

// -1
// 0
