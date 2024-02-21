package geo

// nullInt simple nullable int
type nullInt struct {
	Int   int
	Valid bool
}

// geoNode for KDBush
type geoNode struct {
	// itemID
	itemID nullInt

	left   int
	right  int
	axis   int
	dist   float64
	minLng float64
	minLat float64
	maxLng float64
	maxLat float64

	// Queue index
	index int
}

// geoNodeQueue a priority queue by distance
// Example of Priority Queue is using heap std, see more at https://pkg.go.dev/container/heap#example__priorityQueue
type geoNodeQueue []*geoNode

func (q geoNodeQueue) Len() int { return len(q) }

func (q geoNodeQueue) Less(i, j int) bool {
	return q[i].dist > q[j].dist
}

func (q geoNodeQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

func (q *geoNodeQueue) Push(x any) {
	n := len(*q)
	item := x.(*geoNode)
	item.index = n
	*q = append(*q, item)
}

func (q *geoNodeQueue) Pop() any {
	old := *q
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*q = old[0 : n-1]
	return item
}
