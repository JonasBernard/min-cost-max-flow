package graphlike

type Graphlike[V any] interface {
	NeightboursOf(V) []V
}
