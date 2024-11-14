package swag

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/protocol"
)

func TestGetCompletionItems(t *testing.T) {
	t.Parallel()

	type args struct {
		position protocol.Position
	}
	tests := map[string]struct {
		args    args
		want    *protocol.CompletionList
		wantErr bool
	}{
		"return completion items": {
			args: args{
				position: protocol.Position{
					Line:      0,
					Character: 1,
				},
			},
			want: &protocol.CompletionList{
				IsIncomplete: false,
				Items: []protocol.CompletionItem{
					{
						Label: "@Accept",
						Kind:  protocol.CompletionItemKindKeyword,
						TextEdit: &protocol.TextEdit{
							NewText: "@Accept MIME_TYPE",
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      0,
									Character: 0,
								},
								End: protocol.Position{
									Line:      0,
									Character: 1,
								},
							},
						},
					},
					{
						Label: "@Description",
						Kind:  protocol.CompletionItemKindKeyword,
						TextEdit: &protocol.TextEdit{
							NewText: "@Description DESCRIPTION",
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      0,
									Character: 0,
								},
								End: protocol.Position{
									Line:      0,
									Character: 1,
								},
							},
						},
					},
					{
						Label: "@Failure",
						Kind:  protocol.CompletionItemKindKeyword,
						TextEdit: &protocol.TextEdit{
							NewText: "@Failure STATUS_CODE {DATA_TYPE} GO_TYPE",
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      0,
									Character: 0,
								},
								End: protocol.Position{
									Line:      0,
									Character: 1,
								},
							},
						},
					},
					{
						Label: "@Header",
						Kind:  protocol.CompletionItemKindKeyword,
						TextEdit: &protocol.TextEdit{
							NewText: "@Header STATUS_CODE {DATA_TYPE} HEADER_NAME COMMENT",
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      0,
									Character: 0,
								},
								End: protocol.Position{
									Line:      0,
									Character: 1,
								},
							},
						},
					},
					{
						Label: "@ID",
						Kind:  protocol.CompletionItemKindKeyword,
						TextEdit: &protocol.TextEdit{
							NewText: "@ID ID",
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      0,
									Character: 0,
								},
								End: protocol.Position{
									Line:      0,
									Character: 1,
								},
							},
						},
					},
					{
						Label: "@Param",
						Kind:  protocol.CompletionItemKindKeyword,
						TextEdit: &protocol.TextEdit{
							NewText: "@Param PARAM_NAME PARAM_TYPE GO_TYPE REQUIRED \"DESCRIPTION\"",
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      0,
									Character: 0,
								},
								End: protocol.Position{
									Line:      0,
									Character: 1,
								},
							},
						},
					},
					{
						Label: "@Produce",
						Kind:  protocol.CompletionItemKindKeyword,
						TextEdit: &protocol.TextEdit{
							NewText: "@Produce MIME_TYPE",
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      0,
									Character: 0,
								},
								End: protocol.Position{
									Line:      0,
									Character: 1,
								},
							},
						},
					},
					{
						Label: "@Router",
						Kind:  protocol.CompletionItemKindKeyword,
						TextEdit: &protocol.TextEdit{
							NewText: "@Router PATH [HTTP_METHOD]",
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      0,
									Character: 0,
								},
								End: protocol.Position{
									Line:      0,
									Character: 1,
								},
							},
						},
					},
					{
						Label: "@Success",
						Kind:  protocol.CompletionItemKindKeyword,
						TextEdit: &protocol.TextEdit{
							NewText: "@Success STATUS_CODE {DATA_TYPE} GO_TYPE",
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      0,
									Character: 0,
								},
								End: protocol.Position{
									Line:      0,
									Character: 1,
								},
							},
						},
					},
					{
						Label: "@Summary",
						Kind:  protocol.CompletionItemKindKeyword,
						TextEdit: &protocol.TextEdit{
							NewText: "@Summary SUMMARY",
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      0,
									Character: 0,
								},
								End: protocol.Position{
									Line:      0,
									Character: 1,
								},
							},
						},
					},
					{
						Label: "@Tags",
						Kind:  protocol.CompletionItemKindKeyword,
						TextEdit: &protocol.TextEdit{
							NewText: "@Tags TAG1,TAG2",
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      0,
									Character: 0,
								},
								End: protocol.Position{
									Line:      0,
									Character: 1,
								},
							},
						},
					},
				},
			},
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := GetTagCompletionItems("", tt.args.position)
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
