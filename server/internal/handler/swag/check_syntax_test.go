package swag

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/takaaa220/swaggo-ide/server/internal/handler/protocol"
)

func TestCheckSyntax(t *testing.T) {
	t.Parallel()

	type args struct {
		uri string
		src string
	}
	tests := map[string]struct {
		args args
		want []protocol.Diagnostics
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
			want: []protocol.Diagnostics{},
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
			want: []protocol.Diagnostics{
				{
					Range: protocol.Range{
						Start: protocol.Position{Line: 2, Character: 3},
						End:   protocol.Position{Line: 2, Character: 11},
					},
					Severity: 1,
					Source:   "swag",
					Message:  "Should be `@Summary SUMMARY`.",
				},
				{
					Range: protocol.Range{
						Start: protocol.Position{Line: 3, Character: 3},
						End:   protocol.Position{Line: 3, Character: 15},
					},
					Severity: 1,
					Source:   "swag",
					Message:  "Should be `@Description DESCRIPTION`.",
				},
				{
					Range: protocol.Range{
						Start: protocol.Position{Line: 4, Character: 3},
						End:   protocol.Position{Line: 4, Character: 8},
					},
					Severity: 1,
					Source:   "swag",
					Message:  "Should be `@Tags TAG1,TAG2`.",
				},
				{
					Range: protocol.Range{
						Start: protocol.Position{Line: 5, Character: 11},
						End:   protocol.Position{Line: 5, Character: 18},
					},
					Severity: 1,
					Source:   "swag",
					Message:  "MIME_TYPE should be valid mime type.",
				},
				{
					Range: protocol.Range{
						Start: protocol.Position{Line: 6, Character: 3},
						End:   protocol.Position{Line: 6, Character: 11},
					},
					Severity: 1,
					Source:   "swag",
					Message:  "Should be `@Produce MIME_TYPE`.",
				},
				{
					Range: protocol.Range{
						Start: protocol.Position{Line: 7, Character: 15},
						End:   protocol.Position{Line: 7, Character: 20},
					},
					Severity: 1,
					Source:   "swag",
					Message:  "{DATA_TYPE} should be `string, number, integer, boolean, file or object`.",
				},
				{
					Range: protocol.Range{
						Start: protocol.Position{Line: 8, Character: 26},
						End:   protocol.Position{Line: 8, Character: 29},
					},
					Severity: 1,
					Source:   "swag",
					Message:  "REQUIRED should be `true or false`.",
				},
				{
					Range: protocol.Range{
						Start: protocol.Position{Line: 9, Character: 16},
						End:   protocol.Position{Line: 9, Character: 21},
					},
					Severity: 1,
					Source:   "swag",
					Message:  "{DATA_TYPE} should be `string, number, integer, boolean, file or object`.",
				},
				{
					Range: protocol.Range{
						Start: protocol.Position{Line: 10, Character: 12},
						End:   protocol.Position{Line: 10, Character: 15},
					},
					Severity: 1,
					Source:   "swag",
					Message:  "STATUS_CODE should be integer.",
				},
				{
					Range: protocol.Range{
						Start: protocol.Position{Line: 11, Character: 26},
						End:   protocol.Position{Line: 11, Character: 35},
					},
					Severity: 1,
					Source:   "swag",
					Message:  "[HTTP_METHOD] should be `get, post, put, patch, delete, head, options, trace, or connect`.",
				},
			},
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
