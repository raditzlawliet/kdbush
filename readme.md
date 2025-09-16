# Go - KDBush

[![Go Reference](https://pkg.go.dev/badge/github.com/raditzlawliet/kdbush.svg)](https://pkg.go.dev/github.com/raditzlawliet/kdbush)
[![codecov](https://codecov.io/gh/raditzlawliet/kdbush/graph/badge.svg?token=0H3J4MQK59)](https://codecov.io/gh/raditzlawliet/kdbush)

Golang KD-Bush implementation

A very fast static spatial index for 2D points based on a flat KD-tree and almost Zero-Allocation

- 2 Dimensional Points only — no rectangles.
- Static — you can't add/remove items after initial indexing (You need to rebuild index)
- Faster indexing and search, with lower memory footprint
- Build-in API with almost **Zero-Allocation** (See [#Benchmark](#benchmark))
  - Range: return indexes within 2 point
  - Within: return indexes within radius of point

Extension

- [Geo Ext.](geo) A simple geographic extension for Golang port of KDBush, support get point around location coordinates

This implementation is based on:

- [Javascript - KDBush](https://github.com/mourner/kdbush)

Requirement:

- Go 1.18+ (Generic)

## Usage

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
    BuildIndex(points, kdbush.STANDARD_NODE_SIZE)

// Ready to use
// API Range
indexes := bush.Range(-2.1, 1.0, 2.1, 1.0)
// [2 8 13]
for _, v := range indexes {
  fmt.Println(points[v])
}
// &{1 1} &{-1 1} &{0 1}

// API Within X, Y with Radius
indexes := bush.Within(0, 0, 1)
// [0 9 11 13 15]
for _, v := range indexes {
  fmt.Println(points[v])
}
// &{0 0} &{1 0} &{-1 0} &{0 1} &{0 -1}

```

kDBush library didn't store original slices to avoid duplication slice

Since index are static. if you need rebuild or adding new point for some case, you can call BuildIndex() multiple times.
_Keep in mind, you will need more resources to rebuild..._

```go
bush := kdbush.NewBush().
  BuildIndex(points, kdbush.STANDARD_NODE_SIZE)
points.Add(newPoint)
bush.BuildIndex(points, kdbush.STANDARD_NODE_SIZE)
```

To avoid datarace on concurrency, please make sure lock bush when build index and you also can check it's already indexed or not via KDBush.Indexed()

## API

### BuildIndex(points, nodeSize) \*KDBush

bulding kd-tree index given list of `points` and `nodeSize`, and return `*KDBush`

- `points`: list Point interface `[]Point`
- `nodeSize`: kd-tree node size. Standard Node Size is 64. Higher value means faster indexing but slower search and vice versa. `int`

### Range(minX, minY, maxX, maxY) []int

return all indexes points across 2 point `minX`, `minY`, `maxX`, `maxY`

- `minX`: point X of A `float64`
- `minY`: point Y of A `float64`
- `maxX`: point X of B `float64`
- `maxY`: point Y of B `float64`

### Within(x, y, radius) []int

return all indexes points within `radius` of given single point `x`, `y`

- `x`: X point `float64`
- `y`: Y point `float64`
- `radius`: radius to search within `float64`

## Benchmark

All benchmark are run on Go 1.20.3, Windows 11 & 12th Gen Intel(R) Core(TM) i7-12700H (Laptop version). **Do not trust benchmark**

`go test -bench=. -benchmem`

```sh
BenchmarkBuildIndex/nodeSize_64_total_1000-20              66031             17605 ns/op           24576 B/op          2 allocs/op
BenchmarkBuildIndex/nodeSize_64_total_10000-20              1380            825414 ns/op          245760 B/op          2 allocs/op
BenchmarkBuildIndex/nodeSize_64_total_100000-20              100          10672552 ns/op         2408448 B/op          2 allocs/op
BenchmarkBuildIndex/nodeSize_64_total_1000000-20               8         125372525 ns/op        24010752 B/op          2 allocs/op
BenchmarkBuildIndex/nodeSize_8_total_1000-20               19930             58559 ns/op           24576 B/op          2 allocs/op
BenchmarkBuildIndex/nodeSize_8_total_10000-20               1090           1097253 ns/op          245760 B/op          2 allocs/op
BenchmarkBuildIndex/nodeSize_8_total_100000-20                81          14333335 ns/op         2408448 B/op          2 allocs/op
BenchmarkBuildIndex/nodeSize_8_total_1000000-20                7         164077629 ns/op        24010752 B/op          2 allocs/op
BenchmarkRange/nodeSize_64_total_10-20                  88866507                14.03 ns/op            0 B/op          0 allocs/op
BenchmarkRange/nodeSize_64_total_100-20                 16080229                76.07 ns/op            0 B/op          0 allocs/op
BenchmarkRange/nodeSize_64_total_1000-20                12128733                99.23 ns/op            0 B/op          0 allocs/op
BenchmarkRange/nodeSize_64_total_10000-20               17081582                69.91 ns/op            0 B/op          0 allocs/op
BenchmarkRange/nodeSize_64_total_100000-20              12586558                92.07 ns/op            0 B/op          0 allocs/op
BenchmarkRange/nodeSize_64_total_1000000-20             10143787               119.4 ns/op             0 B/op          0 allocs/op
BenchmarkRange/nodeSize_8_total_10-20                   140908826                8.477 ns/op           0 B/op          0 allocs/op
BenchmarkRange/nodeSize_8_total_100-20                  70526006                17.62 ns/op            0 B/op          0 allocs/op
BenchmarkRange/nodeSize_8_total_1000-20                 47799242                26.17 ns/op            0 B/op          0 allocs/op
BenchmarkRange/nodeSize_8_total_10000-20                32195490                38.34 ns/op            0 B/op          0 allocs/op
BenchmarkRange/nodeSize_8_total_100000-20               24681555                48.29 ns/op            0 B/op          0 allocs/op
BenchmarkRange/nodeSize_8_total_1000000-20              20815878                56.69 ns/op            0 B/op          0 allocs/op
BenchmarkWithin/nodeSize_64_total_10-20                 97434231                11.75 ns/op            0 B/op          0 allocs/op
BenchmarkWithin/nodeSize_64_total_100-20                17370007                66.08 ns/op            0 B/op          0 allocs/op
BenchmarkWithin/nodeSize_64_total_1000-20               14064499                90.85 ns/op            0 B/op          0 allocs/op
BenchmarkWithin/nodeSize_64_total_10000-20              19023159                61.59 ns/op            0 B/op          0 allocs/op
BenchmarkWithin/nodeSize_64_total_100000-20             14507365                81.15 ns/op            0 B/op          0 allocs/op
BenchmarkWithin/nodeSize_64_total_1000000-20            11783131               101.1 ns/op             0 B/op          0 allocs/op
BenchmarkWithin/nodeSize_8_total_10-20                  134839747                8.773 ns/op           0 B/op          0 allocs/op
BenchmarkWithin/nodeSize_8_total_100-20                 62120803                19.70 ns/op            0 B/op          0 allocs/op
BenchmarkWithin/nodeSize_8_total_1000-20                43613366                27.79 ns/op            0 B/op          0 allocs/op
BenchmarkWithin/nodeSize_8_total_10000-20               29173159                40.47 ns/op            0 B/op          0 allocs/op
BenchmarkWithin/nodeSize_8_total_100000-20              23700660                50.03 ns/op            0 B/op          0 allocs/op
BenchmarkWithin/nodeSize_8_total_1000000-20             21675673                55.10 ns/op            0 B/op          0 allocs/op
```
