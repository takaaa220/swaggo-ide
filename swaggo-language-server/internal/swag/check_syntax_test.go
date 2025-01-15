package swag

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCheckSyntax(t *testing.T) {
	t.Parallel()

	type args struct {
		uri string
		src string
	}
	tests := map[string]struct {
		args args
		want []SyntaxError
	}{
		"return empty": {
			args: args{
				uri: "test1",
				src: `package main

		// @Summary Show an account
		// @Description get string by ID
		// @Tags accounts
		// @Accept json
		// @Produce json
		// @Header 200 {string} Location "Location of the newly created resource"
		// @Param id path integer true "Account ID"
		// @Success 200 {object} model.Account
		// @Failure 400 {object} httputil.HTTPError
		// @Router /accounts/{id} [get]
				`,
			},
			want: []SyntaxError{},
		},
		"return diagnostics": {
			args: args{
				uri: "test1",
				src: `package main

// @Summary
// @Description
// @Tags
// @Accept unknown
// @Produce
// @Header 200 {sss} Location "Location of the newly created resource"
// @Param id path integer aaa "Account ID"
// @Success 200 {xxx} model.Account
// @Failure aaa {object} httputil.HTTPError
// @Router /accounts/{id} [unknown]
		`,
			},
			want: []SyntaxError{
				{
					Range: Range{
						Start: Position{Line: 2, Character: 3},
						End:   Position{Line: 2, Character: 11},
					},
					Message: "Should be `@Summary SUMMARY`.",
				},
				{
					Range: Range{
						Start: Position{Line: 3, Character: 3},
						End:   Position{Line: 3, Character: 15},
					},
					Message: "Should be `@Description DESCRIPTION`.",
				},
				{
					Range: Range{
						Start: Position{Line: 4, Character: 3},
						End:   Position{Line: 4, Character: 8},
					},
					Message: "Should be `@Tags TAG1,TAG2`.",
				},
				{
					Range: Range{
						Start: Position{Line: 5, Character: 11},
						End:   Position{Line: 5, Character: 18},
					},
					Message: "MIME_TYPE should be valid mime type.",
				},
				{
					Range: Range{
						Start: Position{Line: 6, Character: 3},
						End:   Position{Line: 6, Character: 11},
					},
					Message: "Should be `@Produce MIME_TYPE`.",
				},
				{
					Range: Range{
						Start: Position{Line: 7, Character: 15},
						End:   Position{Line: 7, Character: 20},
					},
					Message: "{DATA_TYPE} should be `string, number, integer, boolean, file or object`.",
				},
				{
					Range: Range{
						Start: Position{Line: 8, Character: 26},
						End:   Position{Line: 8, Character: 29},
					},
					Message: "REQUIRED should be `true or false`.",
				},
				{
					Range: Range{
						Start: Position{Line: 9, Character: 16},
						End:   Position{Line: 9, Character: 21},
					},
					Message: "{DATA_TYPE} should be `string, number, integer, boolean, file or object`.",
				},
				{
					Range: Range{
						Start: Position{Line: 10, Character: 12},
						End:   Position{Line: 10, Character: 15},
					},
					Message: "STATUS_CODE should be integer.",
				},
				{
					Range: Range{
						Start: Position{Line: 11, Character: 26},
						End:   Position{Line: 11, Character: 35},
					},
					Message: "[HTTP_METHOD] should be `get, post, put, patch, delete, head, options, trace, or connect`.",
				},
			},
		},
		"return diagnostics_when_router_doesn't_exist": {
			args: args{
				uri: "test1",
				src: `package main

// @Summary hello world
func hello() {
}
`,
			},
			want: []SyntaxError{
				{
					Range: Range{
						Start: Position{Line: 2, Character: 0},
						End:   Position{Line: 2, Character: 0},
					},
					Message: "@Router is required.",
				},
			},
		},
		"don't_return diagnostics_when_swaggo_comments_don't_exist": {
			args: args{
				uri: "test1",
				src: `package main

// hello world
func hello() {
}
`,
			},
			want: []SyntaxError{},
		},
	}
	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := CheckSyntax(tt.args.uri, tt.args.src)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
