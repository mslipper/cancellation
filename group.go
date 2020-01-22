package cancellation

import "sync"

// Group represents a group of cancellable goroutines. For example:
//
//     var g Group
//     g.Do(func(k <-chan struct{}) {
//	     for {
//		       select {
//			       case <-k:
//				       return
//	               case <-someOtherThing:
//	                   // do some other thing...
//	         }
//	     }
//     })
//     g.Cancel() // causes the goroutine above to exit
//     g.Wait() // waits for the goroutine above to exit
//
// This is useful for managing the lifecycle of multiple long-running
// goroutines.
type Group struct {
	mtx     sync.Mutex
	members []chan struct{}
	done    bool
	wg      sync.WaitGroup
}

// Do spawns a new goroutine under the control of the Group.
// It is the responsibility of the spawned goroutine to handle
// the passed-in killed channel. Calling Do after the Group is
// cancelled will panic.
func (g *Group) Do(f func(killed <-chan struct{})) {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	if g.done {
		panic("group is cancelled")
	}
	killCh := make(chan struct{}, 1)
	g.members = append(g.members, killCh)
	g.wg.Add(1)
	go func() {
		f(killCh)
		g.wg.Done()
	}()
}

// Cancel cancels all goroutines under the control of the Group.
// Since the kill channel passed to each controlled goroutine is
// buffered, each routine is killed asynchronously. Do not rely
// on this method to determine whether or not all goroutines have
// returned. Instead, use Wait(). Cancel will panic if called more
// than once.
func (g *Group) Cancel() {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	if g.done {
		panic("group is cancelled")
	}
	for _, killCh := range g.members {
		killCh <- struct{}{}
	}
	g.done = true
}

// Wait blocks until all goroutines under the control of the Group
// have returned.
func (g *Group) Wait() {
	g.wg.Wait()
}
