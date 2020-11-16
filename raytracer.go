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

const width = 500
const height = 250
const samplesPerPixel = 50
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
			ch := make(chan Vector3, samplesPerPixel)
			accumulatedColor := Vector3{0, 0, 0}
			// rnd := rand.New(rand.NewSource(time.Now().Unix()))
			for sample := 0; sample < cap(ch); sample++ {
				u := (float64(x) + rand.Float64()) / float64(width-1)
				v := (float64(y) + rand.Float64()) / float64(height-1)
				ray := camera.GetRay(u, v)
				go rayColorPar(ray, world, ch)
				// accumulatedColor = accumulatedColor.Add(rayColor(ray, world, 0, rnd))
			}
			for n := 0; n < cap(ch); n++ {
				sampleColor := <-ch
				accumulatedColor = accumulatedColor.Add(sampleColor)
			}
			pixelColor := accumulatedColor.Scale(1.0 / samplesPerPixel).gammaCorrect().ToColor()
			img.Set(x, height-y, pixelColor)
			newProgress := int(math.Floor(float64(((height-y)*width)+x) / float64(width*height) * 100))
			if newProgress != progress {
				progress = newProgress
				fmt.Printf("progress [%v%%]\n", progress)
			}
		}
	}

	fmt.Println("render took ", time.Since(startTime).Round(time.Millisecond))

	f, error := os.Create(fmt.Sprintf("render%v.png", time.Now().Unix()))
	if error != nil {
		fmt.Println(error)
	}
	png.Encode(f, img)
}

func rayColorPar(r Ray, w World, c chan Vector3) {
	rnd := rand.New(rand.NewSource(time.Now().Unix()))
	c <- rayColor(r, w, 0, rnd)
}

func rayColor(r Ray, w World, depth float64, rnd *rand.Rand) Vector3 {
	if depth > 50 {
		return Vector3{0, 0, 0}
	}
	hitRecord, hit := w.Hit(r, 0.0, math.Inf(1))
	if hit {
		target := hitRecord.Point.Add(hitRecord.Normal).Add(RandomInUnitSphere(rnd))
		bounceRay := Ray{
			Origin:    hitRecord.Point,
			Direction: target.Subtract(hitRecord.Point),
		}
		return rayColor(bounceRay, w, depth+1, rnd)
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

func (v Vector3) gammaCorrect() Vector3 {
	return Vector3{
		X: math.Sqrt(v.X),
		Y: math.Sqrt(v.Y),
		Z: math.Sqrt(v.Z),
	}
}
