package ants

import "sync"

// PanicHandler 线程池 Panic 处理器
type PanicRecover interface {
	Handle(err error, w *Worker)
}

type RestartPanicRecover struct {
	N int
	n int

	lock sync.Mutex
}

func (self *RestartPanicRecover) Handle(err error, w *Worker) {
	if self.n < self.N {
		if fun, ok := w.pool.jobs.Load(w.Tag); ok {
			w.pool.Submit(fun.(f))
			self.lock.Lock()
			self.n++
			self.lock.Unlock()
		}
	}
}
