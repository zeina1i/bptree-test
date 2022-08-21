package bptree

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func ExampleIterator() {
	tree, _ := New()

	tree.Put([]byte("apple"), []byte("sweet"))
	tree.Put([]byte("banana"), []byte("honey"))
	tree.Put([]byte("cinnamon"), []byte("savoury"))

	banana, ok := tree.Get([]byte("banana"))
	if ok {
		fmt.Printf("banana = %s\n", string(banana))
	} else {
		fmt.Println("value for banana not found")
	}

	for it := tree.Iterator(&IteratorOptions{
		reverse: false,
	}); it.HasNext(); {
		key, value := it.Next()
		fmt.Printf("key = %s, value = %s\n", string(key), string(value))
	}

	// Output:
	// banana = honey
	// key = apple, value = sweet
	// key = banana, value = honey
	// key = cinnamon, value = savoury
}

var iteratorCases = []struct {
	key   byte
	value string
}{
	{11, "11"},
	{18, "18"},
	{7, "7"},
	{15, "15"},
	{0, "0"},
	{16, "16"},
	{14, "14"},
	{33, "33"},
	{25, "25"},
	{42, "42"},
	{60, "60"},
	{2, "2"},
	{1, "1"},
	{74, "74"},
}

func TestIterator(t *testing.T) {
	tree, _ := New()
	for _, c := range iteratorCases {
		tree.Put([]byte{c.key}, []byte(c.value))
	}

	actual := make([]byte, 0)
	for it := tree.Iterator(&IteratorOptions{reverse: false}); it.HasNext(); {
		key, _ := it.Next()
		actual = append(actual, key...)
	}

	isSorted := sort.SliceIsSorted(actual, func(i, j int) bool {
		return actual[i] < actual[j]
	})
	if !isSorted {
		t.Fatalf("each does not traverse in sorted order, produced result: %s", actual)
	}

	expected := make([]byte, 0)
	for _, c := range iteratorCases {
		expected = append(expected, c.key)
	}
	sort.Slice(expected, func(i, j int) bool {
		return expected[i] < expected[j]
	})

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("%s != %s", expected, actual)
	}
}

func TestIteratorForEmptyTree(t *testing.T) {
	tree, _ := New()

	for it := tree.Iterator(&IteratorOptions{reverse: false}); it.HasNext(); {
		t.Fatal("call is not expected")
	}
}

func TestIteratorNextPanicForEmptyTree(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Next must panic on the empty tree")
		}
	}()

	tree, _ := New()

	tree.Iterator(&IteratorOptions{reverse: false}).Next()
}

func TestIteratorNextPanicAfterIteration(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Next must panic after the iteration is finished")
		}
	}()

	tree, _ := New()
	tree.Put([]byte{1}, nil)

	it := tree.Iterator(&IteratorOptions{
		reverse: true,
	})
	it.Next()
	it.Next()
}

func ExampleReverseIterator() {
	tree, _ := New()

	tree.Put([]byte("apple"), []byte("sweet"))
	tree.Put([]byte("banana"), []byte("honey"))
	tree.Put([]byte("cinnamon"), []byte("savoury"))
	tree.Put([]byte("apple2"), []byte("sweet2"))
	tree.Put([]byte("banana2"), []byte("honey2"))
	tree.Put([]byte("cinnamon2"), []byte("savoury2"))

	banana, ok := tree.Get([]byte("banana"))
	if ok {
		fmt.Printf("banana = %s\n", string(banana))
	} else {
		fmt.Println("value for banana not found")
	}

	for it := tree.Iterator(&IteratorOptions{reverse: true}); it.HasNext(); {
		key, value := it.Next()
		fmt.Printf("key = %s, value = %s\n", string(key), string(value))
	}
	// Output:
	// banana = honey
	// key = cinnamon2, value = savoury2
	// key = cinnamon, value = savoury
	// key = banana2, value = honey2
	// key = banana, value = honey
	// key = apple2, value = sweet2
	// key = apple, value = sweet
}

func TestReverseIterator(t *testing.T) {
	tree, _ := New()
	for _, c := range iteratorCases {
		tree.Put([]byte{c.key}, []byte(c.value))
	}

	actual := make([]byte, 0)
	for it := tree.Iterator(&IteratorOptions{reverse: true}); it.HasNext(); {
		key, _ := it.Next()
		actual = append(actual, key...)
	}

	isSorted := sort.SliceIsSorted(actual, func(i, j int) bool {
		return actual[i] > actual[j]
	})
	if !isSorted {
		t.Fatalf("each does not traverse in sorted order, produced result: %s", actual)
	}

	expected := make([]byte, 0)
	for _, c := range iteratorCases {
		expected = append(expected, c.key)
	}
	sort.Slice(expected, func(i, j int) bool {
		return expected[i] > expected[j]
	})

	fmt.Println(expected)
	fmt.Println(actual)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("%s != %s", expected, actual)
	}
}
