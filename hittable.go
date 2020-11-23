package main

import (
	"math"
)

// Hittable is an object that can be hit by a Ray
type Hittable interface {
	Hit(Ray, float64, float64) (*HitRecord, bool)
}

// HitRecord holds information of a Ray hitting a Hittable object
type HitRecord struct {
	Point, Normal Vector3
	T             float64
	IsFrontFace   bool
	Material      Material
}

// NewHitRecord initializes a new HitRecord and returns it
func NewHitRecord(point, normal Vector3, ray Ray, t float64, m Material) HitRecord {
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
		Material:    m,
	}
}

// Sphere is a Hittable object
type Sphere struct {
	Position Vector3
	Radius   float64
	Material Material
}

// Hit returns the record of the hit if hit and a boolean denoting if the object was hit
func (s Sphere) Hit(r Ray, tMin, tMax float64) (*HitRecord, bool) {
	oc := r.Origin.Subtract(s.Position)
	a := r.Direction.LengthSquared()
	bHalf := oc.Dot(r.Direction)
	c := oc.LengthSquared() - s.Radius*s.Radius
	discriminant := bHalf*bHalf - a*c
	if discriminant < 0 {
		return nil, false
	}
	discriminantSquared := math.Sqrt(discriminant)
	root := (-bHalf - discriminantSquared) / a
	if root < tMin || root > tMax {
		root = (-bHalf + discriminantSquared) / a
		if root < tMin || root > tMax {
			return nil, false
		}
	}
	hitPoint := r.At(root)
	normal := hitPoint.Subtract(s.Position).Scale(1.0 / s.Radius)
	hitRecord := NewHitRecord(hitPoint, normal, r, root, s.Material)
	return &hitRecord, true
}

// Triangle is a Hittable object
type Triangle struct {
	V0, V1, V2 Vector3
	N0, N1, N2 Vector3
	Material   Material
}

// Hit returns the record of the hit if hit and a boolean denoting if the object was hit
func (tri Triangle) Hit(r Ray, tMin, tMax float64) (*HitRecord, bool) {
	// Source: https://en.wikipedia.org/wiki/M%C3%B6ller%E2%80%93Trumbore_intersection_algorithm
	epsilon := 0.0000001
	edge1 := tri.V1.Subtract(tri.V0)
	edge2 := tri.V2.Subtract(tri.V0)
	h := r.Direction.Cross(edge2)
	a := edge1.Dot(h)
	if a < epsilon && a > -epsilon {
		return nil, false // This ray is parallel to this triangle.
	}
	f := 1.0 / a
	s := r.Origin.Subtract(tri.V0)
	u := f * s.Dot(h)
	if u < 0.0 || u > 1.0 {
		return nil, false
	}
	q := s.Cross(edge1)
	v := f * r.Direction.Dot(q)
	if v < 0.0 || u+v > 1.0 {
		return nil, false
	}
	t := f * edge2.Dot(q)
	if t > epsilon {
		// normal := edge1.Cross(edge2).Unit()
		normal := tri.N0.Scale(1.0 - u - v).
			Add(tri.N1.Scale(u)).
			Add(tri.N2.Scale(v)).
			Unit()
		hitPoint := r.At(t)
		hitRecord := NewHitRecord(hitPoint, normal, r, t, tri.Material)
		return &hitRecord, true
	}
	return nil, false
}
