package main

import (
	"fmt"
	"strings"
)

var count = 0

func LogicDoWork(input string) string {
	count++
	fmt.Println("Go logic Start:", input, "count:", count)
	if strings.TrimSpace(input) == "" {
		return "Type Any text"
	}
	return "Run Now: " + input + " 🔥"
}

func LogicLoadFile(path string) string {
	return "Loaded: " + path
}

func LogicGetCount() int { return count }
func LogicGetProgress() float64 {
	if count > 10 {
		return 1.0
	}
	return float64(count) / 10.0
}
