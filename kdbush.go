// Package kdbush implements kdbush-tree
package kdbush

// STANDARD_NODE_SIZE default nodeSize kdbush-tree. Higher value means faster indexing but slower search and vice versa
const STANDARD_NODE_SIZE = 64

// KDBush an instance
type KDBush struct {
	nodeSize int
	ids      []int
	coords   []float64
	indexed  bool
}

// NewBush return a new pointer of [KDBush]
func NewBush() *KDBush {
	kd := KDBush{}
	return &kd
}

// BuildIndex build kd-tree index given list of Points
func (kd *KDBush) BuildIndex(points []Point, nodeSize int) *KDBush {
	kd.indexed = false
	kd.nodeSize = nodeSize

	kd.ids = make([]int, len(points))
	kd.coords = make([]float64, 2*len(points))

	for i, v := range points {
		kd.ids[i] = i
		kd.coords[i*2] = v.GetX()
		kd.coords[i*2+1] = v.GetY()
	}

	sort(kd.ids, kd.coords, kd.nodeSize, 0, len(kd.ids)-1, 0)

	kd.indexed = true
	return kd
}

// query helper struct for API Range & Within finding result
type query struct {
	left  int
	right int
	axis  int
}

// Range returns all indexes points across [minX], [minY], [maxX], [maxY]
func (kd *KDBush) Range(minX, minY, maxX, maxY float64) []int {
	if !kd.indexed {
		return []int{}
	}

	stack := []query{{0, len(kd.ids) - 1, 0}}
	result := []int{}

	var x, y float64

	for (len(stack)) > 0 {
		left := stack[len(stack)-1].left
		right := stack[len(stack)-1].right
		axis := stack[len(stack)-1].axis
		stack = append(stack[:len(stack)-1], stack[len(stack):]...) // .pop()

		// search linearly
		if right-left <= kd.nodeSize {
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
			stack = append(stack, query{left, m - 1, 1 - axis})
		}
		if (axis == 0 && maxX >= x) || (axis != 0 && maxY >= y) {
			stack = append(stack, query{m + 1, right, 1 - axis})
		}
	}

	return result
}

// Within returns all indexes points within radius of given single [Point]
func (kd *KDBush) Within(qx, qy float64, radius float64) []int {
	if !kd.indexed {
		return []int{}
	}

	stack := []query{{0, len(kd.ids) - 1, 0}}
	result := []int{}

	r2 := radius * radius

	var x, y float64

	for (len(stack)) > 0 {
		left := stack[len(stack)-1].left
		right := stack[len(stack)-1].right
		axis := stack[len(stack)-1].axis
		stack = append(stack[:len(stack)-1], stack[len(stack):]...) // .pop()

		// search linearly
		if right-left <= kd.nodeSize {
			for i := left; i <= right; i++ {
				if sqrtDist(kd.coords[2*i], kd.coords[2*i+1], qx, qy) <= r2 {
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
			stack = append(stack, query{left, m - 1, 1 - axis})
		}
		if (axis == 0 && qx+radius >= x) || (axis != 0 && qy+radius >= y) {
			stack = append(stack, query{m + 1, right, 1 - axis})
		}
	}

	return result
}

//
// Helper get private param
//

// GetNodeSize return current nodesize
func (kd *KDBush) GetNodeSize() int {
	return kd.nodeSize
}

// GetIndexed return all kdtree indexes
func (kd *KDBush) GetIndexes() []int {
	return kd.ids
}

// GetCoords return all coords
func (kd *KDBush) GetCoords() []float64 {
	return kd.coords
}

// Indexed return it's KDBush already indexed or not
func (kd *KDBush) Indexed() bool {
	return kd.indexed
}
