package util

// IntValue covers all integer types supported by [IntPool].
type IntValue interface {
	int | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
}

// IntPool is an implementation using implicit linked lists.
// Implements https://skypjack.github.io/2019-05-06-ecs-baf-part-3/
type IntPool[T IntValue] struct {
	pool              []T
	next              T
	available         uint32
	capacityIncrement uint32
}

// NewIntPool creates a new, initialized Entity pool.
func NewIntPool[T IntValue](capacityIncrement uint32) IntPool[T] {
	return IntPool[T]{
		pool:              make([]T, 0, capacityIncrement),
		next:              0,
		available:         0,
		capacityIncrement: capacityIncrement,
	}
}

// Get returns a fresh or recycled entity.
func (p *IntPool[T]) Get() T {
	if p.available == 0 {
		return p.getNew()
	}
	curr := p.next
	p.next, p.pool[p.next] = p.pool[p.next], p.next
	p.available--
	return p.pool[curr]
}

// Allocates and returns a new entity. For internal use.
func (p *IntPool[T]) getNew() T {
	e := T(len(p.pool))
	if len(p.pool) == cap(p.pool) {
		old := p.pool
		p.pool = make([]T, len(p.pool), len(p.pool)+int(p.capacityIncrement))
		copy(p.pool, old)
	}
	p.pool = append(p.pool, e)
	return e
}

// Recycle hands an entity back for recycling.
func (p *IntPool[T]) Recycle(e T) {
	p.next, p.pool[e] = e, p.next
	p.available++
}

// Reset recycles all entities. Does NOT free the reserved memory.
func (p *IntPool[T]) Reset() {
	p.pool = p.pool[:0]
	p.next = 0
	p.available = 0
}
