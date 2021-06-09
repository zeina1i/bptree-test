# **bp**tree

[![Build Status](https://travis-ci.com/krasun/bptree.svg?branch=main)](https://travis-ci.com/krasun/bptree)
[![codecov](https://codecov.io/gh/krasun/bptree/branch/main/graph/badge.svg?token=8NU6LR4FQD)](https://codecov.io/gh/krasun/bptree)
[![Go Report Card](https://goreportcard.com/badge/github.com/krasun/bptree)](https://goreportcard.com/report/github.com/krasun/bptree)
[![GoDoc](https://godoc.org/https://godoc.org/github.com/krasun/bptree?status.svg)](https://godoc.org/github.com/krasun/bptree)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fkrasun%2Fbptree.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fkrasun%2Fbptree?ref=badge_shield)

An in-memory [B+ tree](https://en.wikipedia.org/wiki/B%2B_tree) implementation for Go with byte-slice keys and values. 

## Installation 

To install, run:

```
go get github.com/krasun/bptree
```

## Quickstart

Feel free to play: 

```go
package main

import (
	"fmt"

	"github.com/krasun/bptree"
)

func main() {
	tree := bptree.New()

	tree.Put([]byte("apple"), []byte("sweet"))
	tree.Put([]byte("banana"), []byte("honey"))
	tree.Put([]byte("cinnamon"), []byte("savoury"))

	banana, ok := tree.Get([]byte("banana"))
	if ok {
		fmt.Printf("banana = %s\n", string(banana))
	} else {
		fmt.Println("value for banana not found")
	}

	tree.ForEach(func(key, value []byte) {
		fmt.Printf("key = %s, value = %s\n", string(key), string(value))
	})

	// Output: 
	// banana = honey
	// key = apple, value = sweet
	// key = banana, value = honey
	// key = cinnamon, value = savoury
}
```

You can use an iterator: 

```go
package main

import (
	"fmt"

	"github.com/krasun/bptree"
)

func main() {
	tree := bptree.New(bptree.Order(3))

	tree.Put([]byte("apple"), []byte("sweet"))
	tree.Put([]byte("banana"), []byte("honey"))
	tree.Put([]byte("cinnamon"), []byte("savoury"))

	banana, ok := tree.Get([]byte("banana"))
	if ok {
		fmt.Printf("banana = %s\n", string(banana))
	} else {
		fmt.Println("value for banana not found")
	}

	for it := tree.Iterator(); it.HasNext(); {
		key, value := it.Next()
		fmt.Printf("key = %s, value = %s\n", string(key), string(value))
	}

	// Output: 
	// banana = honey
	// key = apple, value = sweet
	// key = banana, value = honey
	// key = cinnamon, value = savoury
}
```

An iterator is stateful. You can have multiple iterators without any impact on each other, but make sure to synchronize access to them and the tree in a concurrent environment.

Caution! `Next` panics if there is no next element. Make sure to test for the next element with `HasNext` before.

## Use cases 

1. When you want to use []byte as a key in the map. 
2. When you want to iterate over keys in map in sorted order.

## Limitations 

**Caution!** To guarantee that the B+ tree properties are not violated, keys are copied. 

You should clearly understand what []byte slice is and why it is dangerous to use it as a key. Go language authors do prohibit using byte slice ([]byte) as a map key for a reason. The point is that you can change the values of the key and thus violate the invariants of map: 

```go
// if it worked 
b := []byte{1}
m := make(map[[]byte]int)
m[b] = 1

b[0] = 2 // it would violate the invariants 
m[[]byte{1}] // what do you expect to receive?
```

So to make sure that this situation does not occur in the tree, the key is copied byte by byte.

## Benchmark

Regular Go map is as twice faster for put and get than B+ tree. But if you 
need to iterate over keys in sorted order, the picture is slightly different: 

```
$ go test -benchmem -bench .
goos: darwin
goarch: amd64
pkg: github.com/krasun/bptree
BenchmarkTreePut-8                     	     177	   6379568 ns/op	 2825120 B/op	   99844 allocs/op
BenchmarkMapPut-8                      	     517	   2537254 ns/op	 1732410 B/op	   20151 allocs/op
BenchmarkTreePutRandomized-8           	     140	  10885305 ns/op	 1621450 B/op	   69407 allocs/op
BenchmarkMapPutRandomized-8            	     453	   2599835 ns/op	  981421 B/op	   20110 allocs/op
BenchmarkMapGet-8                      	    1197	    856502 ns/op	   38880 B/op	    9900 allocs/op
BenchmarkTreeGet-8                     	     465	   3217380 ns/op	   38880 B/op	    9900 allocs/op
BenchmarkTreePutAndForEach-8           	     133	   8896404 ns/op	 2825120 B/op	   99844 allocs/op
BenchmarkMapPutAndIterateAfterSort-8   	     166	   6983623 ns/op	 2559120 B/op	   20175 allocs/op
PASS
ok  	github.com/krasun/bptree	14.489s
```

## Tests

Run tests with: 

```
$ go test -cover .
ok  	github.com/krasun/bptree	0.245s	coverage: 100.0% of statements
```

## License 

**bp**tree is released under [the MIT license](LICENSE).