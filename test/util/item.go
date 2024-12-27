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

// Equals implements constraints.Equivalent.
func (t *TestItem) Equals(other *TestItem) bool {
	return t.ID == other.ID
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

type TestItemHasher struct{}

func (t *TestItemHasher) Hash(item *TestItem) uint64 {
	return uint64(item.ID)
}

func Hasher() *TestItemHasher {
	return &TestItemHasher{}
}
