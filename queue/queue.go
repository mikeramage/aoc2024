package queue

type Queue[T any] interface {
	Append(v T)
	PopFront() T
	Len() int
}

// Queue implementation
type QueueImpl[T any] struct {
	elements []T
}

func (q *QueueImpl[T]) Append(v T) {
	q.elements = append(q.elements, v)
}

func (q *QueueImpl[T]) Len() int {
	return len(q.elements)
}

func (q *QueueImpl[T]) PopFront() T {
	if len(q.elements) == 0 {
		panic("Cannot pop from empty queue!")
	}
	var el T
	el, q.elements = q.elements[0], q.elements[1:]

	return el
}

func NewQueue[T any]() Queue[T] {
	return &QueueImpl[T]{elements: make([]T, 0)}
}
