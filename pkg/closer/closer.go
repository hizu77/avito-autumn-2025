package closer

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"slices"
	"sync"
	"sync/atomic"
	"syscall"

	"golang.org/x/sync/errgroup"
)

var globalCloser = newCloser(syscall.SIGINT, syscall.SIGTERM)

type closeFunc func() error

type closer struct {
	mu     *sync.Mutex
	closed atomic.Bool

	ch      chan os.Signal
	errorCh chan error

	priorityByGroup map[string]int
	groupCallbacks  map[string][]closeFunc
}

type Group struct {
	Name     string
	Priority int
}

func newCloser(signals ...os.Signal) *closer {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)

	errorCh := make(chan error, 1)

	closer := &closer{
		mu:              &sync.Mutex{},
		priorityByGroup: make(map[string]int),
		groupCallbacks:  make(map[string][]closeFunc),
		ch:              ch,
		errorCh:         errorCh,
	}

	go func() {
		defer signal.Stop(ch)

		// wait for passed signals
		<-ch

		// if a signal was given, start closing groups in order
		errorCh <- closer.CloseAll()
		close(errorCh)
	}()

	return closer
}

func AddSignals(signals ...os.Signal) {
	signal.Notify(globalCloser.ch, signals...)
}

// AddGroups adds groups with given priority.
// The lower the priority, the earlier the group of function will execute.
func AddGroups(groups ...Group) {
	globalCloser.mu.Lock()
	defer globalCloser.mu.Unlock()

	for _, group := range groups {
		globalCloser.priorityByGroup[group.Name] = group.Priority
	}
}

// AddCallback adds a callback to provided group.
// If the group with passed name does not exist, function returns an error.
func AddCallback(groupName string, callback closeFunc) error {
	globalCloser.mu.Lock()
	defer globalCloser.mu.Unlock()

	if _, ok := globalCloser.priorityByGroup[groupName]; !ok {
		return ErrGroupNotFound
	}

	globalCloser.groupCallbacks[groupName] = append(
		globalCloser.groupCallbacks[groupName],
		callback,
	)

	return nil
}

// CloseAll calls all callbacks from groups in ascending priority order.
// It is called automatically when signal is received but it may be called manually.
// Groups with the same priority are executed in parallel.
func (c *closer) CloseAll() error {
	if !c.closed.CompareAndSwap(false, true) {
		return ErrAlreadyClosed
	}

	var allErrors []error

	groupsByPriority := make(map[int][]string, len(c.priorityByGroup))

	for name, priority := range c.priorityByGroup {
		groupsByPriority[priority] = append(groupsByPriority[priority], name)
	}

	priorities := make([]int, 0, len(groupsByPriority))

	for priority := range groupsByPriority {
		priorities = append(priorities, priority)
	}

	slices.SortFunc(
		priorities,
		func(lhs int, rhs int) int {
			return lhs - rhs
		},
	)

	for _, priority := range priorities {
		g, _ := errgroup.WithContext(context.Background())

		for _, group := range groupsByPriority[priority] {
			for _, callback := range c.groupCallbacks[group] {
				g.Go(callback)
			}
		}

		if err := g.Wait(); err != nil {
			allErrors = append(allErrors, err)
		}
	}

	return errors.Join(allErrors...)
}

// Wait waiting until all callbacks are executed
// i.e the channel with the error receives a value
func Wait() error {
	return <-globalCloser.errorCh
}
