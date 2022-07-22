package devmodeextension

type Span struct {
	Name                  string            `json:"name"`
	SpanID                string            `json:"span_id"`
	TraceID               string            `json:"trace_id"`
	ParentID              string            `json:"parent_id,omitempty"`
	StartTime             int64             `json:"start_time"`
	EndTime               int64             `json:"end_time"`
	Attributes            string            `json:",omitempty"`
	ResourceAttributes    string            `json:",omitempty"`
	AttributesMap         map[string]string `json:"attributes,omitempty"`
	ResourceAttributesMap map[string]string `json:"resource_attributes,omitempty"`
}

//type Attribute struct {
//	Key   sql.NullString
//	Value sql.NullString
//}
