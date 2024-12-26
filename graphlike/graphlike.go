package graphlike

type Graphlike[V] interface {
	NeightboursOf(V) []V
}
