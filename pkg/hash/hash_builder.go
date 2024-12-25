package hash

type HashBuilder interface {
	Write(data []byte) error
	Sum() uint64
}
