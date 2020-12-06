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

const width = 100 * 20
const height = 100 * 20
const samplesPerPixel = 1000
const maxBounces = 50

func main() {
	numThreads := 8
	fmt.Printf("number of available CPUs: %v, spawning %v threads\n", runtime.NumCPU(), numThreads)

	img := createImage()

	// world := newTestWorldIcoSphere()
	// world := newTestWorldTeapot()
	// world := newTestWorldSphereTriangleLight()
	// world := newTestWorldCornellBox()
	// world := newTestWorldPlanet()
	world := newTestWorldStairs()
	// world := newTestWorldPyramid()

	startTime := time.Now()

	var wg sync.WaitGroup
	wg.Add(height)

	jobs := make(chan int)
	progressUpdates := make(chan int)

	go listenForProgress(progressUpdates)

	for i := 0; i < numThreads; i++ {
		rnd := rand.New(rand.NewSource(time.Now().Unix()))
		go lineWorker(world, img, rnd, jobs, progressUpdates, &wg)
	}

	for line := height - 1; line >= 0; line-- {
		jobs <- line
	}

	wg.Wait()

	fmt.Println("render took ", time.Since(startTime).Round(time.Millisecond))
	saveImageAs(img, fmt.Sprintf("render%v.png", time.Now().Unix()))
}

func lineWorker(world World, img *image.RGBA, rnd *rand.Rand, jobs chan int, progressUpdates chan int, wg *sync.WaitGroup) {
	for y := range jobs {
		for x := 0; x < width; x++ {
			accumulatedColor := Vector3{0, 0, 0}
			for sample := 0; sample < samplesPerPixel; sample++ {
				u := (float64(x) + rand.Float64()) / float64(width-1)
				v := (float64(y) + rand.Float64()) / float64(height-1)
				ray := world.Camera.GetRay(u, v, rnd)
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
		emitted := hitRecord.Material.Emit(r, *hitRecord, rnd)
		bounceRay, attenuation, hasScattered := hitRecord.Material.Scatter(r, *hitRecord, rnd)
		if hasScattered {
			return rayColor(*bounceRay, w, depth+1, rnd).
				MultiplyComponents(attenuation).
				Add(emitted)
		}
		return emitted
	}
	return w.AmbientColor(r)
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
