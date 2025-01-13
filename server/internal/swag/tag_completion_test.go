package swag

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetCompletionItems(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		want    []CompletionCandidate
		wantErr bool
	}{
		"return completion items": {
			want: []CompletionCandidate{
				{
					Label:   "@Accept",
					NewText: "@Accept MIME_TYPE",
				},
				{
					Label:   "@Description",
					NewText: "@Description DESCRIPTION",
				},
				{
					Label:   "@Failure",
					NewText: "@Failure STATUS_CODE {DATA_TYPE} GO_TYPE",
				},
				{
					Label:   "@Header",
					NewText: "@Header STATUS_CODE {DATA_TYPE} HEADER_NAME COMMENT",
				},
				{
					Label:   "@ID",
					NewText: "@ID ID",
				},
				{
					Label:   "@Param",
					NewText: "@Param PARAM_NAME PARAM_TYPE GO_TYPE REQUIRED \"DESCRIPTION\"",
				},
				{
					Label:   "@Produce",
					NewText: "@Produce MIME_TYPE",
				},
				{
					Label:   "@Router",
					NewText: "@Router PATH [HTTP_METHOD]",
				},
				{
					Label:   "@Success",
					NewText: "@Success STATUS_CODE {DATA_TYPE} GO_TYPE",
				},
				{
					Label:   "@Summary",
					NewText: "@Summary SUMMARY",
				},
				{
					Label:   "@Tags",
					NewText: "@Tags TAG1,TAG2",
				},
			},
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := GetTagCompletionItems("// @P")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCompletionItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetCompletionItems() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
