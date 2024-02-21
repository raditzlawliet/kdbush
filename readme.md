# Go - KDBush

[![Go Reference](https://pkg.go.dev/badge/github.com/raditzlawliet/kdbush.svg)](https://pkg.go.dev/github.com/raditzlawliet/kdbush)
[![codecov](https://codecov.io/gh/raditzlawliet/kdbush/graph/badge.svg?token=0H3J4MQK59)](https://codecov.io/gh/raditzlawliet/kdbush)

Golang KD-Bush implementation

A very fast static spatial index for 2D points based on a flat KD-tree and almost Zero-Allocation

- 2 Dimensional Points only — no rectangles.
- Static — you can't add/remove items after initial indexing (You need to rebuild index)
- Faster indexing and search, with lower memory footprint
- Build-in API with **Zero-Allocation** (See [#Benchmark](#benchmark))
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

All benchmark are run on Go 1.20.3, Windows 11 & 12th Gen Intel(R) Core(TM) i7-12700H (Laptop version).

**Do not trust benchmark**

`go test -bench=. -benchmem -benchtime=10s`

Benchmark (-benchtime=10s); NodeSize=64

```sh
BenchmarkBuildIndex/nodeSize_64_total_1000-20             566725             20356 ns/op           24576 B/op          2 allocs/op
BenchmarkBuildIndex/nodeSize_64_total_10000-20             13802            876863 ns/op          245760 B/op          2 allocs/op
BenchmarkBuildIndex/nodeSize_64_total_100000-20             1122          10741001 ns/op         2408451 B/op          2 allocs/op
BenchmarkBuildIndex/nodeSize_64_total_1000000-20             100         122507051 ns/op        24010786 B/op          2 allocs/op
BenchmarkBuildIndex/nodeSize_8_total_1000-20              181689             66618 ns/op           24576 B/op          2 allocs/op
BenchmarkBuildIndex/nodeSize_8_total_10000-20              10000           1151113 ns/op          245760 B/op          2 allocs/op
BenchmarkBuildIndex/nodeSize_8_total_100000-20               805          14646166 ns/op         2408448 B/op          2 allocs/op
BenchmarkBuildIndex/nodeSize_8_total_1000000-20               72         161955225 ns/op        24010752 B/op          2 allocs/op
BenchmarkRange/nodeSize_64_total_10-20                  839765822               14.11 ns/op            0 B/op          0 allocs/op
BenchmarkRange/nodeSize_64_total_100-20                 157284577               76.62 ns/op            0 B/op          0 allocs/op
BenchmarkRange/nodeSize_64_total_1000-20                100000000              100.4 ns/op             0 B/op          0 allocs/op
BenchmarkRange/nodeSize_64_total_10000-20               165987278               72.31 ns/op            0 B/op          0 allocs/op
BenchmarkRange/nodeSize_64_total_100000-20              125252017               97.03 ns/op            0 B/op          0 allocs/op
BenchmarkRange/nodeSize_64_total_1000000-20             97950468               125.1 ns/op             0 B/op          0 allocs/op
BenchmarkRange/nodeSize_8_total_10-20                   1000000000               8.681 ns/op           0 B/op          0 allocs/op
BenchmarkRange/nodeSize_8_total_100-20                  696466574               17.35 ns/op            0 B/op          0 allocs/op
BenchmarkRange/nodeSize_8_total_1000-20                 460818502               26.25 ns/op            0 B/op          0 allocs/op
BenchmarkRange/nodeSize_8_total_10000-20                303527074               39.56 ns/op            0 B/op          0 allocs/op
BenchmarkRange/nodeSize_8_total_100000-20               243663084               48.67 ns/op            0 B/op          0 allocs/op
BenchmarkRange/nodeSize_8_total_1000000-20              194835972               60.02 ns/op            0 B/op          0 allocs/op
BenchmarkWithin/nodeSize_64_total_10-20                 984755572               12.07 ns/op            0 B/op          0 allocs/op
BenchmarkWithin/nodeSize_64_total_100-20                176766757               67.20 ns/op            0 B/op          0 allocs/op
BenchmarkWithin/nodeSize_64_total_1000-20               134922337               89.54 ns/op            0 B/op          0 allocs/op
BenchmarkWithin/nodeSize_64_total_10000-20              192972932               62.07 ns/op            0 B/op          0 allocs/op
BenchmarkWithin/nodeSize_64_total_100000-20             151349331               79.71 ns/op            0 B/op          0 allocs/op
BenchmarkWithin/nodeSize_64_total_1000000-20            100000000              101.4 ns/op             0 B/op          0 allocs/op
BenchmarkWithin/nodeSize_8_total_10-20                  1000000000               8.625 ns/op           0 B/op          0 allocs/op
BenchmarkWithin/nodeSize_8_total_100-20                 612873906               19.54 ns/op            0 B/op          0 allocs/op
BenchmarkWithin/nodeSize_8_total_1000-20                425276205               27.61 ns/op            0 B/op          0 allocs/op
BenchmarkWithin/nodeSize_8_total_10000-20               302259685               40.02 ns/op            0 B/op          0 allocs/op
BenchmarkWithin/nodeSize_8_total_100000-20              100000000              136.8 ns/op             0 B/op          0 allocs/op
BenchmarkWithin/nodeSize_8_total_1000000-20             203591168               56.46 ns/op            0 B/op          0 allocs/op
```
