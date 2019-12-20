package wgext

import "sync"

type WaitGroup struct {
	wg       sync.WaitGroup
	errs     chan error
	finished chan bool
}

func NewWaitGroup() *WaitGroup {
	return &WaitGroup{
		wg:       sync.WaitGroup{},
		errs:     make(chan error),
		finished: make(chan bool, 1),
	}
}

func (wg *WaitGroup) Add(delta int) {
	wg.wg.Add(delta)
}

func (wg *WaitGroup) Done() {
	wg.wg.Done()
}

func (wg *WaitGroup) Fail(err error) {
	wg.errs <- err
	wg.wg.Done()
}

func (wg *WaitGroup) Wait() error {

	go func() {
		wg.wg.Wait()
		wg.finished <- true
	}()

	select {
	case err := <-wg.errs:
		return err
	case <-wg.finished:
		return nil
	}
}
