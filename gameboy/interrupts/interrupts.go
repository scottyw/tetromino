package interrupts

// Interrupts captures the current state of interrupts
type Interrupts struct {
	ime bool

	// IE register
	ieHighBits    uint8
	vblankEnabled bool
	statEnabled   bool
	timerEnabled  bool
	serialEnabled bool
	joypadEnabled bool

	// IF register
	vblankRequested bool
	statRequested   bool
	timerRequested  bool
	serialRequested bool
	joypadRequested bool
}

// New Interrupts
func New() *Interrupts {
	i := &Interrupts{}
	i.Enable()
	i.WriteIE(0x00)
	i.WriteIF(0x01)
	return i
}

// IME - Interrupt Master Enable Flag (Write Only)
//   0 - Disable all Interrupts
//   1 - Enable all Interrupts that are enabled in IE Register (FFFF)
// The IME flag is used to disable all interrupts, overriding any enabled bits in
// the IE Register. It isn't possible to access the IME flag by using a I/O
// address, instead IME is accessed directly from the CPU, by the following
// opcodes/operations:
//   EI     ;Enable Interrupts  (ie. IME=1)
//   DI     ;Disable Interrupts (ie. IME=0)
//   RETI   ;Enable Ints & Return (same as the opcode combination EI, RET)
//   <INT>  ;Disable Ints & Call to Interrupt Vector
// Whereas <INT> means the operation which is automatically executed by the CPU
// when it executes an interrupt.

func (i *Interrupts) Enabled() bool {
	return i.ime
}

func (i *Interrupts) Enable() {
	i.ime = true
}

func (i *Interrupts) Disable() {
	i.ime = false
}

// FFFF - IE - Interrupt Enable (R/W)
//   Bit 0: V-Blank  Interrupt Enable  (INT 40h)  (1=Enable)
//   Bit 1: LCD STAT Interrupt Enable  (INT 48h)  (1=Enable)
//   Bit 2: Timer    Interrupt Enable  (INT 50h)  (1=Enable)
//   Bit 3: Serial   Interrupt Enable  (INT 58h)  (1=Enable)
//   Bit 4: Joypad   Interrupt Enable  (INT 60h)  (1=Enable)

// WriteIE handles writes to register IE
func (i *Interrupts) WriteIE(value uint8) {
	// fmt.Printf("> IE - 0x%02x\n", value)
	i.ieHighBits = value & 0xe0
	i.joypadEnabled = value&0x10 > 0
	i.serialEnabled = value&0x08 > 0
	i.timerEnabled = value&0x04 > 0
	i.statEnabled = value&0x02 > 0
	i.vblankEnabled = value&0x01 > 0
}

// ReadIE handles reads from register IE
func (i *Interrupts) ReadIE() uint8 {
	ier := i.ieHighBits
	if i.joypadEnabled {
		ier += 0x10
	}
	if i.serialEnabled {
		ier += 0x08
	}
	if i.timerEnabled {
		ier += 0x04
	}
	if i.statEnabled {
		ier += 0x02
	}
	if i.vblankEnabled {
		ier += 0x01
	}
	// fmt.Printf("< IE - 0x%02x\n", ier)
	return ier
}

// FF0F - IF - Interrupt Flag (R/W)
//   Bit 0: V-Blank  Interrupt Request (INT 40h)  (1=Request)
//   Bit 1: LCD STAT Interrupt Request (INT 48h)  (1=Request)
//   Bit 2: Timer    Interrupt Request (INT 50h)  (1=Request)
//   Bit 3: Serial   Interrupt Request (INT 58h)  (1=Request)
//   Bit 4: Joypad   Interrupt Request (INT 60h)  (1=Request)
// When an interrupt signal changes from low to high, then the corresponding bit
// in the IF register becomes set. For example, Bit 0 becomes set when the LCD
// controller enters into the V-Blank period.

// WriteIF handles writes to register IF
func (i *Interrupts) WriteIF(value uint8) {
	// fmt.Printf("> IF - 0x%02x\n", value)
	i.joypadRequested = value&0x10 > 0
	i.serialRequested = value&0x08 > 0
	i.timerRequested = value&0x04 > 0
	i.statRequested = value&0x02 > 0
	i.vblankRequested = value&0x01 > 0
}

// ReadIF handles reads from register IF
func (i *Interrupts) ReadIF() uint8 {
	// Top 3 bits are always high
	ifr := uint8(0xe0)
	if i.joypadRequested {
		ifr += 0x10
	}
	if i.serialRequested {
		ifr += 0x08
	}
	if i.timerRequested {
		ifr += 0x04
	}
	if i.statRequested {
		ifr += 0x02
	}
	if i.vblankRequested {
		ifr += 0x01
	}
	// fmt.Printf("< IF - 0x%02x\n", ifr)
	return ifr
}

func (i *Interrupts) RequestJoypad() {
	i.joypadRequested = true
}

func (i *Interrupts) RequestSerial() {
	i.serialRequested = true
}

func (i *Interrupts) RequestTimer() {
	i.timerRequested = true
}

func (i *Interrupts) RequestStat() {
	i.statRequested = true
}

func (i *Interrupts) RequestVblank() {
	i.vblankRequested = true
}

func (i *Interrupts) ResetJoypad() {
	i.joypadRequested = false
}

func (i *Interrupts) ResetSerial() {
	i.serialRequested = false
}

func (i *Interrupts) ResetTimer() {
	i.timerRequested = false
}

func (i *Interrupts) ResetStat() {
	i.statRequested = false
}

func (i *Interrupts) ResetVblank() {
	i.vblankRequested = false
}

func (i *Interrupts) JoypadPending() bool {
	return i.joypadEnabled && i.joypadRequested
}

func (i *Interrupts) SerialPending() bool {
	return i.serialEnabled && i.serialRequested
}

func (i *Interrupts) TimerPending() bool {
	return i.timerEnabled && i.timerRequested
}

func (i *Interrupts) StatPending() bool {
	return i.statEnabled && i.statRequested
}

func (i *Interrupts) VblankPending() bool {
	return i.vblankEnabled && i.vblankRequested
}

func (i *Interrupts) Pending() bool {
	return i.JoypadPending() ||
		i.SerialPending() ||
		i.TimerPending() ||
		i.StatPending() ||
		i.VblankPending()
}
