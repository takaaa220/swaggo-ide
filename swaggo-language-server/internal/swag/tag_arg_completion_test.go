package swag

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestGetTagArgCompletionItems(t *testing.T) {
	t.Parallel()

	type args struct {
		line     string
		position Position
	}
	tests := map[string]struct {
		args    args
		want    []CompletionCandidate
		wantErr bool
	}{
		"return_candidates_when_position_is_last": {
			args: args{
				line: `// @Param page `,
				position: Position{
					Line:      0,
					Character: 15,
				},
			},
			want: []CompletionCandidate{
				{
					NewText: "path",
					Label:   "path",
				},
				{
					NewText: "query",
					Label:   "query",
				},
				{
					NewText: "header",
					Label:   "header",
				},
				{
					NewText: "body",
					Label:   "body",
				},
				{
					NewText: "formData",
					Label:   "formData",
				},
				{
					NewText: "object",
					Label:   "object",
				},
			},
		},
		"return_candidates_when_position_is_middle": {
			args: args{
				line: `// @Param page  int "hello"`,
				position: Position{
					Line:      0,
					Character: 15,
				},
			},
			want: []CompletionCandidate{
				{
					NewText: "path",
					Label:   "path",
				},
				{
					NewText: "query",
					Label:   "query",
				},
				{
					NewText: "header",
					Label:   "header",
				},
				{
					NewText: "body",
					Label:   "body",
				},
				{
					NewText: "formData",
					Label:   "formData",
				},
				{
					NewText: "object",
					Label:   "object",
				},
			},
		},
		"return_candidates_when_position_is_last_and_args_count_is_one": {
			args: args{
				line: `// @Accept `,
				position: Position{
					Line:      0,
					Character: 11,
				},
			},
			want: []CompletionCandidate{
				{
					NewText: "json",
					Label:   "json",
				},
				{
					NewText: "application/json",
					Label:   "application/json",
				},
				{
					NewText: "xml",
					Label:   "xml",
				},
				{
					NewText: "text/xml",
					Label:   "text/xml",
				},
				{
					NewText: "plain",
					Label:   "plain",
				},
				{
					NewText: "text/plain",
					Label:   "text/plain",
				},
				{
					NewText: "html",
					Label:   "html",
				},
				{
					NewText: "text/html",
					Label:   "text/html",
				},
				{
					NewText: "mpfd",
					Label:   "mpfd",
				},
				{
					NewText: "multipart/form-data",
					Label:   "multipart/form-data",
				},
				{
					NewText: "x-www-form-urlencoded",
					Label:   "x-www-form-urlencoded",
				},
				{
					NewText: "application/x-www-form-urlencoded",
					Label:   "application/x-www-form-urlencoded",
				},
				{
					NewText: "json-api",
					Label:   "json-api",
				},
				{
					NewText: "application/vnd.api+json",
					Label:   "application/vnd.api+json",
				},
				{
					NewText: "json-stream",
					Label:   "json-stream",
				},
				{
					NewText: "application/x-json-stream",
					Label:   "application/x-json-stream",
				},
				{
					NewText: "octet-stream",
					Label:   "octet-stream",
				},
				{
					NewText: "application/octet-stream",
					Label:   "application/octet-stream",
				},
				{
					NewText: "png",
					Label:   "png",
				},
				{
					NewText: "image/png",
					Label:   "image/png",
				},
				{
					NewText: "jpeg",
					Label:   "jpeg",
				},
				{
					NewText: "image/jpeg",
					Label:   "image/jpeg",
				},
				{
					NewText: "gif",
					Label:   "gif",
				},
				{
					NewText: "image/gif",
					Label:   "image/gif",
				},
			},
		},
		"don't_return_candidates_when_the_arg_doesn't_have_candidates": {
			args: args{
				line: `// @Param  query int true `,
				position: Position{
					Line:      0,
					Character: 10,
				},
			},
			want: []CompletionCandidate{},
		},
		"don't_return_when_the_count_of_tag_args_is_exceeded": {
			args: args{
				line: `// @Summary hello `,
				position: Position{
					Line:      0,
					Character: 18,
				},
			},
			want: []CompletionCandidate{},
		},
		"return_data_type_candidates_for_success_tag": {
			args: args{
				line: `// @Success 200 `,
				position: Position{
					Line:      0,
					Character: 16,
				},
			},
			want: []CompletionCandidate{
				{
					NewText: "{string}",
					Label:   "{string}",
				},
				{
					NewText: "{number}",
					Label:   "{number}",
				},
				{
					NewText: "{integer}",
					Label:   "{integer}",
				},
				{
					NewText: "{boolean}",
					Label:   "{boolean}",
				},
				{
					NewText: "{file}",
					Label:   "{file}",
				},
				{
					NewText: "{object}",
					Label:   "{object}",
				},
				{
					NewText: "{array}",
					Label:   "{array}",
				},
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
			if diff := cmp.Diff(tt.want, got, cmpopts.SortSlices(func(i, j CompletionCandidate) bool {
				return i.Label < j.Label
			})); diff != "" {
				t.Errorf("GetTagArgCompletionItems() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
