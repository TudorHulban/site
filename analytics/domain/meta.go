package domain

import "errors"

var ErrKeyNotFound = errors.New("key not found")

type Meta[T comparable] struct {
	Names  [7]T
	Values [7]uint32
}

func (m *Meta[T]) GetValue(byKey T) (uint32, error) {
	for ix, n := range m.Names {
		if n == byKey {
			return m.Values[ix], nil
		}
	}

	return 0,
		ErrKeyNotFound
}

func (m *Meta[T]) Increment(key T) {
	// 1. Scan for exact match
	for i, name := range m.Names {
		if name == key {
			m.Values[i]++
			return
		}
	}

	// 2. Look for an empty slot (using your zeroT check safely now)
	var zeroT T
	for i, name := range m.Names {
		if name == zeroT {
			m.Names[i] = key
			m.Values[i] = 1
			return
		}
	}

	// 3. Matrix is full: Find the absolute weakest link
	lowestIdx := 0
	lowestVal := m.Values[0]
	for i := 1; i < 7; i++ {
		if m.Values[i] < lowestVal {
			lowestVal = m.Values[i]
			lowestIdx = i
		}
	}

	// 4. Space-Saving Eviction
	// We evict the old loser, install the new key, and increment the counter.
	// This ensures a massive spike will instantly overtake the matrix.
	m.Names[lowestIdx] = key
	m.Values[lowestIdx] = lowestVal + 1
}
