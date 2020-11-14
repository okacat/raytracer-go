package main

// Ray represents a ray with an origin and direction
type Ray struct {
	Origin, Direction Vector3
}

// At returns the position on this ray given t
func (r Ray) At(t float64) Vector3 {
	return r.Origin.Add(r.Direction.Scale(t))
}
