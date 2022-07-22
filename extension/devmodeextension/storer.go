package devmodeextension

type Storer interface {
	StoreSpan(span Span)
}
