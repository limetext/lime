package backend

import (
	"sync"
)

var (
	idCount = 0
	idMutex sync.Mutex
)

type HasId struct {
	id int
}

func (i *HasId) Id() int {
	if i.id == 0 {
		i.id = nextId()
	}
	return i.id
}

func nextId() int {
	idMutex.Lock()
	defer idMutex.Unlock()
	idCount++
	return idCount
}
