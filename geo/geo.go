package geo

import (
	"container/heap"
	"math"

	"github.com/raditzlawliet/kdbush"
)

// Global const for geo calculation
const (
	earthRadius = 6371
	rad         = math.Pi / 180
)

func Around(index kdbush.KDBush, lng, lat float64, maxResults int, maxDistance float64, predicate func(int) bool) []int {
	maxHaverSinDist := 1.0
	result := []int{}

	if maxDistance >= 0 {
		maxHaverSinDist = haverSin(maxDistance / earthRadius)
		_ = maxHaverSinDist
	}

	// a distance-sorted priority queue that will contain both points and kd-tree q
	q := geoNodeQueue{}
	heap.Init(&q)

	// an object that represents the top kd-tree node (the whole Earth)
	node := &geoNode{
		left:   0,                           // left index in the kd-tree array
		right:  len(index.GetIndexes()) - 1, // right index
		axis:   0,                           // 0 for longitude axis and 1 for latitude axis
		dist:   0,                           // will hold the lower bound of children's distances to the query point
		minLng: -180,                        // bounding box of the node
		minLat: -90,
		maxLng: 180,
		maxLat: 90,
	}

	cosLat := math.Cos(lat * rad)

	for node != nil {
		right := node.right
		left := node.left

		if right-left <= index.GetNodeSize() {
			// leaf node

			// add all points of the leaf node to the queue
			for i := left; i <= right; i++ {
				itemId := index.GetIndexes()[i]
				if predicate == nil || (predicate != nil && predicate(itemId)) {
					heap.Push(&q, &geoNode{
						itemID: nullInt{itemId, true},
						dist:   haverSinDist(lng, lat, index.GetCoords()[2*i], index.GetCoords()[2*i+1], cosLat),
					})
				}
			}
		} else {
			// not a leaf node (has child nodes)

			mid := (left + right) >> 1 // middle index
			midLng := index.GetCoords()[2*mid]
			midLat := index.GetCoords()[2*mid+1]

			// add middle point to the queue
			itemId := index.GetIndexes()[mid]
			if predicate == nil || (predicate != nil && predicate(itemId)) {
				heap.Push(&q, &geoNode{
					itemID: nullInt{itemId, true},
					dist:   haverSinDist(lng, lat, midLng, midLat, cosLat),
				})
			}

			nextAxis := (node.axis + 1) % 2

			// first half of the node
			leftNode := &geoNode{
				left:   left,
				right:  mid - 1,
				axis:   nextAxis,
				minLng: node.minLng,
				minLat: node.minLat,
				maxLng: midLng,
				maxLat: midLat,
				dist:   0,
			}

			// second half of the node
			rightNode := &geoNode{
				left:   mid + 1,
				right:  right,
				axis:   nextAxis,
				minLng: midLng,
				minLat: midLat,
				maxLng: node.maxLng,
				maxLat: node.maxLat,
				dist:   0,
			}

			if node.axis == 0 {
				leftNode.maxLat = node.maxLat
				rightNode.minLat = node.minLat
			}
			if node.axis == 1 {
				leftNode.maxLng = node.maxLng
				rightNode.minLng = node.minLng
			}

			leftNode.dist = boxDist(lng, lat, cosLat, leftNode)
			rightNode.dist = boxDist(lng, lat, cosLat, rightNode)

			// add child nodes to the queue
			heap.Push(&q, leftNode)
			heap.Push(&q, rightNode)
		}

		// fetch closest points from the queue; they're guaranteed to be closer than all remaining points (both individual and those in kd-tree nodes), since each node's distance is a lower bound of distances to its children
		for len(q) > 0 && q[0].itemID.Valid {
			candidate := heap.Pop(&q).(*geoNode)
			if candidate.dist > maxHaverSinDist {
				return result
			}

			result = append(result, candidate.itemID.Int)
			if len(result) == maxResults {
				return result
			}
		}

		// the next closest kd-tree node
		if n := heap.Pop(&q); n != nil {
			node = n.(*geoNode)
		} else {
			node = nil
		}
	}

	return result
}
