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

// 1. Build Index (prefered way)
bush := kdbush.NewBush().
    BuildIndexWith(points, kdbush.STANDARD_NODE_SIZE)

// Ready to use
// API Range
indexes := bush.Range(-2.1, 1.0, 2.1, 1.0)
// [2 8 13]
for _, v := range indexes {
  fmt.Println(points[v])
  fmt.Println(bush.Points[v]) // or
}
// &{1 1} &{-1 1} &{0 1}

// API Within Radius
indexes := bush.Within(&kdbush.SimplePoint{0, 0}, 1)
// [0 9 11 13 15]
for _, v := range indexes {
  fmt.Println(points[v])
  fmt.Println(bush.Points[v]) // or
}
// &{0 0} &{1 0} &{-1 0} &{0 1} &{0 -1}

```

You also can use variadic param when adding points. it will use more alloc, but some case may need it.
Since index are static. if you need add new point for some case, you can add Index using Add() and rebuild index using BuildIndex()
_Keep in mind, you will need more resources to rebuild..._

```go
// Build Index using variadic version (use more alloc)
bush := kdbush.NewBush().
    Add(points...).
    BuildIndex(kdbush.STANDARD_NODE_SIZE)

// Other example
bush := kdbush.NewBush().
    BuildIndexWith(points, kdbush.STANDARD_NODE_SIZE)
for _, v := range newPoints {
    Add(v)
}
bush.BuildIndex(kdbush.STANDARD_NODE_SIZE)
```

### Benchmark

All benchmark are run on Go 1.20.3, Windows 11 & 12th Gen Intel(R) Core(TM) i7-12700H (Laptop version).

**Do not trust benchmark**

`go test -bench=BenchmarkBuildIndex -benchmem -benchtime=10s`

Build Index With Benchmark (-benchtime=10s); NodeSize=64

```sh
BenchmarkBuildIndexWith/total_1000-20             775942             15443 ns/op           24576 B/op          2 allocs/op
BenchmarkBuildIndexWith/total_10000-20             22734            522197 ns/op          245760 B/op          2 allocs/op
BenchmarkBuildIndexWith/total_100000-20             1813           6551473 ns/op         2408448 B/op          2 allocs/op
BenchmarkBuildIndexWith/total_1000000-20             146          84477151 ns/op        24010752 B/op          2 allocs/op
```

Build Index (Variadic with Add) Benchmark (-benchtime=10s); NodeSize=64

```sh
BenchmarkBuildIndex/total_1000-20                 654110             17228 ns/op           40960 B/op          3 allocs/op
BenchmarkBuildIndex/total_10000-20                 23164            520225 ns/op          409600 B/op          3 allocs/op
BenchmarkBuildIndex/total_100000-20                 1587           7685670 ns/op         4014080 B/op          3 allocs/op
BenchmarkBuildIndex/total_1000000-20                 134          89170069 ns/op        40017920 B/op          3 allocs/op
```

Range Benchmark; NodeSize=64

```sh
BenchmarkRange/total_1000-20        790533	      1428 ns/op	     120 B/op	       5 allocs/op
BenchmarkRange/total_10000-20       795254	      1424 ns/op	     216 B/op	       9 allocs/op
BenchmarkRange/total_100000-20      647720	      1898 ns/op	     288 B/op	      12 allocs/op
BenchmarkRange/total_1000000-20     525925	      2313 ns/op	     360 B/op	      15 allocs/op
ok  	github.com/raditzlawliet/kdbush	7.733s

```

Within Benchmark; NodeSize=64

```sh
BenchmarkWithin/total_1000-20       725872	      1524 ns/op	     136 B/op	       6 allocs/op
BenchmarkWithin/total_10000-20      771256	      1490 ns/op	     232 B/op	      10 allocs/op
BenchmarkWithin/total_100000-20     619287	      1938 ns/op	     304 B/op	      13 allocs/op
BenchmarkWithin/total_1000000-20    550041	      2350 ns/op	     376 B/op	      16 allocs/op
ok  	github.com/raditzlawliet/kdbush	7.747s
```
