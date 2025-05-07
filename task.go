package zerodown

import (
	"sync"
)

type task struct {
	handler func() error
	async   bool
}

type taskManager struct {
	startupHooks  []task
	shutdownHooks []task
	mu            sync.RWMutex
}

func newTaskManager() *taskManager {
	return &taskManager{
		startupHooks:  make([]task, 0),
		shutdownHooks: make([]task, 0),
	}
}

func (tm *taskManager) executeStartupHooks() error {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	var wg sync.WaitGroup
	errChan := make(chan error, len(tm.startupHooks))

	for _, hook := range tm.startupHooks {
		if hook.async {
			wg.Add(1)
			go func(h task) {
				defer wg.Done()
				if err := h.handler(); err != nil {
					errChan <- err
				}
			}(hook)
		} else {
			if err := hook.handler(); err != nil {
				return err
			}
		}
	}

	wg.Wait()
	close(errChan)

	// 收集异步错误
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}

func (tm *taskManager) executeShutdownHooks() error {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	var wg sync.WaitGroup
	for _, hook := range tm.shutdownHooks {
		if hook.async {
			wg.Add(1)
			go func(h task) {
				defer wg.Done()
				_ = h.handler()
			}(hook)
		} else {
			_ = hook.handler()
		}
	}
	wg.Wait()
	return nil
}

func (tm *taskManager) addStartupHook(handler func() error, async bool) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.startupHooks = append(tm.startupHooks, task{handler: handler, async: async})
}

func (tm *taskManager) addShutdownHook(handler func() error, async bool) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.shutdownHooks = append(tm.shutdownHooks, task{handler: handler, async: async})
}
