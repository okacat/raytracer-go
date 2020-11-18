package main

import "math"

// Hittable is an object that can be hit by a Ray
type Hittable interface {
	Hit(Ray, float64, float64) (*HitRecord, bool)
	GetMaterial() Material
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
	hitRecord := NewHitRecord(hitPoint, normal, r, root, s.GetMaterial())
	return &hitRecord, true
}

// GetMaterial returns the sphere's material
func (s Sphere) GetMaterial() Material {
	return s.Material
}
