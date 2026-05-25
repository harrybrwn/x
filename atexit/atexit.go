// Package atexit allows users to register callbacks that will be run at exit.
package atexit

import (
	"sync"
)

type Handler func() error

var (
	lock     sync.Mutex
	handlers = make([]Handler, 0)
	once     sync.Once
)

func Add(fn Handler) {
	lock.Lock()
	handlers = append(handlers, fn)
	lock.Unlock()
}

func RunAll() (err error) {
	lock.Lock()
	defer lock.Unlock()
	once.Do(func() { err = run() })
	return err
}

func run() (err error) {
	for _, fn := range handlers {
		e := fn()
		if e != nil && err == nil {
			err = e
		}
	}
	return err
}
