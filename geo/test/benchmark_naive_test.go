package geo_test

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/raditzlawliet/kdbush"
	"github.com/raditzlawliet/kdbush/geo"
)

var dst = "output"

func extract() {
	archive, err := zip.OpenReader("cities5000.zip")
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)
		fmt.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			fmt.Println("invalid file path")
			return
		}
		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}
}

// Benchmark BuildIndex func
func BenchmarkGeo(b *testing.B) {
	extract()

	cities := []kdbush.Point{}

	// prepare data benchmark
	file, err := os.Open(filepath.Join(dst, "cities5000.txt"))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	parser := csv.NewReader(file)
	parser.Comma = '\t'

	for {
		record, err := parser.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// lat record[4]
		// lng record[5]
		lat, _ := strconv.ParseFloat(record[4], 64)
		lng, _ := strconv.ParseFloat(record[5], 64)
		cities = append(cities, &geo.MarkerPoint{Lat: lat, Lng: lng})
	}

	var cases = []struct {
		Points []kdbush.Point
		Total  int
	}{
		{Points: []kdbush.Point{}, Total: 1000},
		{Points: []kdbush.Point{}, Total: 10_000},
		{Points: []kdbush.Point{}, Total: 100_000},
		{Points: []kdbush.Point{}, Total: 1_000_000},
	}

	// Naive setup point
	for num := range cases {
		for i := 0; i < cases[num].Total; i++ {
			for len(cases[num].Points) < cases[num].Total {
				cases[num].Points = append(cases[num].Points, cities...)
			}
			// truncate
			cases[num].Points = cases[num].Points[:cases[num].Total]
		}
	}

	// Benchmark query 1000 closest with different total data
	for _, v := range cases {
		bush := kdbush.NewBush().
			BuildIndex(v.Points, kdbush.STANDARD_NODE_SIZE)

		b.Run(fmt.Sprintf("AroundClosest1kWithData_%d", v.Total), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				geo.Around(bush, 112.74996851603348, -7.265850333832262, 1000, -1, nil)
			}
		})
	}

	// Benchmark query 50000 closest with different total data
	for _, v := range cases {
		bush := kdbush.NewBush().
			BuildIndex(v.Points, kdbush.STANDARD_NODE_SIZE)

		b.Run(fmt.Sprintf("AroundClosest5kWithData_%d", v.Total), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				geo.Around(bush, 112.74996851603348, -7.265850333832262, 50000, -1, nil)
			}
		})
	}

	// Benchmark query all data
	for _, v := range cases {
		bush := kdbush.NewBush().
			BuildIndex(v.Points, kdbush.STANDARD_NODE_SIZE)

		b.Run(fmt.Sprintf("AroundWithData_%d", v.Total), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				geo.Around(bush, 112.74996851603348, -7.265850333832262, -1, -1, nil)
			}
		})
	}

	// Benchmark query closest random point
	for _, v := range cases {
		bush := kdbush.NewBush().
			BuildIndex(v.Points, kdbush.STANDARD_NODE_SIZE)

		index := rand.Intn(len(v.Points))

		b.Run(fmt.Sprintf("AroundClosest1RanndomWithData_%d", v.Total), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				geo.Around(bush, v.Points[index].GetX(), v.Points[index].GetY(), 1, -1, nil)
			}
		})
	}

	// Benchmark Distance, get random loc from latest case
	index1 := rand.Intn(len(cases[len(cases)-1].Points))
	point1 := cases[len(cases)-1].Points[index1]
	index2 := rand.Intn(len(cases[len(cases)-1].Points))
	point2 := cases[len(cases)-1].Points[index2]

	b.Run(fmt.Sprintf("Distance"), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			geo.Distance(point1.GetX(), point1.GetY(), point2.GetX(), point2.GetY())
		}
	})

}
