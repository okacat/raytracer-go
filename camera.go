package main

// Camera is a 3d camera
type Camera struct {
	position, horizontal, vertical, lowerLeftCorner Vector3
	aspectRatio                                     float64
}

// NewCamera initializeds and returns a new Camera
func NewCamera(position Vector3, width, height float64) Camera {
	aspectRatio := width / height
	viewportHeight := 2.0
	viewportWidth := viewportHeight * aspectRatio
	focalLength := 1.0

	horizontal := Vector3{viewportWidth, 0, 0}
	vertical := Vector3{0, viewportHeight, 0}
	lowerLeftCorner := position.
		Subtract(horizontal.Scale(0.5)).
		Subtract(vertical.Scale(0.5)).
		Subtract(Vector3{0, 0, focalLength})

	return Camera{
		position:        position,
		aspectRatio:     aspectRatio,
		horizontal:      horizontal,
		vertical:        vertical,
		lowerLeftCorner: lowerLeftCorner,
	}
}

// GetRay returns a ray going from the camera's position into the scene based on the given u and v
func (c Camera) GetRay(u, v float64) Ray {
	return Ray{
		Origin: c.position,
		Direction: c.lowerLeftCorner.
			Add(c.horizontal.Scale(u)).
			Add(c.vertical.Scale(v)).
			Subtract(c.position),
	}
}
