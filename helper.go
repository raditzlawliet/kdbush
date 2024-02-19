package kdbush

import (
	"math"
)

// Number
type Number interface {
	float64 | int
}

// sort
func sort(ids []int, coords []float64, nodeSize, left, right, axis int) {
	if (right - left) <= nodeSize {
		return
	}

	m := (left + right) >> 1

	selection(ids, coords, nodeSize, left, right, axis)

	sort(ids, coords, nodeSize, left, m-1, 1-axis)
	sort(ids, coords, nodeSize, m+1, right, 1-axis)
}

// selection
func selection(ids []int, coords []float64, k, left, right, axis int) {
	for right > left {
		if (right - left) > 600 {
			n := float64(right - left + 1)
			m := float64(k - left + 1)
			z := math.Log(n)
			s := 0.5 * math.Exp(2*z/3)

			sds := 1.0
			if (m - n/2) < 0 {
				sds = -1.0
			}

			sd := 0.5 * math.Sqrt(z*s*(n-s)/n) * sds

			newLeft := max(left, int(math.Floor(float64(k)-m*s/n+sd)))
			newRight := min(right, int(math.Floor(float64(k)+(n-m)*s/n+sd)))

			selection(ids, coords, k, newLeft, newRight, axis)
		}

		t := coords[2*k+axis]
		i := left
		j := right

		swapItem(ids, coords, left, k)
		if coords[2*right+axis] > t {
			swapItem(ids, coords, left, right)
		}

		for i < j {
			swapItem(ids, coords, i, j)
			i += 1
			j -= 1

			for coords[2*i+axis] < t {
				i += 1
			}
			for coords[2*j+axis] > t {
				j -= 1
			}
		}

		if coords[2*left+axis] == t {
			swapItem(ids, coords, left, j)
		} else {
			j += 1
			swapItem(ids, coords, j, right)
		}

		if j <= k {
			left = j + 1
		}
		if k <= j {
			right = j - 1
		}
	}
}

// swapItem
func swapItem(ids []int, coords []float64, i, j int) {
	swap(ids, i, j)
	swap(coords, i*2, j*2)
	swap(coords, i*2+1, j*2+1)
}

// swap
func swap[T Number](arr []T, i, j int) {
	v := arr[i]
	arr[i] = arr[j]
	arr[j] = v
}

// min
func min[T Number](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// max
func max[T Number](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// sqrtDist
func sqrtDist(ax, ay, bx, by float64) float64 {
	dx := ax - bx
	dy := ay - by
	return dx*dx + dy*dy
}
