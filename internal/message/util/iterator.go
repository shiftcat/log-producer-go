package util

type Iterator[T any] interface {
	HasNext() bool
	Next() (int, T)
}

func NewIterator[T any](items []T) Iterator[T] {
	return &iter[T]{pos: -1, items: items, size: len(items)}
}

type iter[T any] struct {
	pos   int
	items []T
	size  int
}

func (i *iter[T]) HasNext() bool {
	if (i.pos + 1) < i.size {
		return true
	}
	return false
}

func (i *iter[T]) Next() (int, T) {
	if i.size <= i.pos {
		i.pos = i.size - 1
	} else {
		i.pos = i.pos + 1
	}
	return i.pos, i.items[i.pos]
}
