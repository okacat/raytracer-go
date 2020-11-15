package main

import (
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"os"
	"time"
)

const width = 1000
const height = 600
const samplesPerPixel = 50

func main() {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	camera := MakeCamera(Vector3{0, 0, 0}, width, height)

	startTime := time.Now()

	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			accumulatedColor := Vector3{0, 0, 0}
			for sample := 0; sample < samplesPerPixel; sample++ {
				u := (float64(x) + rand.Float64()) / float64(width-1)
				v := (float64(y) + rand.Float64()) / float64(height-1)
				ray := camera.GetRay(u, v)
				accumulatedColor = accumulatedColor.Add(rayColor(ray))
			}
			pixelColor := accumulatedColor.Scale(1.0 / samplesPerPixel).ToColor()
			img.Set(x, y, pixelColor)
		}
	}

	fmt.Println("render took ", time.Since(startTime))

	f, error := os.Create("render.png")
	if error != nil {
		fmt.Println(error)
	}
	png.Encode(f, img)
}

func rayColor(r Ray) Vector3 {
	hitRecord, hit := Sphere{Vector3{0, 0, -1}, 0.5}.Hit(r)
	switch {
	case hit:
		return hitRecord.Normal.AddScalar(1.0).Scale(0.5)
	default:
		return skyboxColor(r)
	}
}

func skyboxColor(r Ray) Vector3 {
	unitDirection := r.Direction.Unit()
	t := 0.5 * (unitDirection.Y + 1.0)
	a := Vector3{0.5, 0.7, 1.0}.Scale(t)
	b := Vector3{1.0, 1.0, 1.0}.Scale(1.0 - t)
	return a.Add(b)
}
