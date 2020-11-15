package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"time"
)

const width = 1000
const height = 600

func main() {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	camera := MakeCamera(Vector3{0, 0, 0}, width, height)

	startTime := time.Now()

	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			u := float64(x) / float64(width-1)
			v := float64(y) / float64(height-1)
			ray := camera.GetRay(u, v)
			clr := rayColor(ray)
			img.Set(x, y, clr)
		}
	}

	fmt.Println("render took ", time.Since(startTime))

	f, error := os.Create("render.png")
	if error != nil {
		fmt.Println(error)
	}
	png.Encode(f, img)
}

func rayColor(r Ray) color.Color {
	hitRecord, hit := Sphere{Vector3{0, 0, -1}, 0.5}.Hit(r)
	switch {
	case hit:
		return hitRecord.Normal.AddScalar(1.0).Scale(0.5).ToColor()
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
