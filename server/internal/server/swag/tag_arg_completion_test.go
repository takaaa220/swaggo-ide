package swag

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/takaaa220/go-swag-ide/server/internal/server-sdk/protocol"
)

func TestGetTagArgCompletionItems(t *testing.T) {
	t.Parallel()

	type args struct {
		line     string
		position protocol.Position
	}
	tests := map[string]struct {
		args    args
		want    *protocol.CompletionList
		wantErr bool
	}{
		"return_candidates_when_position_is_last": {
			args: args{
				line: `// @Param page `,
				position: protocol.Position{
					Line:      0,
					Character: 15,
				},
			},
			want: &protocol.CompletionList{
				IsIncomplete: false,
				Items: []protocol.CompletionItem{
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "path",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "query",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "header",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "body",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "formData",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "object",
					},
				},
			},
		},
		"return_candidates_when_position_is_middle": {
			args: args{
				line: `// @Param page  int "hello"`,
				position: protocol.Position{
					Line:      0,
					Character: 15,
				},
			},
			want: &protocol.CompletionList{
				IsIncomplete: false,
				Items: []protocol.CompletionItem{
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "path",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "query",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "header",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "body",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "formData",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "object",
					},
				},
			},
		},
		"return_candidates_when_position_is_last_and_args_count_is_one": {
			args: args{
				line: `// @Accept `,
				position: protocol.Position{
					Line:      0,
					Character: 11,
				},
			},
			want: &protocol.CompletionList{
				IsIncomplete: false,
				Items: []protocol.CompletionItem{
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "json",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "application/json",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "xml",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "text/xml",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "plain",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "text/plain",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "html",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "text/html",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "mpfd",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "multipart/form-data",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "x-www-form-urlencoded",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "application/x-www-form-urlencoded",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "json-api",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "application/vnd.api+json",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "json-stream",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "application/x-json-stream",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "octet-stream",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "application/octet-stream",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "png",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "image/png",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "jpeg",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "image/jpeg",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "gif",
					},
					{
						Kind:  protocol.CompletionItemKindKeyword,
						Label: "image/gif",
					},
				},
			},
		},
		"don't_return_candidates_when_the_arg_doesn't_have_candidates": {
			args: args{
				line: `// @Param  query int true `,
				position: protocol.Position{
					Line:      0,
					Character: 10,
				},
			},
			want: &protocol.CompletionList{
				IsIncomplete: false,
				Items:        []protocol.CompletionItem{},
			},
		},
		"don't_return_when_the_count_of_tag_args_is_exceeded": {
			args: args{
				line: `// @Summary hello `,
				position: protocol.Position{
					Line:      0,
					Character: 18,
				},
			},
			want: &protocol.CompletionList{
				IsIncomplete: false,
				Items:        []protocol.CompletionItem{},
			},
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := GetTagArgCompletionItems(tt.args.line, tt.args.position)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTagArgCompletionItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(protocol.CompletionItem{}, "TextEdit"), cmpopts.SortSlices(func(i, j protocol.CompletionItem) bool {
				return i.Label < j.Label
			})); diff != "" {
				t.Errorf("GetTagArgCompletionItems() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
