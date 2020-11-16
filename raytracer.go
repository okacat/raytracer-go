package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

const width = 500
const height = 250
const samplesPerPixel = 50
const maxBounces = 50

func main() {
	numThreads := 4
	fmt.Printf("number of available CPUs: %v, spawning %v threads\n", runtime.NumCPU(), numThreads)

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	world := World{[]Hittable{
		Sphere{Vector3{0, 0, -1}, 0.5},
		Sphere{Vector3{0, -100.5, -1}, 100},
	}}
	camera := MakeCamera(Vector3{0, 0, 0}, width, height)

	startTime := time.Now()

	// progress := 0
	var wg sync.WaitGroup
	wg.Add(height)

	jobs := make(chan int)

	for i := 0; i < numThreads; i++ {
		fmt.Println("making a thread")
		rnd := rand.New(rand.NewSource(time.Now().Unix()))
		go lineWorker(world, camera, img, rnd, jobs, &wg)
	}

	for line := height - 1; line >= 0; line-- {
		jobs <- line
	}

	wg.Wait()

	fmt.Println("render took ", time.Since(startTime).Round(time.Millisecond))

	f, error := os.Create(fmt.Sprintf("render%v.png", time.Now().Unix()))
	if error != nil {
		fmt.Println(error)
	}
	png.Encode(f, img)
}

func lineWorker(world World, camera Camera, img *image.RGBA, rnd *rand.Rand, jobs chan int, wg *sync.WaitGroup) {
	for y := range jobs {
		for x := 0; x < width; x++ {
			accumulatedColor := Vector3{0, 0, 0}
			for sample := 0; sample < samplesPerPixel; sample++ {
				u := (float64(x) + rand.Float64()) / float64(width-1)
				v := (float64(y) + rand.Float64()) / float64(height-1)
				ray := camera.GetRay(u, v)
				accumulatedColor = accumulatedColor.Add(rayColor(ray, world, 0, rnd))
			}
			pixelColor := accumulatedColor.Scale(1.0 / samplesPerPixel).gammaCorrect().ToColor()
			img.Set(x, height-y, pixelColor)
		}
		wg.Done()
	}
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
