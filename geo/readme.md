# Go - KDBush/geo

A geographic extension for Golang port of KDBush, the fastest static spatial index for points.

It implements fast nearest neighbors queries for locations on Earth, taking Earth curvature and date line wrapping into account. Inspired by sphere-knn, but uses a different algorithm.

This extension works for and require [https://github.com/raditzlawliet/kdbush](https://github.com/raditzlawliet/kdbush) and only use standard library

This implementation is based on:

- [Geospatial Ext. GeoKDBush-tk](https://github.com/tkafka/geokdbush-tk)

Dump file cities5000 for Benchmark are from [geonames.org data](https://download.geonames.org/export/dump/).

## Usage

By default, when you go get kdbush, it will include this extension

```go
import(
    "github.com/raditzlawliet/kdbush"
    "github.com/raditzlawliet/kdbush/geo"
)

//...
points = []kdbush.Point{
    // Surabaya
    &geo.MarkerPoint{Lat: -7.265850333832262, Lng: 112.74996851603348},
    &geo.MarkerPoint{Lat: -7.261669467066506, Lng: 112.74641761874226},
    &geo.MarkerPoint{Lat: -7.262083358514896, Lng: 112.74242319159997},
    &geo.MarkerPoint{Lat: -7.266374558931883, Lng: 112.74432953003436},
    &geo.MarkerPoint{Lat: -7.271393029465542, Lng: 112.74264315372703},
    // Surabaya but not near from previous one
    &geo.MarkerPoint{Lat: -7.279393373923168, Lng: 112.7413233810989},
    // Jakarta
    &geo.MarkerPoint{Lat: -6.199482563158932, Lng: 106.84831233134457},
    &geo.MarkerPoint{Lat: -6.173354331560208, Lng: 106.82685482999992},
    ...
}

// Build Index
bush := kdbush.NewBush().
    BuildIndexWith(points, kdbush.STANDARD_NODE_SIZE)

// Try to search 5 cloests around 10km in Jakarta
results := geo.Around(bush, 106.84831233134457,  -6.199482563158932, 5, 10, nil)
fmt.Println(results)
// [6, 7]
```

## API

### Around(kdbush, longitude, latitude, maxResults, maxDistanceInKm, filterFn)

Returns an array of the closest ids (indices) of points from a given location in order of increasing distance.

- `kdbush`: kdbush pointer `*KDBush`
- `longitude`: query point longitude `float64`
- `latitude`: query point latitude `float64`
- `maxResults`: maximum number of points to return (-1 for all result) `int`.
- `maxDistance`: maximum distance in kilometers to search within (-1 for all distance) `float64`.
- `filterFn`: (optional) a function to filter the results (ids) with `func(int) bool`.

### Distance(longitude1, latitude1, longitude2, latitude2)

Returns great circle distance between two locations in kilometers.

- `longitude1`: query point longitude location A `float64`
- `latitude1`: query point latitude location A `float64`
- `longitude2`: query point longitude location B `float64`
- `latitude2`: query point latitude location B `float64`

## Benchmark

All benchmark are run on Go 1.20.3, Windows 11 & 12th Gen Intel(R) Core(TM) i7-12700H (Laptop version).

**Do not trust benchmark**

`go test -bench=BenchmarkGeo -benchmem -benchtime=10s`

Benchmark Around (-benchtime=10s); NodeSize=64

- Around 1k closest with various data
- Around 5k closest with various data
- All data (_not recommended_)
- Around 1 cloests with various data (_Recommended_)

```sh
BenchmarkGeo/AroundClosest1kWithData_1000-20               74640            162609 ns/op          139080 B/op       1054 allocs/op
BenchmarkGeo/AroundClosest1kWithData_10000-20              57082            209555 ns/op          186408 B/op       1547 allocs/op
BenchmarkGeo/AroundClosest1kWithData_100000-20             64539            188799 ns/op          168296 B/op       1428 allocs/op
BenchmarkGeo/AroundClosest1kWithData_1000000-20            66148            186287 ns/op          178664 B/op       1536 allocs/op
BenchmarkGeo/AroundClosest5kWithData_1000-20               73554            369826 ns/op          139080 B/op       1054 allocs/op
BenchmarkGeo/AroundClosest5kWithData_10000-20               6464           1857184 ns/op         1426632 B/op      10544 allocs/op
BenchmarkGeo/AroundClosest5kWithData_100000-20              1185          10202234 ns/op         7698122 B/op      57933 allocs/op
BenchmarkGeo/AroundClosest5kWithData_1000000-20             1341           8911351 ns/op         7179144 B/op      53549 allocs/op
BenchmarkGeo/AroundWithData_1000-20                        75081            161414 ns/op          139080 B/op       1054 allocs/op
BenchmarkGeo/AroundWithData_10000-20                        6279           1880233 ns/op         1426632 B/op      10544 allocs/op
BenchmarkGeo/AroundWithData_100000-20                        607          19836218 ns/op        14280010 B/op     104140 allocs/op
BenchmarkGeo/AroundWithData_1000000-20                        57         206906467 ns/op        141484369 B/op   1032826 allocs/op
BenchmarkGeo/AroundClosest1RanndomWithData_1000-20       2282091              5206 ns/op            9176 B/op         84 allocs/op
BenchmarkGeo/AroundClosest1RanndomWithData_10000-20      2766841              4309 ns/op            7000 B/op         71 allocs/op
BenchmarkGeo/AroundClosest1RanndomWithData_100000-20     2078661              5800 ns/op            9752 B/op         90 allocs/op
BenchmarkGeo/AroundClosest1RanndomWithData_1000000-20            3971299              3070 ns/op            4568 B/op         50 allocs/op
```
