package graphics

import (
	"fmt"
	"time"

	pixel "github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"golang.org/x/image/colornames"
)

type Window struct {
	*opengl.Window
	KeyMap   map[uint16]pixel.Button
	KeysDown [16]*time.Ticker
}

func NewWindow() (*Window, error) {
	cfg := opengl.WindowConfig{
		Title:  "Chip 8 Emulator",
		Bounds: pixel.R(0, 0, 640, 320),
		VSync:  true,
	}
	w, err := opengl.NewWindow(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating new window: %v", err)
	}
	km := map[uint16]pixel.Button{
		0x1: pixel.Key1, 0x2: pixel.Key2,
		0x3: pixel.Key3, 0xC: pixel.Key4,
		0x4: pixel.KeyQ, 0x5: pixel.KeyW,
		0x6: pixel.KeyE, 0xD: pixel.KeyR,
		0x7: pixel.KeyA, 0x8: pixel.KeyS,
		0x9: pixel.KeyD, 0xE: pixel.KeyF,
		0xA: pixel.KeyZ, 0x0: pixel.KeyX,
		0xB: pixel.KeyC, 0xF: pixel.KeyV,
	}
	return &Window{
		Window:   w,
		KeyMap:   km,
		KeysDown: [16]*time.Ticker{},
	}, nil
}

func (w *Window) DrawGraphics(gfx [64 * 32]uint8) {
	w.Clear(colornames.Black)
	imDraw := imdraw.New(nil)
	imDraw.Color = pixel.RGB(1, 1, 1)
	width, height := 10, 10

	for i := 0; i < 64; i++ {
		for j := 0; j < 32; j++ {
			// If the gfx byte in question is turned off,
			// continue and skip drawing the rectangle
			if gfx[(31-j)*64+i] == 0 {
				continue
			}
			imDraw.Push(pixel.V(float64(width*i), float64(height*j)))
			imDraw.Push(pixel.V(float64(width*i+width), float64(height*j+height)))
			imDraw.Rectangle(0)
		}
	}

	imDraw.Draw(w)
	w.Update()
}
