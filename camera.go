package main

import (
	"math"
	"math/rand"
)

// Camera is a 3d camera
type Camera struct {
	position, horizontal, vertical, lowerLeftCorner Vector3
	aspectRatio, verticalFov, lensRadius            float64
	u, v, w                                         Vector3
}

// NewCamera initializeds and returns a new Camera
func NewCamera(position, lookAt, up Vector3, verticalFov, aperture, focusDistance, width, height float64) Camera {
	theta := Deg2Rad(verticalFov)
	h := math.Tan(theta / 2.0)

	aspectRatio := width / height
	viewportHeight := 2.0 * h
	viewportWidth := viewportHeight * aspectRatio

	w := position.Subtract(lookAt).Unit()
	u := up.Cross(w).Unit()
	v := w.Cross(u)

	horizontal := u.Scale(viewportWidth).Scale(focusDistance)
	vertical := v.Scale(viewportHeight).Scale(focusDistance)
	lowerLeftCorner := position.
		Subtract(horizontal.Scale(0.5)).
		Subtract(vertical.Scale(0.5)).
		Subtract(w.Scale(focusDistance))

	lensRadius := aperture * 0.5

	return Camera{
		position:        position,
		aspectRatio:     aspectRatio,
		horizontal:      horizontal,
		vertical:        vertical,
		lowerLeftCorner: lowerLeftCorner,
		lensRadius:      lensRadius,
		u:               u,
		v:               v,
		w:               w,
	}
}

// GetRay returns a ray going from the camera's position into the scene based on the given u and v
func (c Camera) GetRay(s, t float64, rnd *rand.Rand) Ray {
	random := RandomOnUnitDisk(rnd).Scale(c.lensRadius)
	offset := c.u.Scale(random.X).Add(c.v.Scale(random.Y))
	return Ray{
		Origin: c.position.Add(offset),
		Direction: c.lowerLeftCorner.
			Add(c.horizontal.Scale(s)).
			Add(c.vertical.Scale(t)).
			Subtract(c.position).
			Subtract(offset),
	}
}
