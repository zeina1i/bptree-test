package bptree

type IteratorOptions struct {
	reverse bool
}

// Iterator returns a stateful Iterator for traversing the tree
// in ascending key order.
type Iterator struct {
	next    *node
	i       int
	options *IteratorOptions
}

// Iterator returns a stateful iterator that traverses the tree
// in ascending key order.
func (t *BPTree) Iterator(options *IteratorOptions) *Iterator {
	var i int
	if options.reverse {
		i = t.rightmost.keyNum - 1
		return &Iterator{t.rightmost, i, options}
	} else {
		i = 0
		return &Iterator{t.leftmost, i, options}
	}
}

// HasNext returns true if there is a next element to retrive.
func (it *Iterator) HasNext() bool {
	if it.options.reverse {
		return it.next != nil && it.i >= 0
	}

	return it.next != nil && it.i < it.next.keyNum
}

// Next returns a key and a value at the current position of the iteration
// and advances the iterator.
// Caution! Next panics if called on the nil element.
func (it *Iterator) Next() ([]byte, []byte) {
	if !it.HasNext() {
		// to sleep well
		panic("there is no next node")
	}

	if !it.options.reverse {
		key, value := it.next.keys[it.i], it.next.pointers[it.i].asValue()

		it.i++
		if it.i == it.next.keyNum {
			nextPointer := it.next.next(false)
			if nextPointer != nil {
				it.next = nextPointer.asNode()
			} else {
				it.next = nil
			}

			it.i = 0
		}

		return key, value
	} else {
		key, value := it.next.keys[it.i], it.next.pointers[it.i].asValue()

		it.i--
		if it.i == -1 {
			nextPointer := it.next.next(true)
			if nextPointer != nil {
				it.next = nextPointer.asNode()
				it.i = nextPointer.asNode().keyNum - 1
			} else {
				it.i = 0
				it.next = nil
			}
		}

		return key, value
	}
}
