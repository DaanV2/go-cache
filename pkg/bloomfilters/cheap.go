package bloomfilters

const (
	size_uint64 = 8
	diffuser    = 0x47b5481dbefa4fa4
)

type Cheap struct {
	amount uint64   // Amount of items stored
	words  []uint64 // Bit sizes
}

func NewCheap(amount uint64) *Cheap {
	w := max(amount/size_uint64, 1) * 2
	if (amount % size_uint64) != 0 {
		w++
	}

	return &Cheap{
		words:  make([]uint64, w),
		amount: amount,
	}
}

func (c *Cheap) Has(hash uint64) bool {
	return c.has(hash) || c.has(hash^diffuser)
}

func (c *Cheap) has(hash uint64) bool {
	bucket, bit := index(hash, c.amount)

	word := c.words[bucket]
	mask := uint64(1 << bit)
	return (word & mask) == mask
}

func (c *Cheap) Set(hash uint64) {
	c.set(hash)
	c.set(hash ^ diffuser)
}

func (c *Cheap) set(hash uint64) {
	bucket, bit := index(hash, c.amount)

	v := uint64(1 << bit)
	c.words[bucket] |= v
}

func index(hash, amount uint64) (bucket, bit uint64) {
	bitIndex := hash % amount

	bucket = bitIndex / size_uint64
	bit = bitIndex % size_uint64
	return bucket, bit
}
