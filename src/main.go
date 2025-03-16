package main

import (
	"time"

	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/jssosa10/golang-chip8-emulator/src/emulator"
)

func main() {
	opengl.Run(run)
}

func run() {
	emulator := emulator.Emulator{}
	emulator.Init()
	emulator.LoadProgram("./roms/pong.ch8")
	defer emulator.Window.Destroy()
	done := make(chan struct{})
	ticker := time.NewTicker(1 * time.Second / time.Duration(60))
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if !emulator.Window.Closed() {
					emulator.Cycle()
				}
			}
		}
	}()
	for !emulator.Window.Closed() {

	}
	ticker.Stop()
	done <- struct{}{}

}
