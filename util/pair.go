package util

/**
键值对
*/
type Pair struct {
	K interface{}
	V interface{}
}

func NewPair(K, V interface{}) *Pair {
	return &Pair{K: K, V: V}
}
