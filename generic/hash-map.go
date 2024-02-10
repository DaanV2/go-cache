package generic



type KVPair[K, V any] struct {
	Hashcode int
	Key K
	Value V
}

func NewKVPair[K, V any](key K, value V) *KVPair[K, V] {
	return &KVPair[K, V]{
		Hashcode: 0,
		Key: key,
		Value: value,
	}
}


type HashMap[K, V any] struct {

}
