package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// ReadObj parses a WaveFront .obj file
func ReadObj(filePath string, material Material) []Triangle {
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
			triangles = append(triangles, parseFace(line, verts, normals, material))
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
func parseFace(line string, verts, normals []Vector3, material Material) Triangle {
	groups := strings.Split(line, " ")
	if len(groups) > 4 {
		log.Fatal(".obj models should be triangulated")
	}

	vertexIndices := make([]int, 3)
	normalIndices := make([]int, 3)
	for i := 1; i < 4; i++ {
		splitGroup := strings.Split(groups[i], "/")

		vertexIndex, err1 := strconv.Atoi(splitGroup[0])
		if err1 != nil {
			log.Fatal("Couldn't parse vertex index as integer")
		}
		vertexIndices[i-1] = vertexIndex

		normalIndex, err2 := strconv.Atoi(splitGroup[2])
		if err2 != nil {
			log.Fatal("Couldn't parse normal index as integer")
		}
		normalIndices[i-1] = normalIndex
	}

	return Triangle{
		V0:       verts[vertexIndices[0]-1],
		V1:       verts[vertexIndices[1]-1],
		V2:       verts[vertexIndices[2]-1],
		N0:       normals[normalIndices[0]-1],
		N1:       normals[normalIndices[1]-1],
		N2:       normals[normalIndices[2]-1],
		Material: material,
	}
}
