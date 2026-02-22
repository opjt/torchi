package push

import "sync"

type WaitMap struct {
	mu sync.Mutex
	m  map[string]chan string
}

func NewWaitMap() *WaitMap {
	return &WaitMap{m: make(map[string]chan string)}
}

func (w *WaitMap) Set(id string, ch chan string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.m[id] = ch
}

func (w *WaitMap) Get(id string) (chan string, bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	ch, ok := w.m[id]
	return ch, ok
}

func (w *WaitMap) Delete(id string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.m, id)
}
