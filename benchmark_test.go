package kdbush_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/raditzlawliet/kdbush"
)

// Benchmark BuildIndex func
func BenchmarkBuildIndex(b *testing.B) {
	var points = []struct {
		Points   []kdbush.Point
		Total    int
		NodeSize int
	}{
		{Points: []kdbush.Point{}, NodeSize: kdbush.STANDARD_NODE_SIZE, Total: 1000},
		{Points: []kdbush.Point{}, NodeSize: kdbush.STANDARD_NODE_SIZE, Total: 10_000},
		{Points: []kdbush.Point{}, NodeSize: kdbush.STANDARD_NODE_SIZE, Total: 100_000},
		{Points: []kdbush.Point{}, NodeSize: kdbush.STANDARD_NODE_SIZE, Total: 1_000_000},
		{Points: []kdbush.Point{}, NodeSize: 8, Total: 1000},
		{Points: []kdbush.Point{}, NodeSize: 8, Total: 10_000},
		{Points: []kdbush.Point{}, NodeSize: 8, Total: 100_000},
		{Points: []kdbush.Point{}, NodeSize: 8, Total: 1_000_000},
	}

	// Setup
	for num := range points {
		for i := 0; i < points[num].Total; i++ {
			points[num].Points = append(points[num].Points, &kdbush.SimplePoint{rand.Float64()*24.0 + 24.0, rand.Float64()*24.0 + 24.0})
		}

	}

	for _, v := range points {
		b.Run(fmt.Sprintf("nodeSize_%d_total_%d", v.NodeSize, v.Total), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				kdbush.NewBush().
					BuildIndex(v.Points, v.NodeSize)
			}
		})
	}
}

// Benchmark Range func
func BenchmarkRange(b *testing.B) {
	var points = []struct {
		Points   []kdbush.Point
		Total    int
		KDBush   *kdbush.KDBush
		NodeSize int
	}{
		{Points: []kdbush.Point{}, NodeSize: kdbush.STANDARD_NODE_SIZE, Total: 10},
		{Points: []kdbush.Point{}, NodeSize: kdbush.STANDARD_NODE_SIZE, Total: 100},
		{Points: []kdbush.Point{}, NodeSize: kdbush.STANDARD_NODE_SIZE, Total: 1_000},
		{Points: []kdbush.Point{}, NodeSize: kdbush.STANDARD_NODE_SIZE, Total: 10_000},
		{Points: []kdbush.Point{}, NodeSize: kdbush.STANDARD_NODE_SIZE, Total: 100_000},
		{Points: []kdbush.Point{}, NodeSize: kdbush.STANDARD_NODE_SIZE, Total: 1_000_000},
		{Points: []kdbush.Point{}, NodeSize: 8, Total: 10},
		{Points: []kdbush.Point{}, NodeSize: 8, Total: 100},
		{Points: []kdbush.Point{}, NodeSize: 8, Total: 1_000},
		{Points: []kdbush.Point{}, NodeSize: 8, Total: 10_000},
		{Points: []kdbush.Point{}, NodeSize: 8, Total: 100_000},
		{Points: []kdbush.Point{}, NodeSize: 8, Total: 1_000_000},
	}

	// Setup
	for num := range points {
		for i := 0; i < points[num].Total; i++ {
			points[num].Points = append(points[num].Points, &kdbush.SimplePoint{rand.Float64()*24.0 + 24.0, rand.Float64()*24.0 + 24.0})
		}
		points[num].KDBush = kdbush.NewBush().
			BuildIndex(points[num].Points, points[num].NodeSize)
	}

	for _, v := range points {
		b.Run(fmt.Sprintf("nodeSize_%d_total_%d", v.NodeSize, v.Total), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				v.KDBush.Range(10.0, 10, -10, -10)
			}
		})
	}
}

// Benchmark Within func
func BenchmarkWithin(b *testing.B) {
	var points = []struct {
		Points   []kdbush.Point
		Total    int
		KDBush   *kdbush.KDBush
		NodeSize int
	}{
		{Points: []kdbush.Point{}, NodeSize: kdbush.STANDARD_NODE_SIZE, Total: 10},
		{Points: []kdbush.Point{}, NodeSize: kdbush.STANDARD_NODE_SIZE, Total: 100},
		{Points: []kdbush.Point{}, NodeSize: kdbush.STANDARD_NODE_SIZE, Total: 1_000},
		{Points: []kdbush.Point{}, NodeSize: kdbush.STANDARD_NODE_SIZE, Total: 10_000},
		{Points: []kdbush.Point{}, NodeSize: kdbush.STANDARD_NODE_SIZE, Total: 100_000},
		{Points: []kdbush.Point{}, NodeSize: kdbush.STANDARD_NODE_SIZE, Total: 1_000_000},
		{Points: []kdbush.Point{}, NodeSize: 8, Total: 10},
		{Points: []kdbush.Point{}, NodeSize: 8, Total: 100},
		{Points: []kdbush.Point{}, NodeSize: 8, Total: 1_000},
		{Points: []kdbush.Point{}, NodeSize: 8, Total: 10_000},
		{Points: []kdbush.Point{}, NodeSize: 8, Total: 100_000},
		{Points: []kdbush.Point{}, NodeSize: 8, Total: 1_000_000},
	}

	// Setup
	for num := range points {
		for i := 0; i < points[num].Total; i++ {
			points[num].Points = append(points[num].Points, &kdbush.SimplePoint{rand.Float64()*24.0 + 24.0, rand.Float64()*24.0 + 24.0})
		}
		points[num].KDBush = kdbush.NewBush().
			BuildIndex(points[num].Points, points[num].NodeSize)
	}

	for _, v := range points {
		b.Run(fmt.Sprintf("nodeSize_%d_total_%d", v.NodeSize, v.Total), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				v.KDBush.Within(10, 10, 10)
			}
		})
	}
}
