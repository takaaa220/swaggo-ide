package protocol

import "github.com/takaaa220/swaggo-ide/server/internal/server/transport"

type CompletionFunc func(transport.Context, *CompletionParams) (CompletionResult, error)

// CompletionOptions represents options for completion requests.
type CompletionOptions struct {
	TriggerCharacters   []string             `json:"triggerCharacters,omitempty"`
	AllCommitCharacters []string             `json:"allCommitCharacters,omitempty"`
	ResolveProvider     bool                 `json:"resolveProvider,omitempty"`
	WorkDoneProgress    bool                 `json:"workDoneProgress,omitempty"`
	CompletionItem      CompletionItemOption `json:"completionItem,omitempty"`
}

type CompletionItemOption struct {
	LabelDetailsSupport bool `json:"labelDetailsSupport,omitempty"`
}

type CompletionParams struct {
	TextDocument       TextDocumentIdentifier `json:"textDocument"`
	Position           Position               `json:"position"`
	WorkDoneToken      ProgressToken          `json:"workDoneToken,omitempty"`
	PartialResultToken ProgressToken          `json:"partialResultToken,omitempty"`
	Context            CompletionContext      `json:"context,omitempty"`
}

type CompletionContext struct {
	TriggerKind      CompletionTriggerKind `json:"triggerKind"`
	TriggerCharacter string                `json:"triggerCharacter,omitempty"`
}

type CompletionTriggerKind int

const (
	CompletionTriggerKindInvoked                         CompletionTriggerKind = 1
	CompletionTriggerKindCharacter                                             = 2
	CompletionTriggerKindTriggerForIncompleteCompletions                       = 3
)

type CompletionResult any // []CompletionItem | CompletionList | null

type CompletionList struct {
	IsIncomplete bool                   `json:"isIncomplete"`
	Items        []CompletionItem       `json:"items"`
	ItemDefaults CompletionItemDefaults `json:"itemDefaults,omitempty"`
}

type CompletionItemDefaults struct {
	CommitCharacters string           `json:"commitCharacters,omitempty"`
	EditRange        Range            `json:"editRange,omitempty"`
	InsertTextFormat InsertTextFormat `json:"insertTextFormat,omitempty"`
	InsertTextMode   InsertTextMode   `json:"insertTextMode,omitempty"`
	Data             any              `json:"data,omitempty"`
}

type InsertTextFormat int

const (
	InsertTextFormatPlainText InsertTextFormat = 1
	InsertTextFormatSnippet                    = 2
)

type InsertTextMode int

const (
	InsertTextModeAsIs              InsertTextMode = 1
	InsertTextModeAdjustIndentation                = 2
)

type CompletionItem struct {
	Label    string              `json:"label"`
	Kind     CompletionItemKind  `json:"kind,omitempty"`
	Tags     []CompletionItemTag `json:"tags,omitempty"`
	TextEdit any                 `json:"textEdit,omitempty"` // TextEdit | InsertReplaceEdit
}

type CompletionItemKind int

const (
	CompletionItemKindText          CompletionItemKind = 1
	CompletionItemKindMethod        CompletionItemKind = 2
	CompletionItemKindFunction      CompletionItemKind = 3
	CompletionItemKindConstructor   CompletionItemKind = 4
	CompletionItemKindField         CompletionItemKind = 5
	CompletionItemKindVariable      CompletionItemKind = 6
	CompletionItemKindClass         CompletionItemKind = 7
	CompletionItemKindInterface     CompletionItemKind = 8
	CompletionItemKindModule        CompletionItemKind = 9
	CompletionItemKindProperty      CompletionItemKind = 10
	CompletionItemKindUnit          CompletionItemKind = 11
	CompletionItemKindValue         CompletionItemKind = 12
	CompletionItemKindEnum          CompletionItemKind = 13
	CompletionItemKindKeyword       CompletionItemKind = 14
	CompletionItemKindSnippet       CompletionItemKind = 15
	CompletionItemKindColor         CompletionItemKind = 16
	CompletionItemKindFile          CompletionItemKind = 17
	CompletionItemKindReference     CompletionItemKind = 18
	CompletionItemKindFolder        CompletionItemKind = 19
	CompletionItemKindEnumMember    CompletionItemKind = 20
	CompletionItemKindConstant      CompletionItemKind = 21
	CompletionItemKindStruct        CompletionItemKind = 22
	CompletionItemKindEvent         CompletionItemKind = 23
	CompletionItemKindOperator      CompletionItemKind = 24
	CompletionItemKindTypeParameter CompletionItemKind = 25
)

type CompletionItemTag int

const (
	CompletionItemTagDeprecated CompletionItemTag = 1
)

type TextEdit struct {
	NewText string `json:"newText"`
	Range   Range  `json:"range"`
}

type InsertReplaceEdit struct {
	NewText string `json:"newText"`
	Insert  Range  `json:"insert"`
	Replace Range  `json:"replace"`
}
