package kdbush_test

import (
	"math/rand"
	"testing"

	"github.com/raditzlawliet/kdbush"
	"github.com/stretchr/testify/assert"
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

// Test Build Index with more than 1k point
func TestGeneration(t *testing.T) {
	points2 := []kdbush.Point{}
	for i := 0; i < 1_000; i++ {
		points2 = append(points2, &kdbush.SimplePoint{rand.Float64()*24.0 + 24.0, rand.Float64()*24.0 + 24.0})
	}

	bush := kdbush.NewBush()
	assert.Equal(t, bush.Indexed(), false, "should not indexed")
	assert.Equal(t, bush.Within(0, 0, 0), []int{}, "should return empty slice of int")
	assert.Equal(t, bush.Range(-5, 0, 5, 0), []int{}, "should return empty slice of int")

	bush.BuildIndex(points2, kdbush.STANDARD_NODE_SIZE)

	assert.Equal(t, len(bush.GetIndexes()), len(points2), "they should be have same length")
	assert.Equal(t, bush.GetNodeSize(), kdbush.STANDARD_NODE_SIZE, "nodesize should be same")
	assert.Equal(t, bush.Indexed(), true, "should indexed")
	assert.Equal(t, len(bush.GetIndexes()), len(points2), "indexes length should same with points")
	assert.Equal(t, len(bush.GetCoords()), len(points2)*2, "coords length should 2x with points")
}

// Test Range func
func TestRange(t *testing.T) {
	bush := kdbush.NewBush().
		BuildIndex(points, kdbush.STANDARD_NODE_SIZE)

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

	for _, testCase := range testCases {
		indexes := bush.Range(testCase.Input[0], testCase.Input[1], testCase.Input[2], testCase.Input[3])
		assert.ElementsMatch(t, indexes, testCase.Index, "they should be have same elements")
	}
}

// Test Range func
func TestRangeLowerNodeSize(t *testing.T) {
	bush := kdbush.NewBush().
		BuildIndex(points, 4)

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

	for _, testCase := range testCases {
		indexes := bush.Range(testCase.Input[0], testCase.Input[1], testCase.Input[2], testCase.Input[3])
		assert.ElementsMatch(t, indexes, testCase.Index, "they should be have same elements")
	}
}

// Test Within func
func TestWithin(t *testing.T) {
	bush := kdbush.NewBush().
		BuildIndex(points, kdbush.STANDARD_NODE_SIZE)

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

	for _, testCase := range testCases {
		indexes := bush.Within(testCase.Input[0], testCase.Input[1], testCase.Input[2])
		assert.ElementsMatch(t, indexes, testCase.Index, "they should be have same elements")
	}
}

// Test Within func with Lower Node Size
func TestWithinLowNodeSize(t *testing.T) {
	bush := kdbush.NewBush().
		BuildIndex(points, 4)

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

	for _, testCase := range testCases {
		indexes := bush.Within(testCase.Input[0], testCase.Input[1], testCase.Input[2])
		assert.ElementsMatch(t, indexes, testCase.Index, "they should be have same elements")
	}
}
