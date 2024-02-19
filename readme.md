## Go - KDBush

Golang KD-Bush implementation

A very fast static spatial index for 2D points based on a flat KD-tree

- 2 Dimensional Points only — no rectangles.
- Static — you can't add/remove items after initial indexing (You need to rebuild index)
- Faster indexing and search, with lower memory footprint
- Build-in feat
  - Range: return indexes within 2 point
  - Within: return indexes within radius of point

This implementation is based on:

- [Javascript - KDBush](https://github.com/mourner/kdbush)

Requirement:

- Go 1.18+ (Generic)

### Usage

Install like usually

```sh
go get github.com/raditzlawliet/kdbush
```

Before building index, you need to create list of Points by implement Point interface

```go
// Point interface
type Point interface {
	GetX() (X float64)
	GetY() (Y float64)
}

// Example Simple Point implement Point Interface
type SimplePoint struct {
	X, Y float64
}

func (p *SimplePoint) GetX() (X float64) {
	return p.X
}

func (p *SimplePoint) GetY() (Y float64) {
	return p.Y
}
```

Then, you can building KDBush index.
Standard Node Size is 64. Higher value means faster indexing but slower search and vice versa.

```go
// 0. Creating array Points
points = []kdbush.Point{
    &kdbush.SimplePoint{0.0, 0.0},
    &kdbush.SimplePoint{2.0, 2.0},
    &kdbush.SimplePoint{1.0, 1.0},
    &kdbush.SimplePoint{-2.0, -2.0},
    &kdbush.SimplePoint{-1.0, -1.0},
    &kdbush.SimplePoint{2.0, -2.0},
    &kdbush.SimplePoint{1.0, -1.0},
    &kdbush.SimplePoint{-2.0, 2.0},
    &kdbush.SimplePoint{-1.0, 1.0},
    &kdbush.SimplePoint{1.0, 0.0},
    &kdbush.SimplePoint{2.0, 0.0},
    &kdbush.SimplePoint{-1.0, 0.0},
    &kdbush.SimplePoint{-2.0, 0.0},
    &kdbush.SimplePoint{0.0, 1.0},
    &kdbush.SimplePoint{0.0, 2.0},
    &kdbush.SimplePoint{0.0, -1.0},
    &kdbush.SimplePoint{0.0, -2.0},

// 1. Build Index
bush := kdbush.NewBush().
    BuildIndex([]kdbush.Point(points), kdbush.STANDARD_NODE_SIZE)

// Ready to use
// API Range
indexes := bush.Range(-2.1, 1.0, 2.1, 1.0)
// [2, 8, 13]

// API Within
indexes := bush.Within(&kdbush.SimplePoint{0, 0}, 1)
// [0, 9, 11, 13, 15]
```

### Benchmark

All benchmark are run on Go 1.20.3, Windows 11 & 12th Gen Intel(R) Core(TM) i7-12700H (Laptop version).

**Do not trust benchmark**

Build Index Benchmark (-benchtime=10s)

```sh
BenchmarkBuildIndex/total_1000-20                 795007             17781 ns/op           24576 B/op          2 allocs/op
BenchmarkBuildIndex/total_10000-20                 22180            540004 ns/op          245760 B/op          2 allocs/op
BenchmarkBuildIndex/total_100000-20                 1635           7414973 ns/op         2408448 B/op          2 allocs/op
BenchmarkBuildIndex/total_1000000-20                 141          85518609 ns/op        24010752 B/op          2 allocs/op
ok      github.com/raditzlawliet/kdbush 75.859s
```

Range Benchmark

```sh
BenchmarkRange/total_1000-20        790533	      1428 ns/op	     120 B/op	       5 allocs/op
BenchmarkRange/total_10000-20       795254	      1424 ns/op	     216 B/op	       9 allocs/op
BenchmarkRange/total_100000-20      647720	      1898 ns/op	     288 B/op	      12 allocs/op
BenchmarkRange/total_1000000-20     525925	      2313 ns/op	     360 B/op	      15 allocs/op
ok  	github.com/raditzlawliet/kdbush	7.733s

```

Within Benchmark

```sh
BenchmarkWithin/total_1000-20       725872	      1524 ns/op	     136 B/op	       6 allocs/op
BenchmarkWithin/total_10000-20      771256	      1490 ns/op	     232 B/op	      10 allocs/op
BenchmarkWithin/total_100000-20     619287	      1938 ns/op	     304 B/op	      13 allocs/op
BenchmarkWithin/total_1000000-20    550041	      2350 ns/op	     376 B/op	      16 allocs/op
ok  	github.com/raditzlawliet/kdbush	7.747s
```
