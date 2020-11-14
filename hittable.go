package main

import "math"

// Hittable is an object that can be hit by a Ray
type Hittable interface {
	Hit(Ray) (Vector3, bool)
}

// HitRecord holds information of a Ray hitting a Hittable object
type HitRecord struct {
	T             float64
	Point, Normal Vector3
}

// Sphere is a Hittable object
type Sphere struct {
	Position Vector3
	Radius   float64
}

// Hit returns the record of the hit if hit and a boolean denoting if the object was hit
func (s Sphere) Hit(r Ray) (*HitRecord, bool) {
	oc := r.Origin.Subtract(s.Position)
	a := r.Direction.LengthSquared()
	bHalf := oc.Dot(r.Direction)
	c := oc.LengthSquared() - s.Radius*s.Radius
	discriminant := bHalf*bHalf - a*c
	if discriminant < 0 {
		return nil, false
	}
	t := (-bHalf - math.Sqrt(discriminant)) / a
	hitPoint := r.At(t)
	normal := hitPoint.Subtract(Vector3{0, 0, -1}).Unit()
	return &HitRecord{
		T:      t,
		Point:  hitPoint,
		Normal: normal,
	}, true
}
