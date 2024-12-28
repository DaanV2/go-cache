package large

import (
	"sync"

	"github.com/daanv2/go-locks"
	optimal "github.com/daanv2/go-optimal"
)

// SetBase is the base struct for all sets.
type SetBase struct {
	bucket_size int
	bucket_lock *sync.RWMutex
	items_lock  *locks.Pool
}

// NewSetBase creates a new instance of SetBase with the default bucket size.
func NewSetBase[T any]() SetBase {
	return SetBase{
		bucket_size: optimal.SliceSize[T](),
		bucket_lock: &sync.RWMutex{},
		items_lock:  locks.NewPool(),
	}
}
