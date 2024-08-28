package kdbush_test

import (
	"math/rand"
	"testing"

	"github.com/raditzlawliet/kdbush"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	Input        []float64
	ResultPoints []kdbush.Point
}

func init() {
	maxPoint := 10
	for i := -maxPoint; i <= maxPoint; i++ {
		for j := -maxPoint; j <= maxPoint; j++ {
			point := &kdbush.SimplePoint{
				X: float64(i),
				Y: float64(j),
			}
			points = append(points, point)
		}

	}
}

var (
	points = []kdbush.Point{}
)

// Test Build Index with more than 1k point
func TestGeneration(t *testing.T) {
	var rng = rand.New(rand.NewSource(1))

	_points := []kdbush.Point{}
	for i := 0; i < 1_000; i++ {
		_points = append(_points, &kdbush.SimplePoint{rng.Float64()*24.0 + 24.0, rng.Float64()*24.0 + 24.0})
	}

	bush := kdbush.NewBush()
	assert.Equal(t, bush.Indexed(), false, "should not indexed")
	assert.Equal(t, bush.Within(0, 0, 0), []int{}, "should return empty slice of int")
	assert.Equal(t, bush.Range(-5, 0, 5, 0), []int{}, "should return empty slice of int")

	bush.BuildIndex(_points, kdbush.STANDARD_NODE_SIZE)

	assert.Equal(t, len(bush.GetIndexes()), len(_points), "they should be have same length")
	assert.Equal(t, bush.GetNodeSize(), kdbush.STANDARD_NODE_SIZE, "nodesize should be same")
	assert.Equal(t, bush.Indexed(), true, "should indexed")
	assert.Equal(t, len(bush.GetIndexes()), len(_points), "indexes length should same with points")
	assert.Equal(t, len(bush.GetCoords()), len(_points)*2, "coords length should 2x with points")
}

// Test Range func
func TestRange(t *testing.T) {
	testCases := []testCase{
		// straight line on x=0
		{
			[]float64{-2.1, 0.0, 2.1, 0},
			[]kdbush.Point{
				&kdbush.SimplePoint{-2, 0},
				&kdbush.SimplePoint{-1, 0},
				&kdbush.SimplePoint{0, 0},
				&kdbush.SimplePoint{1, 0},
				&kdbush.SimplePoint{2, 0},
			},
		},
		// straight on x=1.0 to x=2.0
		{
			[]float64{-2.1, 1.0, 2.1, 2.0},
			[]kdbush.Point{
				&kdbush.SimplePoint{-2, 1},
				&kdbush.SimplePoint{-1, 1},
				&kdbush.SimplePoint{0, 1},
				&kdbush.SimplePoint{1, 1},
				&kdbush.SimplePoint{2, 1},
				&kdbush.SimplePoint{-2, 2},
				&kdbush.SimplePoint{-1, 2},
				&kdbush.SimplePoint{0, 2},
				&kdbush.SimplePoint{1, 2},
				&kdbush.SimplePoint{2, 2},
			},
		},
	}

	// Standard Node Size
	{
		bush := kdbush.NewBush().BuildIndex(points, kdbush.STANDARD_NODE_SIZE)

		for _, testCase := range testCases {
			indexes := bush.Range(testCase.Input[0], testCase.Input[1], testCase.Input[2], testCase.Input[3])
			assert.Equal(t, len(indexes), len(testCase.ResultPoints), "it should be has same count result")
			for _, index := range indexes {
				assert.Contains(t, testCase.ResultPoints, points[index], "it should be has same elements")
			}
		}
	}

	// Small Node Size
	{
		bush := kdbush.NewBush().BuildIndex(points, 4)

		for _, testCase := range testCases {
			indexes := bush.Range(testCase.Input[0], testCase.Input[1], testCase.Input[2], testCase.Input[3])
			assert.Equal(t, len(indexes), len(testCase.ResultPoints), "it should be has same count result")

			for _, index := range indexes {
				assert.Contains(t, testCase.ResultPoints, points[index], "it should be has same elements")
			}
		}
	}
}

// Test Within func
func TestWithin(t *testing.T) {
	testCases := []testCase{
		// test within inner radius
		{
			[]float64{0, 0, 0.999},
			[]kdbush.Point{
				&kdbush.SimplePoint{0, 0},
			},
		},
		// test radius more wide, make sure it's circle (not square)
		{
			[]float64{0, 0, 1},
			[]kdbush.Point{
				&kdbush.SimplePoint{1, 0},
				&kdbush.SimplePoint{0, 0},
				&kdbush.SimplePoint{-1, 0},
				&kdbush.SimplePoint{0, 1},
				&kdbush.SimplePoint{0, -1},
			},
		},
		// test radius not center
		{
			[]float64{0.3, 0.2, 0.8},
			[]kdbush.Point{
				&kdbush.SimplePoint{0, 0},
				&kdbush.SimplePoint{1, 0},
			},
		},
	}

	// Standard Node Size
	{
		bush := kdbush.NewBush().BuildIndex(points, kdbush.STANDARD_NODE_SIZE)

		for _, testCase := range testCases {
			indexes := bush.Within(testCase.Input[0], testCase.Input[1], testCase.Input[2])
			assert.Equal(t, len(indexes), len(testCase.ResultPoints), "it should be has same count result")

			for _, index := range indexes {
				assert.Contains(t, testCase.ResultPoints, points[index], "it should be has same elements")
			}
		}
	}

	// Small Node Size
	{
		bush := kdbush.NewBush().BuildIndex(points, 4)

		for _, testCase := range testCases {
			indexes := bush.Within(testCase.Input[0], testCase.Input[1], testCase.Input[2])
			assert.Equal(t, len(indexes), len(testCase.ResultPoints), "it should be has same count result")

			for _, index := range indexes {
				assert.Contains(t, testCase.ResultPoints, points[index], "it should be has same elements")
			}
		}
	}
}
