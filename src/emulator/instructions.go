package emulator

import "math/rand"

func (e *Emulator) _0x00E0() {
	e.Gfx = [64 * 32]byte{}
	e.Pc += 2
}

func (e *Emulator) _0x00EE() {
	e.Pc = e.Stack[e.Sp] + 2
	e.Sp--
}

func (e *Emulator) _0x1000(nnn uint16) {
	e.Pc = nnn
}

func (e *Emulator) _0x2000(nnn uint16) {
	e.Sp++
	e.Stack[e.Sp] = e.Pc
	e.Pc = nnn
}

func (e *Emulator) _0x3000(x uint16, nn uint8) {
	if e.V[x] == nn {
		e.Pc += 4
	} else {
		e.Pc += 2
	}
}

func (e *Emulator) _0x4000(x uint16, nn uint8) {
	if e.V[x] != nn {
		e.Pc += 4
	} else {
		e.Pc += 2
	}
}

func (e *Emulator) _0x5000(x, y uint16) {
	if e.V[x] == e.V[y] {
		e.Pc += 4
	} else {
		e.Pc += 2
	}
}

func (e *Emulator) _0x6000(x uint16, nn uint8) {
	e.V[x] = nn
	e.Pc += 2
}

func (e *Emulator) _0x7000(x uint16, nn uint8) {
	e.V[x] += nn
	e.Pc += 2
}

func (e *Emulator) _0x0000(x, y uint16) {
	e.V[x] = e.V[y]
	e.Pc += 2
}

func (e *Emulator) _0x0001(x, y uint16) {
	e.V[x] |= e.V[y]
	e.Pc += 2
}

func (e *Emulator) _0x0002(x, y uint16) {
	e.V[x] &= e.V[y]
	e.Pc += 2
}

func (e *Emulator) _0x0003(x, y uint16) {
	e.V[x] ^= e.V[y]
	e.Pc += 2
}

// Set VF to 01 if a carry occurs
// Set VF to 00 if a carry does not occur
func (e *Emulator) _0x0004(x, y uint16) {
	if e.V[y] > (0xFF - e.V[x]) {
		e.V[0xF] = 1
	} else {
		e.V[0xF] = 0
	}
	e.V[x] += e.V[y]
	e.Pc += 2
}

// Set VF to 00 if a borrow occurs
// Set VF to 01 if a borrow does not occur
func (e *Emulator) _0x0005(x, y uint16) {
	if e.V[y] > e.V[x] {
		e.V[0xF] = 0
	} else {
		e.V[0xF] = 1
	}
	e.V[x] -= e.V[y]
	e.Pc += 2
}

// Set register VF to the least significant bit prior to the shift
func (e *Emulator) _0x0006(x, y uint16) {
	e.V[x] = e.V[y] >> 1
	e.V[0xF] = e.V[y] & 0x01
	e.Pc += 2
}

// Set VF to 00 if a borrow occurs
// Set VF to 01 if a borrow does not occur
func (e *Emulator) _0x0007_1(x, y uint16) {
	if e.V[x] > e.V[y] {
		e.V[0xF] = 0
	} else {
		e.V[0xF] = 1
	}
	e.V[x] = e.V[y] - e.V[x]
	e.Pc += 2
}

// Set register VF to the most significant bit prior to the shift
func (e *Emulator) _0x000E(x, y uint16) {
	e.V[x] = e.V[y] << 1
	e.V[0xF] = e.V[y] & 0x80
	e.Pc += 2
}

func (e *Emulator) _0x9000(x, y uint16) {
	if e.V[x] != e.V[y] {
		e.Pc += 4
	} else {
		e.Pc += 2
	}
}

func (e *Emulator) _0xA000(nnn uint16) {
	e.I = nnn
	e.Pc += 2
}

func (e *Emulator) _0xB000(nnn uint16) {
	e.Pc = nnn + uint16(e.V[0])
	e.Pc += 2
}

func (e *Emulator) _0xC000(x uint16, nn uint8) {
	e.V[x] = uint8(rand.Float32()*255) & nn
	e.Pc += 2
}

// Set VF to 01 if any set pixels are changed to unset, and 00 otherwise
// Get the starting x and y coordinates of the graphics array.
func (e *Emulator) _0xD000(x, y uint16) {
	x = uint16(e.V[x])
	y = uint16(e.V[y])
	e.drawSprite(x, y)
	e.Pc += 2
}

func (e *Emulator) _0x009E(x uint16) {
	if e.Keypad[e.V[x]] == 1 {
		e.Pc += 4
		e.Keypad[e.V[x]] = 0
	} else {
		e.Pc += 2
	}
}

func (e *Emulator) _0x00A1(x uint16) {
	if e.Keypad[e.V[x]] == 0 {
		e.Pc += 4
	} else {
		e.Keypad[e.V[x]] = 0
		e.Pc += 2
	}
}

func (e *Emulator) _0x0007_2(x uint16) {
	e.V[x] = e.DelayTimer
	e.Pc += 2
}

func (e *Emulator) _0x000A(x uint16) {
	for i, k := range e.Keypad {
		if k != 0 {
			e.V[x] = uint8(i)
			e.Pc += 2
			break
		}
	}
	e.Keypad[e.V[x]] = 0
}

func (e *Emulator) _0x0015(x uint16) {
	e.DelayTimer = e.V[x]
	e.Pc += 2
}

func (e *Emulator) _0x0018(x uint16) {
	e.SoundTimer = e.V[x]
	e.Pc += 2
}

func (e *Emulator) _0x001E(x uint16) {
	e.I += uint16(e.V[x])
	e.Pc += 2
}

func (e *Emulator) _0x0029(x uint16) {
	e.I = uint16(e.V[x]) * 5
	e.Pc += 2
}

func (e *Emulator) _0x0033(x uint16) {
	e.Memory[e.I] = e.V[x] / 100
	e.Memory[e.I+1] = (e.V[x] / 10) % 10
	e.Memory[e.I+2] = (e.V[x] % 100) % 10
	e.Pc += 2
}

// i is set to i+x+1 after operation
func (e *Emulator) _0x0065(x uint16) {
	for ind := uint16(0); ind <= x; ind++ {
		e.V[ind] = e.Memory[e.I+ind]
	}
	e.Pc += 2
}

// i is set to i+x+1 after operation
func (e *Emulator) _0x0055(x uint16) {
	for ind := uint16(0); ind <= x; ind++ {
		e.Memory[e.I+ind] = e.V[ind]
	}
	e.Pc += 2
}
