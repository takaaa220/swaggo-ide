package swag

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetAttribute(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		line    string
		want    *Attribute
		wantErr bool
	}{
		"return_nil_when_line_is_not_comment": {
			line:    "package main",
			want:    nil,
			wantErr: false,
		},
		"return_nil_when_line_is_not_swag_comment": {
			line:    "// This is a comment",
			want:    nil,
			wantErr: false,
		},
		"return_nil_when_swag_tag_is_invalid": {
			line:    "// @invalid",
			want:    nil,
			wantErr: false,
		},
		"return_attribute_when_swag_tag_is_valid": {
			line: "// @Summary This is a summary",
			want: &Attribute{
				Title:       "@Summary SUMMARY",
				Description: "A short summary of the operation.",
			},
			wantErr: false,
		},
	}
	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := GetAttribute(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAttribute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetAttribute() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
