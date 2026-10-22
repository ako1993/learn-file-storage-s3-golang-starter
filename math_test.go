package main

import (
	"fmt"
	"testing"
)

func TestNearestStandardRatio(t *testing.T) {
	width := 608
	height := 1080
	w, h, err := nearestStandardRatio(width, height)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%v:%v", w, h)

}
