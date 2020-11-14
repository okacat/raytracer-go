package main

import (
	"image/color"
	"math"
)

// Vector3 is a 3d vector
type Vector3 struct {
	X, Y, Z float64
}

// Add adds this and the other vector and returns the result as a new vector
func (a Vector3) Add(b Vector3) Vector3 {
	return Vector3{
		X: a.X + b.X,
		Y: a.Y + b.Y,
		Z: a.Z + b.Z,
	}
}

// Subtract subtracts this vector with the other vector and returns the result as a new vector
func (a Vector3) Subtract(b Vector3) Vector3 {
	return Vector3{
		X: a.X - b.X,
		Y: a.Y - b.Y,
		Z: a.Z - b.Z,
	}
}

// Scale scales this Vector3 and returns the result as a new vector
func (a Vector3) Scale(s float64) Vector3 {
	return Vector3{
		X: a.X * s,
		Y: a.Y * s,
		Z: a.Z * s,
	}
}

// Length returns the length of this Vector3
func (a Vector3) Length() float64 {
	return math.Sqrt(a.X*a.X + a.Y*a.Y + a.Z*a.Z)
}

// Unit returns a new unit vector
func (a Vector3) Unit() Vector3 {
	return a.Scale(1.0 / a.Length())
}

// Dot returns the dot product of this and the other vector
func (a Vector3) Dot(b Vector3) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

// Cross returns the cross product of this and the other vector in a new vector
func (a Vector3) Cross(b Vector3) Vector3 {
	return Vector3{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

// ToColor converts the vector to a image.Color.RGBA and returns the result
// The values are expected to be [0..1]
func (a Vector3) ToColor() color.Color {
	return color.RGBA{
		R: uint8(a.X * 255),
		G: uint8(a.Y * 255),
		B: uint8(a.Z * 255),
		A: 255,
	}
}