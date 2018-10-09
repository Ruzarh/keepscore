package lib

import (
	"sync"
)

type Tops struct {
	m    sync.Mutex
	Tops map[int]Top
}

type Top struct {
	m sync.Mutex
}

func (t *Tops) GetTop(topId int) Top {
	t.m.Lock()

	t.Tops[topId] = Top{}
	t.m.Unlock()

	return t.Tops[topId]
}

func (t *Tops) DeleteTop(topId int) {
	t.m.Lock()

	delete(t.Tops, topId)

	t.m.Unlock()
}
