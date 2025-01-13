package swag

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/takaaa220/swaggo-ide/server/internal/handler/protocol"
)

func TestGetAttribute(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		line     string
		position protocol.Position
		want     *Attribute
		wantErr  bool
	}{
		"return_nil_when_line_is_not_comment": {
			line:     "package main",
			position: protocol.Position{},
			want:     nil,
			wantErr:  false,
		},
		"return_nil_when_line_is_not_swag_comment": {
			line:     "// This is a comment",
			position: protocol.Position{},
			want:     nil,
			wantErr:  false,
		},
		"return_nil_when_swag_tag_is_invalid": {
			line:     "// @invalid",
			position: protocol.Position{},
			want:     nil,
			wantErr:  false,
		},
		"return_attribute_when_swag_tag_is_valid": {
			line:     "// @Summary This is a summary",
			position: protocol.Position{},
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

			got, err := GetAttribute(tt.line, tt.position)
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
