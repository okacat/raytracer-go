package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"math/rand"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
)

const width = 600
const height = 300
const samplesPerPixel = 100
const maxBounces = 50

func main() {
	numThreads := 4
	fmt.Printf("number of available CPUs: %v, spawning %v threads\n", runtime.NumCPU(), numThreads)

	img := createImage()

	world := World{[]Hittable{
		Sphere{
			Position: Vector3{0, 0, -1},
			Radius:   0.5,
			Material: Lambertian{Color: Vector3{0.8, 0.8, 0.8}}},
		Sphere{
			Position: Vector3{-1.1, 0, -1},
			Radius:   0.5,
			Material: Dielectric{IndexOfRefraction: 1.5}},
		Sphere{
			Position: Vector3{1.1, 0, -1},
			Radius:   0.5,
			Material: Metal{Color: Vector3{0.9, 0.2, 0.2}, Glosiness: 0.7}},
		Sphere{
			Position: Vector3{0, -100.5, -1},
			Radius:   100,
			Material: Lambertian{Color: Vector3{0.2, 0.8, 0.2}}},
	}}

	position := Vector3{-2, 2, 1}
	lookAt := Vector3{0, 0, -1}
	up := Vector3{0, 1, 0}
	camera := NewCamera(position, lookAt, up, width, height, 90.0)

	startTime := time.Now()

	var wg sync.WaitGroup
	wg.Add(height)

	jobs := make(chan int)
	progressUpdates := make(chan int)

	go listenForProgress(progressUpdates)

	for i := 0; i < numThreads; i++ {
		rnd := rand.New(rand.NewSource(time.Now().Unix()))
		go lineWorker(world, camera, img, rnd, jobs, progressUpdates, &wg)
	}

	for line := height - 1; line >= 0; line-- {
		jobs <- line
	}

	wg.Wait()

	fmt.Println("render took ", time.Since(startTime).Round(time.Millisecond))
	saveImageAs(img, fmt.Sprintf("render%v.png", time.Now().Unix()))
}

func lineWorker(world World, camera Camera, img *image.RGBA, rnd *rand.Rand, jobs chan int, progressUpdates chan int, wg *sync.WaitGroup) {
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
		progressUpdates <- 1
		wg.Done()
	}
}

func listenForProgress(progressUpdates chan int) {
	linesCompleted := 0
	for p := range progressUpdates {
		linesCompleted += p
		percent := math.Floor(100 * float64(linesCompleted) / float64(height))
		fmt.Printf("rendered %v/%v lines [%v%%]\n", linesCompleted, height, percent)
	}
}

func rayColor(r Ray, w World, depth int, rnd *rand.Rand) Vector3 {
	if depth > 50 {
		return Vector3{0, 0, 0}
	}
	hitRecord, hit := w.Hit(r, 0.001, math.Inf(1))
	if hit {
		// return hitRecord.Normal.Add(Vector3{1, 1, 1}).Scale(0.5) // render normals
		bounceRay, attenuation, hasScattered := hitRecord.Material.Scatter(r, *hitRecord, rnd)
		if hasScattered {
			return rayColor(bounceRay, w, depth+1, rnd).MultiplyComponents(attenuation)
		}
		return Vector3{0, 0, 0}
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

func createImage() *image.RGBA {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	return image.NewRGBA(image.Rectangle{upLeft, lowRight})
}

func saveImageAs(img *image.RGBA, filename string) {
	os.Mkdir("output", 0775)
	f, error := os.Create(path.Join("output", filename))
	if error != nil {
		fmt.Println(error)
	}
	png.Encode(f, img)
}
