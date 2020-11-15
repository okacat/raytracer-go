package main

// World holds the objects in it
type World struct {
	Hittables []Hittable
}

// Hit returns a HitRecord and true if any hits, nil and false otherwise
func (w *World) Hit(r Ray, tMin, tMax float64) (*HitRecord, bool) {
	hitAnything := false
	closestT := tMax
	var hitRecord *HitRecord
	for _, hittable := range w.Hittables {
		record, hit := hittable.Hit(r, tMin, closestT)
		if hit {
			hitAnything = true
			closestT = record.T
			hitRecord = record
		}
	}
	return hitRecord, hitAnything
}
