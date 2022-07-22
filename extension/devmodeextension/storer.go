package devmodeextension

type Storer interface {
	StoreTrace(span Span)
}
