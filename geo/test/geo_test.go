package geo_test

import (
	"testing"

	"github.com/raditzlawliet/kdbush"
	"github.com/raditzlawliet/kdbush/geo"
	"github.com/stretchr/testify/assert"
)

var (
	points = []kdbush.Point{
		// Surabaya
		&geo.MarkerPoint{Lat: -7.265850333832262, Lng: 112.74996851603348}, // Submarine Monument
		&geo.MarkerPoint{Lat: -7.261669467066506, Lng: 112.74641761874226}, // Panglima Besar Djendral Soedirman
		&geo.MarkerPoint{Lat: -7.262083358514896, Lng: 112.74242319159997}, // Statue of Gubernur Suryo
		&geo.MarkerPoint{Lat: -7.266374558931883, Lng: 112.74432953003436}, // Monumen Bambu Runcing
		&geo.MarkerPoint{Lat: -7.271393029465542, Lng: 112.74264315372703}, // Karapan Sapi Statue
		// Surabaya but not near from previous one
		&geo.MarkerPoint{Lat: -7.279393373923168, Lng: 112.7413233810989}, // Monumen Perjuangan Polisi Republik Indonesia
		// Jakarta
		&geo.MarkerPoint{Lat: -6.199482563158932, Lng: 106.84831233134457}, // Tugu Proklamasi
		&geo.MarkerPoint{Lat: -6.173354331560208, Lng: 106.82685482999992}, // Monas
	}
)

func TestDistance(t *testing.T) {
	dist := geo.Distance(points[0].GetX(), points[0].GetY(), points[1].GetX(), points[1].GetY())
	assert.GreaterOrEqual(t, dist, 0.0, "distance should be more than 0km")
	assert.LessOrEqual(t, dist, 1.0, "distance should be less than 1km")

	dist2 := geo.Distance(points[0].GetX(), points[0].GetY(), points[7].GetX(), points[7].GetY())
	assert.GreaterOrEqual(t, dist2, 500.0, "distance should be more than 500km")
	assert.LessOrEqual(t, dist2, 1000.0, "distance should be less than 1000km")
}

func TestSimpleAround(t *testing.T) {
	bush := kdbush.NewBush().
		BuildIndexWith(points, kdbush.STANDARD_NODE_SIZE)

	assert.ElementsMatch(t, bush.Points, points, "they should be have same elements")

	testCases := []struct {
		Name         string
		Input        kdbush.Point
		MaxResult    int
		Distance     float64
		Result       []int
		ResultPoints []kdbush.Point
	}{
		{
			"All Surabaya closest 5 within 10km",
			points[0],
			5,
			10,
			[]int{0, 1, 2, 3, 4},
			[]kdbush.Point{points[0], points[1], points[2], points[3], points[4]},
		},
		{
			"All Jakarta closest 2 within 10km",
			points[7],
			5,
			10,
			[]int{6, 7},
			[]kdbush.Point{points[6], points[7]},
		},
		{
			"All",
			points[6],
			-1,
			-1,
			[]int{0, 1, 2, 3, 4, 5, 6, 7},
			points,
		},
		{
			"Closest 2 all distance",
			points[7],
			2,
			-1,
			[]int{6, 7},
			[]kdbush.Point{points[6], points[7]},
		},
		{
			"All point within distance",
			points[6],
			-1,
			10,
			[]int{6, 7},
			[]kdbush.Point{points[6], points[7]},
		},
	}

	for _, testCase := range testCases {
		results := geo.Around(bush, testCase.Input.GetX(), testCase.Input.GetY(), testCase.MaxResult, testCase.Distance, nil)
		assert.ElementsMatch(t, results, testCase.Result, "[%v] Result element index should same", testCase.Name)

		resultPoints := []kdbush.Point{}
		for _, v := range results {
			resultPoints = append(resultPoints, points[v])
		}
		assert.ElementsMatch(t, resultPoints, testCase.ResultPoints, "[%v] Result element point should same", testCase.Name)
	}
}
