package priorityqueue

import "container/heap"

// Adapted from Go example here: https://pkg.go.dev/container/heap
// This implementation maintains an internal map of value to index for
// quick update by value rather than index. The expected use case assumes all
// values are unique, but if not, it will update the most recently added/updated
// item with the requested value.

type Item[T any] struct {
	Value    T
	Priority int
	Index    int
}

type PriorityQueue[T any] struct {
	items []*Item[T]
}

func NewPriorityQueue[T any]() *PriorityQueue[T] {
	return &PriorityQueue[T]{
		items: []*Item[T]{},
	}
}

func (pq PriorityQueue[T]) Len() int {
	return len(pq.items)
}

func (pq PriorityQueue[T]) Less(i, j int) bool {
	// Lower priority values are better!
	return pq.items[i].Priority < pq.items[j].Priority
}

func (pq PriorityQueue[T]) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].Index = i
	pq.items[j].Index = j
}

func (pq *PriorityQueue[T]) PushItem(item *Item[T]) {
	heap.Push(pq, item)
}

func (pq *PriorityQueue[T]) PopItem() *Item[T] {
	popped := heap.Pop(pq)
	if popped == nil {
		return nil
	}
	return popped.(*Item[T])
}

// Annoying. I want this to take a generic, but have to pass any
func (pq *PriorityQueue[T]) Push(x any) {
	item, ok := x.(*Item[T])
	if !ok {
		panic("Push expects *Item[T]")
	}
	n := len(pq.items)
	item.Index = n
	pq.items = append(pq.items, item)
}

func (pq *PriorityQueue[T]) Pop() any {
	old := pq.items
	n := len(old)
	if n == 0 {
		return nil
	}
	item := old[n-1]
	old[n-1] = nil //Allow GC
	pq.items = old[:n-1]
	return item
}

// Finds item by index and updates the value and priority
func (pq *PriorityQueue[T]) Update(item *Item[T], value T, priority int) {
	item.Value = value
	item.Priority = priority
	heap.Fix(pq, item.Index)
}

// func (pq *PriorityQueue[T]) GetItem(value T) *Item[T] {
// 	index, exists := pq.indexMap[value]
// 	if !exists {
// 		return nil
// 	}
// 	item := pq.items[index]
// 	return item
// }
