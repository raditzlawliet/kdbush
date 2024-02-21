package geo

import "math"

// helper for geo extension

// boxDist calculate lower bound for distance from a location to points inside a bounding box
func boxDist(lng, lat, cosLat float64, node *geoNode) float64 {
	// query point is between minimum and maximum longitudes
	if lng >= node.minLng && lng <= node.maxLng {
		if lat < node.minLat {
			return haverSin((lat - node.minLat) * rad)
		} else if lat > node.maxLat {
			return haverSin((lat - node.maxLat) * rad)
		}
		return 0
	}

	// query point is west or east of the bounding box;
	// calculate the extremum for great circle distance from query point to the closest longitude;
	haverSinDLng := math.Min(haverSin((lng-node.minLng)*rad), haverSin((lng-node.maxLng)*rad))
	extremumLat := vertexLat(lat, haverSinDLng)

	// if extremum is inside the box, return the distance to it
	if extremumLat > node.minLat && extremumLat < node.maxLat {
		return haverSinDistPartial(haverSinDLng, cosLat, lat, extremumLat)
	}

	// otherwise return the distan e to one of the bbox corners (whichever is closest)
	return math.Min(haverSinDistPartial(haverSinDLng, cosLat, lat, node.minLat), haverSinDistPartial(haverSinDLng, cosLat, lat, node.maxLat))
}

// haverSin
func haverSin(theta float64) float64 {
	s := math.Sin(theta / 2)
	return s * s
}

// haverSinDistPartial
func haverSinDistPartial(haverSinDLng, cosLat1, lat1, lat2 float64) float64 {
	return cosLat1*math.Cos(lat2*rad)*haverSinDLng + haverSin((lat1-lat2)*rad)
}

// haverSinDist
func haverSinDist(lng1, lat1, lng2, lat2, cosLat1 float64) float64 {
	haverSinDLng := haverSin((lng1 - lng2) * rad)
	return haverSinDistPartial(haverSinDLng, cosLat1, lat1, lat2)
}

// vertexLat
func vertexLat(lat, haverSinDLng float64) float64 {
	cosDLng := 1 - 2*haverSinDLng
	if cosDLng <= 0 {
		if lat > 0 {
			return 90
		}
		return -90
	}
	return math.Atan(math.Tan(lat*rad)/cosDLng) / rad
}

// Distance calculate distance between 2 point lnglat, useful func
func Distance(lng1, lat1, lng2, lat2 float64) float64 {
	h := haverSinDist(lng1, lat1, lng2, lat2, math.Cos(lat1*rad))
	return 2 * earthRadius * math.Asin(math.Sqrt(h))
}
