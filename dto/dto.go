package dto

type RegularExpression struct {
	RE string `json:"regular_expression"`
}

const Epsilon = 2

type ENFAResponse struct {
	TransitionTableSize int
}

type TransitionKey struct {
	SourceState int
	InputSymbol string
}

type StateSet map[int]bool

type MetricName string

var (
	TotalRequests             MetricName = "http_requests_total"
	RequestDuration           MetricName = "http_request_duration_seconds"
	RegexSize                 MetricName = "regex_size_total"
	EnfaTransitionTableSize   MetricName = "enfa_transition_table_size"
	RegexProcessedTotal       MetricName = "regex_processed_total"
	ProcessingTimeByRegexSize MetricName = "processing_time_by_regex_size"
)

type Edge struct {
	Src   int
	Input int
	Dst   int
}

type Closure struct {
	Src int
	Dst int
}

type ResponseFormat struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type Login struct {
	ID       int    `json:"id"`
	Password string `json:"password"`
}

type TransitionTable struct {
	TransitionTable []map[string]string `json:"transition_table"`
}
