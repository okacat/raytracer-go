package main

import "math"

// Hittable is an object that can be hit by a Ray
type Hittable interface {
	Hit(Ray) (Vector3, bool)
}

// Sphere is a Hittable object
type Sphere struct {
	Position Vector3
	Radius   float64
}

// Hit returns the normal and a boolean denoting if the object was hit
func (s Sphere) Hit(r Ray) (Vector3, bool) {
	oc := r.Origin.Subtract(s.Position)
	a := r.Direction.LengthSquared()
	bHalf := oc.Dot(r.Direction)
	c := oc.LengthSquared() - s.Radius*s.Radius
	discriminant := bHalf*bHalf - a*c
	if discriminant < 0 {
		return Vector3{0, 0, 0}, false
	}
	t := (-bHalf - math.Sqrt(discriminant)) / a
	normal := r.At(t).Subtract(Vector3{0, 0, -1}).Unit()
	return normal, true
}
