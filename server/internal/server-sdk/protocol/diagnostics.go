package protocol

type PublishDiagnosticsParams struct {
	Uri         DocumentUri   `json:"uri"`
	Version     int           `json:"version,omitempty"`
	Diagnostics []Diagnostics `json:"diagnostics"`
}

type Diagnostics struct {
	Range              Range                          `json:"range"`
	Severity           DiagnosticsSeverity            `json:"severity,omitempty"`
	Code               string                         `json:"code,omitempty"`
	CodeDescription    CodeDescription                `json:"codeDescription,omitempty"`
	Source             string                         `json:"source,omitempty"`
	Message            string                         `json:"message"`
	Tags               []DiagnosticTag                `json:"tags,omitempty"`
	RelatedInformation []DiagnosticRelatedInformation `json:"relatedInformation,omitempty"`
	Data               any                            `json:"data,omitempty"`
}

type DiagnosticsSeverity int

const (
	DiagnosticsSeverityError DiagnosticsSeverity = 1
	DiagnosticsSeverityWarn  DiagnosticsSeverity = 2
	DiagnosticsSeverityInfo  DiagnosticsSeverity = 3
	DiagnosticsSeverityHint  DiagnosticsSeverity = 4
)

type CodeDescription struct {
	Href string `json:"href"`
}

type DiagnosticTag int

const (
	DiagnosticTagUnnecessary DiagnosticTag = 1
	DiagnosticTagDeprecated  DiagnosticTag = 2
)

type DiagnosticRelatedInformation struct {
	Location Location `json:"location"`
	Message  string   `json:"message"`
}
