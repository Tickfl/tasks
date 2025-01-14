package semaphore

type Semaphore struct {
	ch chan struct{}
}

func NewSemaphore(limit int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, limit),
	}
}

func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.ch
}
