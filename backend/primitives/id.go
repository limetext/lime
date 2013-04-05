package primitives

import (
	"sync"
)

var (
	idCount = Id(0)
	idMutex sync.Mutex
)

type (
	Id    int
	HasId struct {
		id Id
	}
)

func (i *HasId) Id() Id {
	if i.id == 0 {
		i.id = nextId()
	}
	return i.id
}

func nextId() Id {
	idMutex.Lock()
	defer idMutex.Unlock()
	idCount++
	return idCount
}
