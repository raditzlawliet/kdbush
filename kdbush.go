package kdbush

const STANDARD_NODE_SIZE = 64

// KDBush
type KDBush struct {
	Points []Point // pointer

	nodeSize int
	ids      []int
	coords   []float64
}

// NewBush
func NewBush() *KDBush {
	kd := KDBush{
		Points: []Point{},
	}
	return &kd
}

// Add
func (kd *KDBush) Add(points ...Point) *KDBush {
	kd.Points = append(kd.Points, points...)
	return kd
}

// BuildIndexWith
func (kd *KDBush) BuildIndexWith(points []Point, nodeSize int) *KDBush {
	kd.Points = points
	return kd.BuildIndex(nodeSize)
}

// BuildIndex
func (kd *KDBush) BuildIndex(nodeSize int) *KDBush {
	kd.nodeSize = nodeSize

	kd.ids = make([]int, len(kd.Points))
	kd.coords = make([]float64, 2*len(kd.Points))

	for i, v := range kd.Points {
		kd.ids[i] = i
		kd.coords[i*2] = v.GetX()
		kd.coords[i*2+1] = v.GetY()
	}

	sort(kd.ids, kd.coords, kd.nodeSize, 0, len(kd.ids)-1, 0)

	return kd
}

// Range
func (kd *KDBush) Range(minX, minY, maxX, maxY float64) []int {
	stack := [][]int{{0, len(kd.ids) - 1, 0}}
	result := []int{}

	var x, y float64

	for (len(stack)) > 0 {
		axis := stack[len(stack)-1][0]
		right := stack[len(stack)-1][1]
		left := stack[len(stack)-1][2]
		stack = append(stack[:len(stack)-1], stack[len(stack):]...) // .pop()

		// search linearly
		if (right - left) <= kd.nodeSize {
			for i := left; i <= right; i++ {
				x = kd.coords[2*i]
				y = kd.coords[2*i+1]
				if x >= minX && x <= maxX && y >= minY && y <= maxY {
					result = append(result, kd.ids[i])
				}
			}
			continue
		}

		// find in the middle index
		m := (left + right) >> 1

		// include middle item within range
		x = kd.coords[2*m]
		y = kd.coords[2*m+1]
		if x >= minX && x <= maxX && y >= minY && y <= maxY {
			result = append(result, kd.ids[m])
		}

		// queue search in halves that intersect the query
		if (axis == 0 && minX <= x) || (axis != 0 && minY <= y) {
			stack = append(stack, []int{left, m - 1, 1 - axis})
		}
		if (axis == 0 && maxX >= x) || (axis != 0 && maxY >= y) {
			stack = append(stack, []int{m + 1, right, 1 - axis})
		}
	}

	return result
}

// Within
func (kd *KDBush) Within(point Point, radius float64) []int {
	stack := [][]int{{0, len(kd.ids) - 1, 0}}
	result := []int{}

	// r2 := radius * 2
	r2 := radius * radius

	qx, qy := point.GetX(), point.GetY()
	var x, y float64

	for (len(stack)) > 0 {
		axis := stack[len(stack)-1][0]
		right := stack[len(stack)-1][1]
		left := stack[len(stack)-1][2]
		stack = append(stack[:len(stack)-1], stack[len(stack):]...) // .pop()

		// search linearly
		if right-left <= kd.nodeSize {
			for i := left; i <= right; i++ {
				a := sqrtDist(kd.coords[2*i], kd.coords[2*i+1], qx, qy)
				if a <= r2 {
					result = append(result, kd.ids[i])
				}
			}
			continue
		}

		// find in the middle index
		m := (left + right) >> 1

		// include the middle item within range
		x = kd.coords[2*m]
		y = kd.coords[2*m+1]
		if sqrtDist(x, y, qx, qy) <= r2 {
			result = append(result, kd.ids[m])
		}

		// queue search in halves that intersect the query
		if (axis == 0 && qx-radius <= x) || (axis != 0 && qy-radius <= y) {
			stack = append(stack, []int{left, m - 1, 1 - axis})
		}
		if (axis == 0 && qx+radius >= x) || (axis != 0 && qy+radius >= y) {
			stack = append(stack, []int{m + 1, right, 1 - axis})
		}
	}

	return result
}
