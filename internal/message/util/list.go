package util

import (
	"fmt"
	"sync"
)

type List[T any] interface {
	Size() int
	Add(i T)
	AddAll(items []T)
	Remove(i int)
	Get(i int) T
	GetAll() []T
	Take(i int) T
	Iterator() Iterator[T]
}

func NewList[T any](items []T) List[T] {
	l := new(list[T])
	l.items = items
	return l
}

type list[T any] struct {
	items []T
	sync.RWMutex
}

func (l *list[T]) Add(i T) {
	l.Lock()
	l.items = append(l.items, i)
	l.Unlock()
}

func (l *list[T]) AddAll(items []T) {
	l.Lock()
	l.items = append(l.items, items...)
	l.Unlock()
}

func (l *list[T]) removeItem(i int) {
	var size = len(l.items)
	if (size - 1) == i {
		l.items = append(l.items[:i])
	} else {
		l.items = append(l.items[:i], append(l.items[i+1:])...)
	}
}

func (l *list[T]) Remove(i int) {
	l.Lock()
	var size = len(l.items)
	if (size - 1) < i {
		panic(fmt.Sprint("Invalid index number: ", i))
	}
	l.removeItem(i)
	l.Unlock()
}

func (l *list[T]) Take(i int) T {
	l.Lock()
	var size = len(l.items)
	if size <= i {
		panic(fmt.Sprint("Invalid index number: ", i))
	}
	item := l.items[i]
	l.removeItem(i)
	l.Unlock()
	return item
}

func (l *list[T]) Get(i int) T {
	l.RLock()
	var size = len(l.items)
	if size <= i {
		panic(fmt.Sprint("Invalid index number: ", i))
	}
	item := l.items[i]
	l.RUnlock()
	return item
}

func (l *list[T]) GetAll() []T {
	l.RLock()
	rs := append([]T{}, l.items...)
	l.RUnlock()
	return rs
}

func (l *list[T]) Size() int {
	l.RLock()
	size := len(l.items)
	l.RUnlock()
	return size
}

func (l *list[T]) Iterator() Iterator[T] {
	l.RLock()
	copyItems := append([]T{}, l.items...)
	l.RUnlock()
	return NewIterator(copyItems)
}
