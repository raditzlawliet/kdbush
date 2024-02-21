package geo

// MarkerPoint a basic point struct implement [KDBush.Point]
type MarkerPoint struct {
	Lat, Lng float64
}

func (p *MarkerPoint) GetX() (X float64) {
	return p.Lng
}

func (p *MarkerPoint) GetY() (Y float64) {
	return p.Lat
}
