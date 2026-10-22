package main

import (
	"math"
)

func nearestStandardRatio(w, h int) (int, int, error) {
	// List of common aspect ratios
	standards := [][2]int{
		{4, 3},
		{16, 9},
		{9, 16},
		{1, 1},
		{21, 9},
	}

	best := standards[0]
	bestDiff := math.MaxFloat64

	for _, r := range standards {
		ratio := float64(r[0]) / float64(r[1])
		diff := math.Abs(float64(w)/float64(h) - ratio)
		if diff < bestDiff {
			bestDiff = diff
			best = r
		}
	}
	return best[0], best[1], nil

}
