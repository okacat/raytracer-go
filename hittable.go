package main

import "math"

// Hittable is an object that can be hit by a Ray
type Hittable interface {
	Hit(Ray) (Vector3, bool)
}

// HitRecord holds information of a Ray hitting a Hittable object
type HitRecord struct {
	Point, Normal Vector3
	T             float64
	IsFrontFace   bool
}

// MakeHitRecord initializes a new HitRecord and returns it
func MakeHitRecord(point, normal Vector3, ray Ray, t float64) HitRecord {
	rayDotNormal := ray.Direction.Dot(normal)
	isFrontFace := rayDotNormal < 0
	outwardNormal := normal
	if !isFrontFace {
		outwardNormal = normal.Scale(-1)
	}
	return HitRecord{
		Point:       point,
		Normal:      outwardNormal,
		T:           t,
		IsFrontFace: isFrontFace,
	}
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
	hitRecord := MakeHitRecord(hitPoint, normal, r, t)
	return &hitRecord, true
}
