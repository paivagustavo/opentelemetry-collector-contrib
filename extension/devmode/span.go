package devmode

type Span struct {
	SpanID             string
	TraceID            string
	ParentID           string
	StartTime          int
	EndTime            int
	attributes         []Attribute
	resourceAttributes []Attribute
}

type Attribute struct {
	Key   string
	Value string
}
