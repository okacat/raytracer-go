# RAYTRACER
A raytracer written in Go based on https://raytracing.github.io/books/RayTracingInOneWeekend.html

## Running
Run `go run *.go` or `go build` and then run the binary `./raytracer`

## Features
- unidirectional path tracing
- spheres and triangles as primitives
- diffuse, glossy and refractive materials
- positionable camera with depth of field
- lights
- **very** basic `.obj` parsing, supports triangulated meshes only

### Future wish list
- bounding volume hiearchies to improve speed
- better sampling of lights (light sampling or bidirectional path tracing)
- scene graph
