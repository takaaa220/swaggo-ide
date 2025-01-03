package protocol

// file:///xxx/yyy/zzz
type DocumentUri string

type TextDocumentIdentifier struct {
	Uri DocumentUri `json:"uri"`
}

type Position struct {
	Line      uint32 `json:"line"`      // 0-based
	Character uint32 `json:"character"` // 0-based
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Location struct {
	Uri   DocumentUri `json:"uri"`
	Range Range       `json:"range"`
}

type ProgressToken any // string | int
