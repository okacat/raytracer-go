package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

const width = 300
const height = 200

func main() {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	aspectRatio := float64(width) / float64(height)
	viewportHeight := 2.0
	viewportWidth := viewportHeight * aspectRatio
	focalLength := 1.0

	origin := Vector3{0, 0, 0}
	horizontal := Vector3{viewportWidth, 0, 0}
	vertical := Vector3{0, viewportHeight, 0}
	lowerLeftCorner := origin.
		Subtract(horizontal.Scale(0.5)).
		Subtract(vertical.Scale(0.5)).
		Subtract(Vector3{0, 0, focalLength})
	fmt.Println("lower left", lowerLeftCorner)

	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			u := float64(x) / float64(width-1)
			v := float64(y) / float64(height-1)
			rayDirection := lowerLeftCorner.
				Add(horizontal.Scale(u)).
				Add(vertical.Scale(v)).
				Subtract(origin)
			ray := Ray{origin, rayDirection}
			clr := rayColor(ray)
			img.Set(x, y, clr)
		}
	}

	f, error := os.Create("render.png")
	if error != nil {
		fmt.Println(error)
	}
	png.Encode(f, img)
}

func rayColor(r Ray) color.Color {
	normal, hit := hitSphere(r, Vector3{0, 0, -1}, 0.5)
	switch {
	case hit:
		return normal.AddScalar(1.0).Scale(0.5).ToColor()
	default:
		return skyboxColor(r)
	}
}

func skyboxColor(r Ray) color.Color {
	unitDirection := r.Direction.Unit()
	t := 0.5 * (unitDirection.Y + 1.0)
	a := Vector3{0.5, 0.7, 1.0}.Scale(t)
	b := Vector3{1.0, 1.0, 1.0}.Scale(1.0 - t)
	return a.Add(b).ToColor()
}

func hitSphere(ray Ray, center Vector3, radius float64) (Vector3, bool) {
	oc := ray.Origin.Subtract(center)
	a := ray.Direction.Dot(ray.Direction)
	b := oc.Dot(ray.Direction) * 2.0
	c := oc.Dot(oc) - radius*radius
	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return Vector3{0, 0, 0}, false
	}
	t := (-b - math.Sqrt(discriminant)) / (2.0 * a)
	normal := ray.At(t).Subtract(Vector3{0, 0, -1}).Unit()
	return normal, true
}
