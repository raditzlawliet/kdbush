package kdbush_test

import (
	"math/rand"
	"testing"

	"github.com/raditzlawliet/kdbush"
)

type testCase struct {
	Input  []float64
	Index  []int
	Points []kdbush.Point
}

var (
	points = []kdbush.Point{
		&kdbush.SimplePoint{0.0, 0.0},
		// 1-8
		&kdbush.SimplePoint{2.0, 2.0},
		&kdbush.SimplePoint{1.0, 1.0},
		&kdbush.SimplePoint{-2.0, -2.0},
		&kdbush.SimplePoint{-1.0, -1.0},
		&kdbush.SimplePoint{2.0, -2.0},
		&kdbush.SimplePoint{1.0, -1.0},
		&kdbush.SimplePoint{-2.0, 2.0},
		&kdbush.SimplePoint{-1.0, 1.0},
		// 9-16
		&kdbush.SimplePoint{1.0, 0.0},
		&kdbush.SimplePoint{2.0, 0.0},
		&kdbush.SimplePoint{-1.0, 0.0},
		&kdbush.SimplePoint{-2.0, 0.0},
		&kdbush.SimplePoint{0.0, 1.0},
		&kdbush.SimplePoint{0.0, 2.0},
		&kdbush.SimplePoint{0.0, -1.0},
		&kdbush.SimplePoint{0.0, -2.0},
	}
)

// Test Build Index with more than 1 million point
func TestGeneration(t *testing.T) {
	points2 := []kdbush.Point{}
	for i := 0; i < 1_000_000; i++ {
		points2 = append(points2, &kdbush.SimplePoint{rand.Float64()*24.0 + 24.0, rand.Float64()*24.0 + 24.0})
	}

	bush := kdbush.NewBush().
		BuildIndex([]kdbush.Point(points2), kdbush.STANDARD_NODE_SIZE)

	if len(bush.Points) != len(points2) {
		t.Fatalf("Kdbush.Point size %v not same with points2 size %v", len(bush.Points), len(points2))
	}

	for i := range points2 {
		if points2[i] != bush.Points[i] {
			t.Fatalf("Index %v Kdbush.Point %v not same with points2  %v", i, bush.Points[i], points[i])
		}
	}
}

// Test Range func
func TestRange(t *testing.T) {
	bush := kdbush.NewBush().
		BuildIndex([]kdbush.Point(points), kdbush.STANDARD_NODE_SIZE)

	if len(bush.Points) != len(points) {
		t.Fatalf("Kdbush.Point size %v not same with points size %v", len(bush.Points), len(points))
	}

	testCases := []testCase{
		// straight line on x=0
		{
			[]float64{-2.1, 0.0, 2.1, 0},
			[]int{0, 9, 10, 11, 12},
			[]kdbush.Point{
				points[0],
				points[9],
				points[10],
				points[11],
				points[12],
			},
		},
		// straight line on x=1.0
		{
			[]float64{-2.1, 1.0, 2.1, 1.0},
			[]int{2, 8, 13},
			[]kdbush.Point{
				points[2],
				points[8],
				points[13],
			},
		},
	}

	for i, testCase := range testCases {
		index := bush.Range(testCase.Input[0], testCase.Input[1], testCase.Input[2], testCase.Input[3])
		if len(index) != len(testCase.Index) {
			t.Fatalf("case index %v test index size %v not same with result index size %v", i, len(testCase.Index), len(index))

		}
		for k, v := range index {
			if points[v] != testCase.Points[k] {
				t.Fatalf("case index %v result index %v points %v not same with result points %v by index", i, v, points[v], testCase.Points[k])
			}
		}
	}
}

// Test WIthin func
func TestWithin(t *testing.T) {
	bush := kdbush.NewBush().
		BuildIndex([]kdbush.Point(points), kdbush.STANDARD_NODE_SIZE)

	if len(bush.Points) != len(points) {
		t.Fatalf("Kdbush.Point size %v not same with points size %v", len(bush.Points), len(points))
	}

	testCases := []testCase{
		// test within inner radius
		{
			[]float64{0, 0, 0.999},
			[]int{0},
			[]kdbush.Point{
				points[0],
			},
		},
		// test radius more wide, make sure it's circle (not square)
		{
			[]float64{0, 0, 1},
			[]int{0, 9, 11, 13, 15},
			[]kdbush.Point{
				points[0],
				points[9],
				points[11],
				points[13],
				points[15],
			},
		},
		// test radius not center
		{
			[]float64{0.3, 0.2, 0.8},
			[]int{0, 9},
			[]kdbush.Point{
				points[0],
				points[9],
			},
		},
	}

	for i, testCase := range testCases {
		index := bush.Within(&kdbush.SimplePoint{testCase.Input[0], testCase.Input[1]}, testCase.Input[2])
		if len(index) != len(testCase.Index) {
			t.Fatalf("case index %v test index size %v not same with result index size %v", i, len(testCase.Index), len(index))

		}
		for k, v := range index {
			if points[v] != testCase.Points[k] {
				t.Fatalf("case index %v result index %v points %v not same with result points %v by index", i, v, points[v], testCase.Points[k])
			}
		}
	}
}
