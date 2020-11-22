package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// ReadObj parses a WaveFront .obj file
func ReadObj(filePath string) []Triangle {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var verts []Vector3
	var normals []Vector3
	var triangles []Triangle

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "v "):
			verts = append(verts, parseVertex(line))
		case strings.HasPrefix(line, "vn "):
			normals = append(normals, parseNormal(line))
		case strings.HasPrefix(line, "f "):
			triangles = append(triangles, parseFace(line, verts, normals))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Parsed %s\n%d triangles\n", filePath, len(triangles))

	return triangles
}

// format: v x y z
// example: v 0.361800 0.276393 0.262860
func parseVertex(line string) Vector3 {
	var x, y, z float64
	_, err := fmt.Sscanf(line, "v %f %f %f", &x, &y, &z)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Parsed vertex %v\n", Vector3{x, y, z})
	return Vector3{x, y, z}
}

// format: vn x y z
// example: vn 0.1876 -0.7947 0.5774
func parseNormal(line string) Vector3 {
	var x, y, z float64
	_, err := fmt.Sscanf(line, "vn %f %f %f", &x, &y, &z)
	if err != nil {
		log.Fatal(err)
	}
	return Vector3{x, y, z}
}

// format: f v1/vt1/vn1 v2/vt2/vn2 v3/vt3/vn3
// example: f 5/5/2 6/6/2 7/7/2
// note: .obj is 1-indexed
func parseFace(line string, verts, normals []Vector3) Triangle {
	var v0, v1, v2, t0, t1, t2, n0, n1, n2 int64
	_, err := fmt.Sscanf(line, "f %d/%d/%d %d/%d/%d %d/%d/%d", &v0, &t0, &n0, &v1, &t1, &n1, &v2, &t2, &n2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Parsed face %d %d %d\n", v0, v1, v2)
	return Triangle{
		V0:       verts[v0-1],
		V1:       verts[v1-1],
		V2:       verts[v2-1],
		Material: Lambertian{Color: Vector3{0.8, 0.8, 0.8}},
	}
}
