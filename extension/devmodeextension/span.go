package devmode

import "database/sql"

//type Span struct {
//	SpanID             string
//	TraceID            string
//	ParentID           string
//	StartTime          int
//	EndTime            int
//	attributes         []Attribute
//	resourceAttributes []Attribute
//}
//
//type Attribute struct {
//	Key   string
//	Value string
//}

type Span struct {
	SpanID             sql.NullString
	TraceID            sql.NullString
	ParentID           sql.NullString
	StartTime          sql.NullInt64
	EndTime            sql.NullInt64
	attributes         []Attribute
	resourceAttributes []Attribute
}

type Attribute struct {
	Key   sql.NullString
	Value sql.NullString
}
