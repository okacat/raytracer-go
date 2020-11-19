package main

import "math"

// Clamp returns x clamped between min and max
func Clamp(x, min, max float64) float64 {
	switch {
	case x < min:
		return min
	case x > max:
		return max
	default:
		return x
	}
}

// Deg2Rad converts degrees to radians
func Deg2Rad(deg float64) float64 {
	return deg * math.Pi / 180.0
}
