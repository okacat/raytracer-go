package main

import (
	"math"
	"math/rand"
)

// Material determines how rays scatter
type Material interface {
	Scatter(Ray, HitRecord, *rand.Rand) (*Ray, Vector3, bool)
	Emit(Ray, HitRecord, *rand.Rand) Vector3
}

// Lambertian is a diffuse material
type Lambertian struct {
	Color Vector3
}

// Scatter returns the scattered ray and it's attenuation
func (l Lambertian) Scatter(r Ray, h HitRecord, rnd *rand.Rand) (*Ray, Vector3, bool) {
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
	return &scatteredRay, l.Color, true
}

// Emit returns black, since Lambertian doesn't emit light
func (l Lambertian) Emit(r Ray, h HitRecord, rnd *rand.Rand) Vector3 {
	return Vector3{0, 0, 0}
}

// Metal is a reflective material
type Metal struct {
	Color     Vector3
	Glosiness float64
}

// Scatter returns the scattered ray and it's attenuation
func (m Metal) Scatter(r Ray, h HitRecord, rnd *rand.Rand) (*Ray, Vector3, bool) {
	reflected := r.Direction.
		Unit().
		Reflect(h.Normal).
		Add(RandomInUnitSphere(rnd).Scale(1.0 - m.Glosiness))
	scatteredRay := Ray{
		Origin:    h.Point,
		Direction: reflected,
	}
	hasScattered := reflected.Dot(h.Normal) > 0
	return &scatteredRay, m.Color, hasScattered
}

// Emit returns black, since Metal doesn't emit light
func (m Metal) Emit(r Ray, h HitRecord, rnd *rand.Rand) Vector3 {
	return Vector3{0, 0, 0}
}

// Dielectric is a transparent material than refracts light
type Dielectric struct {
	IndexOfRefraction float64
}

// Emit returns black, since Dielectric doesn't emit light
func (d Dielectric) Emit(r Ray, h HitRecord, rnd *rand.Rand) Vector3 {
	return Vector3{0, 0, 0}
}

// Scatter returns the scattered ray and it's attenuation
func (d Dielectric) Scatter(r Ray, h HitRecord, rnd *rand.Rand) (*Ray, Vector3, bool) {
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
	return &refractedRay, Vector3{1.0, 1.0, 1.0}, true
}

// Light is an emissive material
type Light struct {
	Emission Vector3
}

// Scatter returns nil and false since Light doesn't bounce or refract rays
func (l Light) Scatter(r Ray, h HitRecord, rnd *rand.Rand) (*Ray, Vector3, bool) {
	return nil, Vector3{0, 0, 0}, false
}

// Emit returns the light's emission, components can be > 1.0
func (l Light) Emit(r Ray, h HitRecord, rnd *rand.Rand) Vector3 {
	return l.Emission
}

func reflectance(cosine, coefficient float64) float64 {
	// Schlick's approximation
	r0 := (1.0 - coefficient) / (1.0 + coefficient)
	r0 = r0 * r0
	return r0 + (1.0-r0)*math.Pow(1.0-cosine, 5)
}
