package devmodeextension

import "database/sql"

type Span struct {
	SpanID             sql.NullString
	TraceID            sql.NullString
	ParentID           sql.NullString
	StartTime          sql.NullInt64
	EndTime            sql.NullInt64
	Attributes         sql.NullString
	ResourceAttributes sql.NullString
}

//type Attribute struct {
//	Key   sql.NullString
//	Value sql.NullString
//}
