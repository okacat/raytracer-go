package main

import "math"

// Camera is a 3d camera
type Camera struct {
	position, horizontal, vertical, lowerLeftCorner Vector3
	aspectRatio, verticalFov                        float64
}

// NewCamera initializeds and returns a new Camera
func NewCamera(position, lookAt, up Vector3, width, height, verticalFov float64) Camera {
	theta := Deg2Rad(verticalFov)
	h := math.Tan(theta / 2.0)

	aspectRatio := width / height
	viewportHeight := 2.0 * h
	viewportWidth := viewportHeight * aspectRatio

	w := position.Subtract(lookAt).Unit()
	u := up.Cross(w).Unit()
	v := w.Cross(u)

	horizontal := u.Scale(viewportWidth)
	vertical := v.Scale(viewportHeight)
	lowerLeftCorner := position.
		Subtract(horizontal.Scale(0.5)).
		Subtract(vertical.Scale(0.5)).
		Subtract(w)

	return Camera{
		position:        position,
		aspectRatio:     aspectRatio,
		horizontal:      horizontal,
		vertical:        vertical,
		lowerLeftCorner: lowerLeftCorner,
	}
}

// GetRay returns a ray going from the camera's position into the scene based on the given u and v
func (c Camera) GetRay(s, t float64) Ray {
	return Ray{
		Origin: c.position,
		Direction: c.lowerLeftCorner.
			Add(c.horizontal.Scale(s)).
			Add(c.vertical.Scale(t)).
			Subtract(c.position),
	}
}
