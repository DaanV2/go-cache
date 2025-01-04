package test_util

import (
	"fmt"

	"github.com/daanv2/go-cache/pkg/constraints"
	"github.com/daanv2/go-cache/pkg/hash"
)

var _ constraints.Equivalent[*TestItem] = &TestItem{}

type TestItem struct {
	ID   int
	Data string
}

// Equal implements constraints.Equivalent.
func (t *TestItem) Equal(other *TestItem) bool {
	if t == nil || other == nil {
		return t == other
	}

	return t.ID == other.ID
}

func (t *TestItem) String() string {
	return fmt.Sprintf("TestItem{ID: %v, Data: %v}", t.ID, t.Data)
}

func NewItem(id int) *TestItem {
	return &TestItem{
		ID:   id,
		Data: fmt.Sprintf("id=%v", id),
	}
}

func Generate(amount int) []*TestItem {
	items := make([]*TestItem, 0, amount)

	for id := range amount {
		items = append(items, NewItem(id))
	}

	return items
}

var _ hash.Hasher[*TestItem] = &TestItemHasher{}

type TestItemHasher struct {
	base intHasher[int]
}

func (t *TestItemHasher) Hash(item *TestItem) uint64 {
	return t.base.Hash(item.ID)
}

func Hasher() *TestItemHasher {
	return &TestItemHasher{
		base: intHasher[int]{},
	}
}
