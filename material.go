package main

import "math/rand"

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
	Color Vector3
}

// Scatter returns the scattered ray and it's attenuation
func (m Metal) Scatter(r Ray, h HitRecord, rnd *rand.Rand) (Ray, Vector3, bool) {
	reflected := r.Direction.Unit().Reflect(h.Normal)
	scatteredRay := Ray{
		Origin:    h.Point,
		Direction: reflected,
	}
	hasScattered := reflected.Dot(h.Normal) > 0
	return scatteredRay, m.Color, hasScattered
}
