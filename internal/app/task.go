package app

import (
	"context"
	"errors"
	"sync"
)

type TaskFunc func(ctx context.Context) error

type TaskManager struct {
	mu      sync.Mutex
	cancels map[string]context.CancelFunc
}

func NewTaskManager() *TaskManager {
	return &TaskManager{cancels: map[string]context.CancelFunc{}}
}

func (m *TaskManager) Start(id string, fn TaskFunc) {
	m.mu.Lock()
	ctx, cancel := context.WithCancel(context.Background())
	m.cancels[id] = cancel
	m.mu.Unlock()

	go func() {
		_ = fn(ctx)
		m.mu.Lock()
		delete(m.cancels, id)
		m.mu.Unlock()
	}()
}

func (m *TaskManager) Cancel(id string) error {
	m.mu.Lock()
	cancel, ok := m.cancels[id]
	m.mu.Unlock()
	if !ok {
		return errors.New("task not found")
	}
	cancel()
	return nil
}
