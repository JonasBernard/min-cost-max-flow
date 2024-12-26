package graphlike

type RootedTreelike[V any] interface {
	Graphlike[V]
	Root() T
	Parent(V) V
	Children(V) []V
	IsLeaf(V) bool
}