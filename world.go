package main

// World holds the objects in it
type World struct {
	Camera    Camera
	Hittables []Hittable
}

// Hit returns a HitRecord and true if any hits, nil and false otherwise
func (w *World) Hit(r Ray, tMin, tMax float64) (*HitRecord, bool) {
	hitAnything := false
	closestT := tMax
	var hitRecord *HitRecord
	for _, hittable := range w.Hittables {
		record, hit := hittable.Hit(r, tMin, closestT)
		if hit && closestT > record.T {
			hitAnything = true
			closestT = record.T
			hitRecord = record
		}
	}
	return hitRecord, hitAnything
}

// NewTestWorld returns a new test world
func NewTestWorld() World {
	position := Vector3{1, 0, 0}
	lookAt := Vector3{0, 0, -1.0}
	up := Vector3{0, 1, 0}
	aperture := 1.0 / 16.0
	focusDistance := position.Subtract(lookAt).Length()
	camera := NewCamera(position, lookAt, up, 90.0, aperture, focusDistance, width, height)
	return World{
		Camera: camera,
		Hittables: []Hittable{
			Triangle{
				V0:       Vector3{-2.0, -1.0, -2.5},
				V1:       Vector3{0.0, 2.0, -2.5},
				V2:       Vector3{2.0, -1.0, -2.5},
				Material: Metal{Color: Vector3{0.8, 0.8, 0.8}, Glosiness: 0.99}},
			Sphere{
				Position: Vector3{0, 0, -1},
				Radius:   0.5,
				Material: Lambertian{Color: Vector3{0.8, 0.8, 0.8}}},
			Sphere{
				Position: Vector3{0, -100.5, -1},
				Radius:   100,
				Material: Lambertian{Color: Vector3{0.2, 0.8, 0.2}}},
		}}
}

// NewTestWorldTriangles returns a new test world
func NewTestWorldTriangles() World {
	position := Vector3{0, 0, 0}
	lookAt := Vector3{0, 0, -1.0}
	up := Vector3{0, 1, 0}
	aperture := 0.0
	focusDistance := position.Subtract(lookAt).Length()
	camera := NewCamera(position, lookAt, up, 90.0, aperture, focusDistance, width, height)
	return World{
		Camera: camera,
		Hittables: []Hittable{
			Triangle{
				V0:       Vector3{-2.0, -1.0, -2.5},
				V1:       Vector3{0.0, -1.0, -2.5},
				V2:       Vector3{-1.0, 1.0, -2.5},
				Material: Lambertian{Color: Vector3{0.8, 0.2, 0.2}}},
			Triangle{
				V0:       Vector3{-1.5, -1.0, -2.8},
				V1:       Vector3{0.5, -1.0, -2.8},
				V2:       Vector3{-0.5, 1.0, -2.8},
				Material: Lambertian{Color: Vector3{0.2, 0.8, 0.2}}},
			Triangle{
				V0:       Vector3{-1.0, -1.0, -3.0},
				V1:       Vector3{1.0, -1.0, -3.0},
				V2:       Vector3{0.0, 1.0, -3.0},
				Material: Lambertian{Color: Vector3{0.2, 0.2, 0.8}}},
			Sphere{
				Position: Vector3{0, 0, -1},
				Radius:   0.1,
				Material: Lambertian{Color: Vector3{0.8, 0.8, 0.8}}},
			Sphere{
				Position: Vector3{0, -100.5, -1},
				Radius:   100,
				Material: Lambertian{Color: Vector3{0.2, 0.8, 0.2}}},
		}}
}

// NewTestWorldFromObj returns a new test world reading from obj
func NewTestWorldFromObj() World {
	position := Vector3{0, 0.5, 1}
	lookAt := Vector3{0, 0, -1.0}
	up := Vector3{0, 1, 0}
	aperture := 1.0 / 16.0
	focusDistance := position.Subtract(lookAt).Length()
	camera := NewCamera(position, lookAt, up, 90.0, aperture, focusDistance, width, height)
	triangles := ReadObj("test_scene.obj")
	hittables := make([]Hittable, len(triangles))
	for i := range triangles {
		hittables[i] = triangles[i]
	}
	return World{camera, hittables}
}

func NewTestWorldIcoSphere() World {
	position := Vector3{0, 0.5, 1}
	lookAt := Vector3{0, 0, -1.0}
	up := Vector3{0, 1, 0}
	aperture := 1.0 / 16.0
	focusDistance := position.Subtract(lookAt).Length()
	camera := NewCamera(position, lookAt, up, 90.0, aperture, focusDistance, width, height)
	triangles := ReadObj("icosphere.obj")
	hittables := make([]Hittable, len(triangles))
	for i := range triangles {
		hittables[i] = triangles[i]
	}
	hittables = append(hittables, Sphere{
		Position: Vector3{0, -100.5, -1},
		Radius:   100,
		Material: Lambertian{Color: Vector3{0.2, 0.8, 0.2}}})
	return World{camera, hittables}
}
