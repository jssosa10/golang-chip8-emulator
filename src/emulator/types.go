package emulator

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jssosa10/golang-chip8-emulator/src/graphics"
)

const (
	maxRomSize   = 0xFFF - 0x200
	keyRepeatDur = time.Second / 5
)

var chip8FontSet = [80]uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

type Emulator struct {
	Opcode uint16
	Memory [4096]uint8
	// CPU registries
	V [16]uint8
	// Index registry
	I uint16
	// program counter
	Pc         uint16
	Gfx        [64 * 32]uint8
	DelayTimer uint8
	SoundTimer uint8
	Stack      [16]uint16
	Sp         uint16
	Keypad     [16]uint8
	Window     *graphics.Window
	drawFlag   bool
}

func (e *Emulator) Init() {
	e.Pc = 0x200
	e.Opcode = 0x0000
	e.I = 0
	e.Sp = 0
	e.drawFlag = false

	//clear Memory
	for i := 0; i < 4096; i++ {
		e.Memory[i] = 0
	}
	for i := 0; i < 16; i++ {
		e.Keypad[i] = 0
		e.Stack[i] = 0
		e.V[i] = 0
	}
	for i := 0; i < 64*32; i++ {
		e.Gfx[i] = 0
	}

	for i := 0; i < 80; i++ {
		e.Memory[i] = chip8FontSet[i]

	}
	window, err := graphics.NewWindow()
	if err != nil {

	}
	e.Window = window
	window.DrawGraphics(e.Gfx)

}

func (e *Emulator) LoadProgram(path string) error {
	rom, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if len(rom) > maxRomSize {
		return errors.New("error: rom too large. Max size: 3583")
	}

	for i := range len(rom) {
		e.Memory[0x200+i] = rom[i] // Write memory with pc offset
	}

	return nil
}

func (e *Emulator) Cycle() {

	e.fetchOpcode()
	err := e.parseOpcode()
	if err != nil {
		fmt.Printf("error parsing opcode: %v", err)
	}
	e.drawOrUpdate()
	e.handleKeyInput()
	e.updateTimers()

}

func (e *Emulator) fetchOpcode() {
	e.Opcode = uint16(e.Memory[e.Pc])<<8 | uint16(e.Memory[e.Pc+1])
}

func (e *Emulator) updateTimers() {
	if e.DelayTimer > 0 {
		e.DelayTimer = e.DelayTimer - 1
	}
	if e.SoundTimer > 0 {
		e.SoundTimer = e.SoundTimer - 1
		if e.SoundTimer == 0 {
			// Play Sound
		}
	}

}

func (e *Emulator) drawSprite(x, y uint16) {
	height := e.Opcode & 0x000F
	e.V[0xF] = 0
	var pix uint16

	for yLine := uint16(0); yLine < height; yLine++ {
		pix = uint16(e.Memory[e.I+yLine])

		for xLine := uint16(0); xLine < 8; xLine++ {
			ind := (x + xLine + ((y + yLine) * 64))
			if ind >= uint16(len(e.Gfx)) {
				continue
			}
			if (pix & (0x80 >> xLine)) != 0 {
				if e.Gfx[ind] == 1 {
					e.V[0xF] = 1
				}
				e.Gfx[ind] ^= 1
			}
		}
	}

	e.drawFlag = true
}

func (e *Emulator) setKeyDown(index byte) {
	e.Keypad[index] = 1
}

func (e Emulator) unknownOp(opcode uint16) error {
	return fmt.Errorf("unknown opcode: %x", opcode)
}

func (e *Emulator) handleKeyInput() {
	for i, key := range e.Window.KeyMap {
		if e.Window.JustReleased(key) && e.Window.KeysDown[i] != nil {
			e.Window.KeysDown[i].Stop()
			e.Window.KeysDown[i] = nil
		} else if e.Window.JustPressed(key) {
			if e.Window.KeysDown[i] == nil {
				e.Window.KeysDown[i] = time.NewTicker(keyRepeatDur)
			}
			e.setKeyDown(byte(i))
		}

		if e.Window.KeysDown[i] == nil {
			continue
		}

		select {
		case <-e.Window.KeysDown[i].C:
			e.setKeyDown(uint8(i))
		default:
		}
	}
}

func (e *Emulator) drawOrUpdate() {
	if e.drawFlag {
		e.Window.DrawGraphics(e.Gfx)
	} else {
		e.Window.UpdateInput()
	}
}

