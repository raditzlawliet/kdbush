package kdbush

// Point interface will be used for KDBush
type Point interface {
	GetX() (X float64)
	GetY() (Y float64)
}

type SimplePoint struct {
	X, Y float64
}

func (p *SimplePoint) GetX() (X float64) {
	return p.X
}

func (p *SimplePoint) GetY() (Y float64) {
	return p.Y
}
