package controller

import (
	"fmt"
	"time"
)

// Button represents a direction pad or button control
type Button int

const (
	// Up on the control pad
	Up Button = iota
	//Down on the control pad
	Down Button = iota
	// Left on the control pad
	Left Button = iota
	// Right on the control pad
	Right Button = iota
	// A button
	A Button = iota
	// B button
	B Button = iota
	// Start button
	Start Button = iota
	// Select button
	Select Button = iota
)

// Action represents emulator controls
type Action int

const (
	// TakeScreenshot of the current LCD
	TakeScreenshot Action = iota
)

type Controller struct {
	restartCPU     func()
	joyp           uint8
	directionInput uint8
	buttonInput    uint8
}

func New(restartCPU func()) *Controller {
	return &Controller{
		restartCPU:     restartCPU,
		joyp:           0x0f,
		directionInput: 0x0f,
		buttonInput:    0x0f,
	}
}

func (c *Controller) ReadJOYP() uint8 {
	// Bit 5 - P15 Select Button Keys      (0=Select)
	// Bit 4 - P14 Select Direction Keys   (0=Select)
	// First two bits are always high
	if c.joyp&0x10 == 0 {
		return c.joyp&0xf0 | c.directionInput&0x0f | 0xc0
	}
	if c.joyp&0x20 == 0 {
		return c.joyp&0xf0 | c.buttonInput&0x0f | 0xc0
	}
	return c.joyp | 0xcf
}

func (c *Controller) WriteJOYP(value uint8) {
	c.joyp = value
}

// ButtonAction turns UI key presses into emulator button presses corresponding to the Gameboy controls
func (c *Controller) ButtonAction(button Button, pressed bool) {

	// Start the CPU in case it was stopped waiting for input
	c.restartCPU()

	// Bit 3 - P13 Input Down  or Start    (0=Pressed) (Read Only)
	// Bit 2 - P12 Input Up    or Select   (0=Pressed) (Read Only)
	// Bit 1 - P11 Input Left  or Button B (0=Pressed) (Read Only)
	// Bit 0 - P10 Input Right or Button A (0=Pressed) (Read Only)

	switch button {

	// FIXME it shouldn't be possible to press left and right at once or up and down at once

	case Start:
		if pressed {
			c.buttonInput &^= 0x8
		} else {
			c.buttonInput |= 0x8
		}

	case Select:
		if pressed {
			c.buttonInput &^= 0x4
		} else {
			c.buttonInput |= 0x4
		}

	case B:
		if pressed {
			c.buttonInput &^= 0x2
		} else {
			c.buttonInput |= 0x2
		}

	case A:
		if pressed {
			c.buttonInput &^= 0x1
		} else {
			c.buttonInput |= 0x1
		}

	case Down:
		if pressed {
			c.directionInput &^= 0x8
		} else {
			c.directionInput |= 0x8
		}

	case Up:
		if pressed {
			c.directionInput &^= 0x4
		} else {
			c.directionInput |= 0x4
		}

	case Left:
		if pressed {
			c.directionInput &^= 0x2
		} else {
			c.directionInput |= 0x2
		}

	case Right:
		if pressed {
			c.directionInput &^= 0x1
		} else {
			c.directionInput |= 0x1
		}
	}
}

// EmulatorAction turns UI key presses into actions controlling the emulator itself
func (c *Controller) EmulatorAction(action Action) {
	switch action {
	case TakeScreenshot:
		t := time.Now()
		filename := fmt.Sprintf("tetromino-%d%02d%02d-%02d%02d%02d.png",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second())
		fmt.Println("Writing screenshot to", filename)
		// gb.lcd.Screenshot(filename) //FIXME
	}
}
