package swag

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_tokenizer_tokenize(t *testing.T) {
	t.Parallel()

	type tokenizerArgs struct {
		str           string
		maxTokenCount int
	}

	tests := []struct {
		args           tokenizerArgs
		wantFirstToken token
		wantRest       []token
	}{
		{
			args: tokenizerArgs{
				str:           `@Param 			id   path      int  true  "Account ID"`,
				maxTokenCount: -1,
			},
			wantFirstToken: token{
				Text:  "@Param",
				Start: 0,
				End:   6,
			},
			wantRest: []token{
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
			args: tokenizerArgs{
				str:           ` // @Param 			id   path      int  true  "Account ID"`,
				maxTokenCount: -1,
			},
			wantFirstToken: token{
				Text:  "@Param",
				Start: 4,
				End:   10,
			},
			wantRest: []token{
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
			args: tokenizerArgs{
				str:           `@Success      200  {object} model.Account `,
				maxTokenCount: -1,
			},
			wantFirstToken: token{
				Text:  "@Success",
				Start: 0,
				End:   8,
			},
			wantRest: []token{
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
			args: tokenizerArgs{
				str:           `@Router       /accounts/{id} [get]`,
				maxTokenCount: -1,
			},
			wantFirstToken: token{
				Text:  "@Router",
				Start: 0,
				End:   7,
			},
			wantRest: []token{
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
			args: tokenizerArgs{
				str:           `@Accept`,
				maxTokenCount: 1,
			},
			wantFirstToken: token{
				Text:  "@Accept",
				Start: 0,
				End:   7,
			},
			wantRest: []token{},
		},
		{
			args: tokenizerArgs{
				str:           `@Summary  hello world test`,
				maxTokenCount: 1,
			},
			wantFirstToken: token{
				Text:  "@Summary",
				Start: 0,
				End:   8,
			},
			wantRest: []token{
				{
					Text:  "hello world test",
					Start: 10,
					End:   26,
				},
			},
		},
		{
			args: tokenizerArgs{
				str:           `@Summary  hello`,
				maxTokenCount: 3,
			},
			wantFirstToken: token{
				Text:  "@Summary",
				Start: 0,
				End:   8,
			},
			wantRest: []token{
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
		t.Run(fmt.Sprintf("%s_%d", tt.args.str, tt.args.maxTokenCount), func(t *testing.T) {
			t.Parallel()

			firstToken, tokenizeRest := tokenize(tt.args.str)
			rest := []token{}
			for token := range tokenizeRest(tt.args.maxTokenCount) {
				rest = append(rest, token)
			}

			if diff := cmp.Diff(tt.wantFirstToken, firstToken); diff != "" {
				t.Errorf("first tokenize element mismatch (-want +got):\n%s", diff)
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
