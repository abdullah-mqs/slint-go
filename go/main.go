package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

// GUI event loop Requird OS thread
func init() {
	runtime.LockOSThread()
}

// findFile
func findFile(rel string) string {
	for _, p := range []string{rel, "../" + rel} {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return rel
}

func main() {
	b := NewBridge()
	b.Init()

	// Callbacks: Slint -> Go (LoadUI and Register)
	b.OnCallback("do-work", func(input string) string {
		return LogicDoWork(input)
	})
	b.OnCallback("load-file", func(path string) string {
		return LogicLoadFile(path)
	})
	b.OnCallback("open-video", func(url string) string {
		return "playing"
	})

	// 1) First UI load
	if err := b.LoadUI(findFile("ui/app.slint")); err != nil {
		log.Fatal(err)
	}

	// 2) Basic props
	b.SetString("my-data", "Hello! Slint + Go now started")
	b.SetInt("count", 7)
	b.SetBool("is-loading", false)
	b.SetFloat("progress", 0.5)
	b.SetColor("bg-color", "#1A1A1A")
	b.SetImage("avatar", findFile("ui/avatar.png"))

	// 3) Array: Go slice -> Slint [type]
	b.SetStringArray("names", []string{"Ali", "Ahmed", "Jani"})
	b.SetIntArray("scores", []int{10, 20, 30})
	b.SetFloatArray("weights", []float64{0.25, 0.5, 0.75})

	p := NewStruct()
	p.SetString("name", "Jani")
	p.SetInt("score", 99)
	p.SetColor("tint", "#FF5733")
	b.SetStruct("player", p)
	p.Free()

	b.SetEnum("mode", "beta")

	fmt.Println("App chal rahi hai... window dekho")
	b.Run()
}