func (e *Emulator) parseOpcode() error {
	x := (e.Opcode & 0x0F00) >> 8 // Decode Vx register identifier.
	y := (e.Opcode & 0x00F0) >> 4 // Decode Vy register identifier
	nn := byte(e.Opcode & 0x00FF) // load last 8-bits
	nnn := e.Opcode & 0x0FFF      // load last 12-bits

	switch e.Opcode & 0xF000 {
	case 0x0000: // 0NNN -> Execute machine language subroutine at address NNN
		switch e.Opcode & 0x00FF {
		case 0x00E0:
			e._0x00E0() // 00E0 -> Clear the screen
		case 0x00EE:
			e._0x00EE() // 00EE -> Return from a subroutine.
		default:
			return e.unknownOp(e.Opcode & 0x00FF)
		}
	case 0x1000:
		e._0x1000(nnn) // 1NNN -> Jump to address NNN
	case 0x2000:
		e._0x2000(nnn) // 2NNN -> Execute subroutine starting at address NNN
	case 0x3000:
		e._0x3000(x, nn) // 3XNN -> Skip the following instruction if the value of register VX == NN
	case 0x4000:
		e._0x4000(x, nn) // 4XNN -> Skip the following instruction if the value of register VX != NN
	case 0x5000:
		e._0x5000(x, y) // 5XY0 -> Skip the following instruction if the value of register VX == VY
	case 0x6000:
		e._0x6000(x, nn) // 6XNN -> Store number NN in register VX
	case 0x7000:
		e._0x7000(x, nn) // 7XNN -> Add the value NN to register VX
	case 0x8000:
		switch e.Opcode & 0x000F {
		case 0x0000:
			e._0x0000(x, y) // 8XYO -> Store the value of register VY in register VX
		case 0x0001:
			e._0x0001(x, y) // 8XY1 -> Set VX to VX OR VY
		case 0x0002:
			e._0x0002(x, y) // 8XY2 -> Set VX to VX AND VY
		case 0x0003:
			e._0x0003(x, y) // 8XY3 -> Set VX to VX XOR VY
		case 0x0004:
			e._0x0004(x, y) // 8XY4 -> Add the value of register VY to register VX
		case 0x0005:
			e._0x0005(x, y) // 8XY5 -> Subtract the value of register VY from register VX
		case 0x0006:
			e._0x0006(x, y) // 8XY6 -> Store the value of register VY shifted right one bit in register VX
		case 0x0007:
			e._0x0007_1(x, y) // 8XY7 -> Set register VX to the value of VY minus VX
		case 0x000E:
			e._0x000E(x, y) // 8XYE -> Store the value of register VY shifted left one bit in register VX
		default:
			return e.unknownOp(e.Opcode & 0x000F)
		}
	case 0x9000:
		e._0x9000(x, y) // 9XY0 -> Skip the following instruction if the value of VX != value of VY
	case 0xA000:
		e._0xA000(nnn) // ANNN -> Store memory address NNN in index register
	case 0xB000:
		e._0xB000(nnn) // BNNN -> Jump to address NNN + V0
	case 0xC000:
		e._0xC000(x, nn) // CXNN -> Set VX to a random number from 0-255 with a mask of NN
	case 0xD000:
		e._0xD000(x, y) // DXYN -> Draw a sprite at position VX, VY with N bytes of sprite data starting at the address stored in index register
	case 0xE000:
		switch e.Opcode & 0x00FF {
		case 0x009E:
			e._0x009E(x) // EX9E -> Skip the following instruction if the key corresponding to the hex value currently stored in register VX is pressed
		case 0x00A1:
			e._0x00A1(x) // EXA1 -> Skip the following instruction if the key corresponding to the hex value currently stored in register VX is not pressed
		default:
			return e.unknownOp(e.Opcode & 0x00FF)
		}
	case 0xF000:
		switch e.Opcode & 0x00FF {
		case 0x0007:
			e._0x0007_2(x) // FX07 -> Store the current value of the delay timer in register VX
		case 0x000A:
			e._0x000A(x) // FX0A -> Wait for a keypress and store the result in register VX
		case 0x0015:
			e._0x0015(x) // FX15 -> Set the delay timer to the value of register VX
		case 0x0018:
			e._0x0018(x) // FX18 -> Set the sound timer to the value of register VX
		case 0x001E:
			e._0x001E(x) // FX1E -> Add the value stored in register VX to index register
		case 0x0029:
			e._0x0029(x) // FX29 -> Set index register to the memory address of the sprite data corresponding to the hexadecimal digit stored in register VX
		case 0x0033:
			e._0x0033(x) // FX33 -> Store the binary-coded decimal equivalent of the value stored in register VX at addresses i, i+1, and i+2
		case 0x0055:
			e._0x0055(x) // FX55 -> Store the values of registers V0 to VX inclusive in memory starting at address i
		case 0x0065:
			e._0x0065(x) // FX65 -> Fill registers V0 to VX inclusive with the values stored in memory starting at address i
		default:
			return e.unknownOp(e.Opcode & 0x00FF)
		}
	default:
		return e.unknownOp(e.Opcode & 0x00FF)
	}
	return nil
}
