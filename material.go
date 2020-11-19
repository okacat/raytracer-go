package main

import (
	"math"
	"math/rand"
)

// Material determines how rays scatter
type Material interface {
	Scatter(Ray, HitRecord, *rand.Rand) (Ray, Vector3, bool)
}

// Lambertian is a diffuse material
type Lambertian struct {
	Color Vector3
}

// Scatter returns the scattered ray and it's attenuation
func (l Lambertian) Scatter(r Ray, h HitRecord, rnd *rand.Rand) (Ray, Vector3, bool) {
	scatterDirection := h.Normal.Add(RandomInUnitHemisphere(h.Normal, rnd))
	// scatterDirection := h.Normal.Add(RandomInUnitSphere(rnd))

	// Catch degenerate scatter direction
	if scatterDirection.IsNearZero() {
		scatterDirection = h.Normal
	}

	scatteredRay := Ray{
		Origin:    h.Point,
		Direction: scatterDirection,
	}
	return scatteredRay, l.Color, true
}

// Metal is a reflective material
type Metal struct {
	Color     Vector3
	Glosiness float64
}

// Scatter returns the scattered ray and it's attenuation
func (m Metal) Scatter(r Ray, h HitRecord, rnd *rand.Rand) (Ray, Vector3, bool) {
	reflected := r.Direction.
		Unit().
		Reflect(h.Normal).
		Add(RandomInUnitSphere(rnd).Scale(1.0 - m.Glosiness))
	scatteredRay := Ray{
		Origin:    h.Point,
		Direction: reflected,
	}
	hasScattered := reflected.Dot(h.Normal) > 0
	return scatteredRay, m.Color, hasScattered
}

// Dielectric is a transparent material than refracts light
type Dielectric struct {
	IndexOfRefraction float64
}

// Scatter returns the scattered ray and it's attenuation
func (d Dielectric) Scatter(r Ray, h HitRecord, rnd *rand.Rand) (Ray, Vector3, bool) {
	refractionRatio := d.IndexOfRefraction
	if h.IsFrontFace {
		refractionRatio = 1.0 / d.IndexOfRefraction
	}

	unitDirection := r.Direction.Unit()
	cosTheta := math.Min(1.0, unitDirection.Scale(-1.0).Dot(h.Normal))
	sinTheta := math.Sqrt(1.0 - cosTheta*cosTheta)

	var newDirection Vector3
	cannotRefract := refractionRatio*sinTheta > 1.0
	if cannotRefract || reflectance(cosTheta, refractionRatio) > rnd.Float64() {
		newDirection = unitDirection.Reflect(h.Normal)
	} else {
		newDirection = unitDirection.Refract(h.Normal, refractionRatio)
	}

	refractedRay := Ray{
		Origin:    h.Point,
		Direction: newDirection,
	}
	return refractedRay, Vector3{1.0, 1.0, 1.0}, true
}

func reflectance(cosine, coefficient float64) float64 {
	// Schlick's approximation
	r0 := (1.0 - coefficient) / (1.0 + coefficient)
	r0 = r0 * r0
	return r0 + (1.0-r0)*math.Pow(1.0-cosine, 5)
}
