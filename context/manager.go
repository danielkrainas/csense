package context

import (
	"net/http"
	"sync"
)

type Manager struct {
	contexts map[*http.Request]Context
	mutex    sync.Mutex
}

var DefaultContextManager = NewManager()

func NewManager() *Manager {
	return &Manager{
		contexts: make(map[*http.Request]Context),
	}
}

func (m *Manager) Context(parent Context, w http.ResponseWriter, r *http.Request) Context {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if ctx, ok := m.contexts[r]; ok {
		return ctx
	}

	if parent == nil {
		parent = Background()
	}

	ctx := WithRequest(parent, r)
	ctx, w = WithResponseWriter(ctx, w)
	ctx = WithLogger(ctx, GetLogger(ctx))
	m.contexts[r] = ctx
	return ctx
}

func (m *Manager) Release(ctx Context) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	r, err := GetRequest(ctx)
	if err != nil {
		GetLogger(ctx).Error("no request found in context at release")
		return
	}

	delete(m.contexts, r)
}
