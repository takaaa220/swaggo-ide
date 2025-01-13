package swag

type Range struct {
	Start Position
	End   Position
}

type Position struct {
	Line      uint32
	Character uint32
}

type CompletionCandidate struct {
	Label   string
	NewText string
}
