package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"math/rand"
	"os"
	"time"
)

const width = 1000
const height = 600
const samplesPerPixel = 10
const maxBounces = 50

func main() {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	world := World{[]Hittable{
		Sphere{Vector3{0, 0, -1}, 0.5},
		Sphere{Vector3{0, -100.5, -1}, 100},
	}}
	camera := MakeCamera(Vector3{0, 0, 0}, width, height)

	startTime := time.Now()

	progress := 0
	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			accumulatedColor := Vector3{0, 0, 0}
			for sample := 0; sample < samplesPerPixel; sample++ {
				u := (float64(x) + rand.Float64()) / float64(width-1)
				v := (float64(y) + rand.Float64()) / float64(height-1)
				ray := camera.GetRay(u, v)
				accumulatedColor = accumulatedColor.Add(rayColor(ray, world, 0))
			}
			pixelColor := accumulatedColor.Scale(1.0 / samplesPerPixel).ToColor()
			img.Set(x, height-y, pixelColor)
			newProgress := int(math.Floor(float64(((height-y)*width)+x) / float64(width*height) * 100))
			if newProgress != progress {
				progress = newProgress
				fmt.Printf("progress [%v%%]\n", progress)
			}
		}
	}

	fmt.Println("render took ", time.Since(startTime))

	f, error := os.Create("render.png")
	if error != nil {
		fmt.Println(error)
	}
	png.Encode(f, img)
}

func rayColor(r Ray, w World, depth float64) Vector3 {
	if depth > 50 {
		return Vector3{0, 0, 0}
	}
	hitRecord, hit := w.Hit(r, 0, math.Inf(1))
	if hit {
		target := hitRecord.Point.Add(hitRecord.Normal).Add(RandomInUnitSphere())
		bounceRay := Ray{
			Origin:    hitRecord.Point,
			Direction: target.Subtract(hitRecord.Point),
		}
		return rayColor(bounceRay, w, depth+1)
	}
	return skyboxColor(r)
}

func skyboxColor(r Ray) Vector3 {
	unitDirection := r.Direction.Unit()
	t := 0.5 * (unitDirection.Y + 1.0)
	a := Vector3{0.5, 0.7, 1.0}.Scale(t)
	b := Vector3{1.0, 1.0, 1.0}.Scale(1.0 - t)
	return a.Add(b)
}
