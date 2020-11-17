package main

import "math/rand"

// Material determines how rays scatter
type Material interface {
	Scatter(Ray, HitRecord, *rand.Rand) (Ray, Vector3)
}

// Lambertian is a diffuse material
type Lambertian struct {
	Albedo Vector3
}

// Scatter returns the scattered ray and it's attenuation
func (l Lambertian) Scatter(r Ray, h HitRecord, rnd *rand.Rand) (Ray, Vector3) {
	scatterDirection := h.Normal.Add(RandomInUnitHemisphere(h.Normal, rnd))

	// Catch degenerate scatter direction
	if scatterDirection.IsNearZero() {
		scatterDirection = h.Normal
	}

	scatteredRay := Ray{
		Origin:    h.Point,
		Direction: scatterDirection,
	}
	return scatteredRay, l.Albedo
}
