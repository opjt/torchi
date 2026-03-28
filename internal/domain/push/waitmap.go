package push

import "sync"

type WaitResult struct {
	Reaction string
	Deleted  bool
}

type WaitMap struct {
	mu sync.Mutex
	m  map[string]chan WaitResult
}

func NewWaitMap() *WaitMap {
	return &WaitMap{m: make(map[string]chan WaitResult)}
}

func (w *WaitMap) Set(id string, ch chan WaitResult) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.m[id] = ch
}

func (w *WaitMap) Get(id string) (chan WaitResult, bool) {
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
