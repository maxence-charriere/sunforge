package conduit

import "sync"

type queue[T any] struct {
	mutex sync.Mutex
	data  []T
}

func newQueue[T any](capacity int) *queue[T] {
	return &queue[T]{
		data: make([]T, 0, capacity),
	}
}

func (q *queue[T]) Enqueue(v T) {
	q.mutex.Lock()
	q.data = append(q.data, v)
	q.mutex.Unlock()
}

func (q *queue[T]) Dequeue() (T, bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	var zero T
	if len(q.data) == 0 {
		return zero, false
	}

	v := q.data[0]
	copy(q.data, q.data[1:])
	q.data = q.data[:len(q.data)-1]

	return v, true
}

func (q *queue[T]) Len() int {
	q.mutex.Lock()
	l := len(q.data)
	q.mutex.Unlock()
	return l
}

func (q *queue[T]) Cap() int {
	q.mutex.Lock()
	c := cap(q.data)
	q.mutex.Unlock()
	return c
}
