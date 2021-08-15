package audio

import "fmt"

type control struct {
	on             bool
	ch1Right       bool
	ch2Right       bool
	ch3Right       bool
	ch4Right       bool
	ch1Left        bool
	ch2Left        bool
	ch3Left        bool
	ch4Left        bool
	vinLeftEnable  bool
	volumeLeft     uint8
	vinRightEnable bool
	volumeRight    uint8
}

// FF10 - NR10 - Channel 1 Sweep register (R/W)
//   Bit 6-4 - Sweep Time
//   Bit 3   - Sweep Increase/Decrease
//              0: Addition    (frequency increases)
//              1: Subtraction (frequency decreases)
//   Bit 2-0 - Number of sweep shift (n: 0-7)
// Sweep Time:
//   000: sweep off - no freq change
//   001: 7.8 ms  (1/128Hz)
//   010: 15.6 ms (2/128Hz)
//   011: 23.4 ms (3/128Hz)
//   100: 31.3 ms (4/128Hz)
//   101: 39.1 ms (5/128Hz)
//   110: 46.9 ms (6/128Hz)
//   111: 54.7 ms (7/128Hz)
//
// The change of frequency (NR13,NR14) at each shift is calculated by the
// following formula where X(0) is initial freq & X(t-1) is last freq:
//   X(t) = X(t-1) +/- X(t-1)/2^n

// WriteNR10 handles writes to sound register NR10
func (a *Audio) WriteNR10(value uint8) {
	if !a.control.on {
		return
	}
	// fmt.Printf("> NR10 - 0x%02x - %+v\n", value, a.ch1)
	a.ch1.sweepPeriod = (value >> 4) & 0x07
	a.ch1.sweepIncrease = (value>>3)&0x01 == 0
	a.ch1.sweepShift = value & 0x07
	if a.ch1.sweepIncrease && a.ch1.sweepDescending {
		a.ch1.enabled = false
	}
	a.ch1.sweepDescending = false
}

// ReadNR10 handles reads from sound register NR10
func (a *Audio) ReadNR10() uint8 {
	nr10 := 0x80 | a.ch1.sweepPeriod<<4 | a.ch1.sweepShift
	if !a.ch1.sweepIncrease {
		nr10 += 0x08
	}
	// fmt.Printf("< NR10 - 0x%02x - %+v\n", nr10, a.ch1)
	return nr10
}

// FF11 - NR11 - Channel 1 Sound length/Wave pattern duty (R/W)
//   Bit 7-6 - Wave Pattern Duty (Read/Write)
//   Bit 5-0 - Sound length data (Write Only) (t1: 0-63)
// Wave Duty:
//   00: 12.5% ( _-------_-------_------- )
//   01: 25%   ( __------__------__------ )
//   10: 50%   ( ____----____----____---- ) (normal)
//   11: 75%   ( ______--______--______-- )
// Sound Length = (64-t1)*(1/256) seconds
// The Length value is used only if Bit 6 in NR14 is set.

// WriteNR11 handles writes to sound register NR11
func (a *Audio) WriteNR11(value uint8) {
	// fmt.Printf("> NR11 - 0x%02x - %+v\n", value, a.ch1)
	if a.control.on {
		a.ch1.duty = value >> 6
	}
	a.ch1.length = 64 - (value & 0x3f)
}

// ReadNR11 handles reads from sound register NR11
func (a *Audio) ReadNR11() uint8 {
	nr11 := 0x3f | a.ch1.duty<<6
	// fmt.Printf("< NR11 - 0x%02x - %v - %+v\n", nr11, a.control.on, a.ch1)
	return nr11
}

// FF12 - NR12 - Channel 1 Volume Envelope (R/W)
//   Bit 7-4 - Initial Volume of envelope (0-0Fh) (0=No Sound)
//   Bit 3   - Envelope Direction (0=Decrease, 1=Increase)
//   Bit 2-0 - Number of envelope sweep (n: 0-7)
//             (If zero, stop envelope operation.)
// Length of 1 step = n*(1/64) seconds

// WriteNR12 handles writes to sound register NR12
func (a *Audio) WriteNR12(value uint8) {
	if !a.control.on {
		return
	}
	// fmt.Printf("> NR12 - 0x%02x - %+v\n", value, a.ch1)
	a.ch1.initialVolume = value >> 4
	a.ch1.envelopeIncrease = (value>>3)&0x01 > 0
	a.ch1.envelopeSweep = value & 0x07
	a.ch1.dacEnabled = a.ch1.initialVolume > 0 || a.ch1.envelopeIncrease
	if !a.ch1.dacEnabled {
		a.ch1.enabled = false
	}
}

// ReadNR12 handles reads from sound register NR12
func (a *Audio) ReadNR12() uint8 {
	nr12 := a.ch1.initialVolume<<4 | a.ch1.envelopeSweep
	if a.ch1.envelopeIncrease {
		nr12 += 0x08
	}
	// fmt.Printf("< NR12 - 0x%02x - %+v\n", nr12, a.ch1)
	return nr12
}

// FF13 - NR13 - Channel 1 Frequency lo (Write Only)
//
// Lower 8 bits of 11 bit frequency (x).
// Next 3 bit are in NR14 ($FF14)

// WriteNR13 handles writes to sound register NR13
func (a *Audio) WriteNR13(value uint8) {
	if !a.control.on {
		return
	}
	// fmt.Printf("> NR13 - 0x%02x - %+v\n", value, a.ch1)
	// fmt .Printf("NR13- 0x%02x\n", value)
	a.ch1.frequency = (a.ch1.frequency & 0xff00) | uint16(value) // Update low byte
	a.ch1.timer = (2048 - a.ch1.frequency) * 4
}

// ReadNR13 handles reads from sound register NR13
func (a *Audio) ReadNR13() uint8 {
	// fmt.Printf("< NR13 - 0xff - %+v\n", a.ch1)
	return 0xff
}

// FF14 - NR14 - Channel 1 Frequency hi (R/W)
//   Bit 7   - Initial (1=Restart Sound)     (Write Only)
//   Bit 6   - Counter/consecutive selection (Read/Write)
//             (1=Stop output when length in NR11 expires)
//   Bit 2-0 - Frequency's higher 3 bits (x) (Write Only)
// Frequency = 131072/(2048-x) Hz

// WriteNR14 handles writes to sound register NR14
func (a *Audio) WriteNR14(value uint8) {
	if !a.control.on {
		return
	}
	// fmt.Printf("> NR14 - 0x%02x - %+v\n", value, a.ch1)
	a.ch1.frequency = (a.ch1.frequency & 0x00ff) | uint16(value&0x07)<<8 // Update high byte
	trigger := (value>>7)&0x01 > 0
	lengthEnable := (value>>6)&0x01 > 0
	if !a.ch1.lengthEnable && lengthEnable && a.ch1.length > 0 && a.frameSeqTicks%2 == 1 {
		a.ch1.length--
		if a.ch1.length == 0 && !trigger {
			a.ch1.enabled = false
		}
	}
	if trigger {
		a.ch1.trigger()
		if lengthEnable && a.ch1.length == 64 && a.frameSeqTicks%2 == 1 {
			a.ch1.length--
		}
	}
	a.ch1.lengthEnable = lengthEnable
}

// ReadNR14 handles reads from sound register NR14
func (a *Audio) ReadNR14() uint8 {
	if a.ch1.lengthEnable {
		// fmt.Printf("< NR14 - 0xff - %+v\n", a.ch1)
		return 0xff
	}
	// fmt.Printf("< NR14 - 0xbf - %+v\n", a.ch1)
	return 0xbf
}

// FF16 - NR21 - Channel 2 Sound Length/Wave Pattern Duty (R/W)
//   Bit 7-6 - Wave Pattern Duty (Read/Write)
//   Bit 5-0 - Sound length data (Write Only) (t1: 0-63)
// Wave Duty:
//   00: 12.5% ( _-------_-------_------- )
//   01: 25%   ( __------__------__------ )
//   10: 50%   ( ____----____----____---- ) (normal)
//   11: 75%   ( ______--______--______-- )
// Sound Length = (64-t1)*(1/256) seconds
// The Length value is used only if Bit 6 in NR24 is set.

// WriteNR21 handles writes to sound register NR21
func (a *Audio) WriteNR21(value uint8) {
	if a.control.on {
		a.ch2.duty = value >> 6
	}
	a.ch2.length = 64 - (value & 0x3f)
}

// ReadNR21 handles reads from sound register NR21
func (a *Audio) ReadNR21() uint8 {
	return 0x3f | a.ch2.duty<<6
}

// FF17 - NR22 - Channel 2 Volume Envelope (R/W)
//   Bit 7-4 - Initial Volume of envelope (0-0Fh) (0=No Sound)
//   Bit 3   - Envelope Direction (0=Decrease, 1=Increase)
//   Bit 2-0 - Number of envelope sweep (n: 0-7)
//             (If zero, stop envelope operation.)
// Length of 1 step = n*(1/64) seconds

// WriteNR22 handles writes to sound register NR22
func (a *Audio) WriteNR22(value uint8) {
	if !a.control.on {
		return
	}
	a.ch2.initialVolume = value >> 4
	a.ch2.envelopeIncrease = (value>>3)&0x01 > 0
	a.ch2.envelopeSweep = value & 0x07
	a.ch2.dacEnabled = a.ch2.initialVolume > 0 || a.ch2.envelopeIncrease
	if !a.ch2.dacEnabled {
		a.ch2.enabled = false
	}
}

// ReadNR22 handles reads from sound register NR22
func (a *Audio) ReadNR22() uint8 {
	nr22 := a.ch2.initialVolume<<4 | a.ch2.envelopeSweep
	if a.ch2.envelopeIncrease {
		nr22 += 0x08
	}
	return nr22
}

// FF18 - NR23 - Channel 2 Frequency lo data (W)
// Frequency's lower 8 bits of 11 bit data (x).
// Next 3 bits are in NR24 ($FF19).

// WriteNR23 handles writes to sound register NR23
func (a *Audio) WriteNR23(value uint8) {
	if !a.control.on {
		return
	}
	a.ch2.frequency = (a.ch2.frequency & 0xff00) | uint16(value) // Update low byte
}

// ReadNR23 handles reads from sound register NR23
func (a *Audio) ReadNR23() uint8 {
	return 0xff
}

// FF19 - NR24 - Channel 2 Frequency hi data (R/W)
//   Bit 7   - Initial (1=Restart Sound)     (Write Only)
//   Bit 6   - Counter/consecutive selection (Read/Write)
//             (1=Stop output when length in NR21 expires)
//   Bit 2-0 - Frequency's higher 3 bits (x) (Write Only)
// Frequency = 131072/(2048-x) Hz

// WriteNR24 handles writes to sound register NR24
func (a *Audio) WriteNR24(value uint8) {
	if !a.control.on {
		return
	}
	a.ch2.frequency = (a.ch2.frequency & 0x00ff) | uint16(value&0x07)<<8 // Update high byte
	trigger := (value>>7)&0x01 > 0
	lengthEnable := (value>>6)&0x01 > 0
	if !a.ch2.lengthEnable && lengthEnable && a.ch2.length > 0 && a.frameSeqTicks%2 == 1 {
		a.ch2.length--
		if a.ch2.length == 0 && !trigger {
			a.ch2.enabled = false
		}
	}
	if trigger {
		a.ch2.trigger()
		if lengthEnable && a.ch2.length == 64 && a.frameSeqTicks%2 == 1 {
			a.ch2.length--
		}
	}
	a.ch2.lengthEnable = lengthEnable
}

// ReadNR24 handles reads from sound register NR24
func (a *Audio) ReadNR24() uint8 {
	if a.ch2.lengthEnable {
		return 0xff
	}
	return 0xbf
}

// FF1A - NR30 - Channel 3 ch Enable/off (R/W)
//   Bit 7 - Sound Channel 3 Off  (0=Stop, 1=Playback)  (Read/Write)

// WriteNR30 handles writes to sound register NR30
func (a *Audio) WriteNR30(value uint8) {
	if !a.control.on {
		return
	}
	a.ch3.dacEnabled = (value>>7)&0x01 > 0
	if !a.ch3.dacEnabled {
		a.ch3.enabled = false
	}
}

// ReadNR30 handles reads from sound register NR30
func (a *Audio) ReadNR30() uint8 {
	if a.ch3.dacEnabled {
		return 0xff
	}
	return 0x7f
}

// FF1B - NR31 - Channel 3 Sound Length
//   Bit 7-0 - Sound length (t1: 0 - 255)
// Sound Length = (256-t1)*(1/256) seconds
// This value is used only if Bit 6 in NR34 is set.

// WriteNR31 handles writes to sound register NR31
func (a *Audio) WriteNR31(value uint8) {
	a.ch3.length = 256 - uint16(value)
}

// ReadNR31 handles reads from sound register NR31
func (a *Audio) ReadNR31() uint8 {
	return 0xff
}

// FF1C - NR32 - Channel 3 Select output level (R/W)
//   Bit 6-5 - Select output level (Read/Write)
// Possible Output levels are:
//   0: Mute (No sound)
//   1: 100% Volume (Produce Wave Pattern RAM Data as it is)
//   2:  50% Volume (Produce Wave Pattern RAM data shifted once to the right)
//   3:  25% Volume (Produce Wave Pattern RAM data shifted twice to the right)

// WriteNR32 handles writes to sound register NR32
func (a *Audio) WriteNR32(value uint8) {
	if !a.control.on {
		return
	}
	a.ch3.outputLevel = (value >> 5) & 0x03
}

// ReadNR32 handles reads from sound register NR32
func (a *Audio) ReadNR32() uint8 {
	return 0x9f | a.ch3.outputLevel<<5
}

// FF1D - NR33 - Channel 3 Frequency's lower data (W)
// Lower 8 bits of an 11 bit frequency (x).

// WriteNR33 handles writes to sound register NR33
func (a *Audio) WriteNR33(value uint8) {
	if !a.control.on {
		return
	}
	a.ch3.frequency = (a.ch3.frequency & 0xff00) | uint16(value) // Update low byte
}

// ReadNR33 handles reads from sound register NR33
func (a *Audio) ReadNR33() uint8 {
	return 0xff
}

// FF1E - NR34 - Channel 3 Frequency's higher data (R/W)
//   Bit 7   - Initial (1=Restart Sound)     (Write Only)
//   Bit 6   - Counter/consecutive selection (Read/Write)
//             (1=Stop output when length in NR31 expires)
//   Bit 2-0 - Frequency's higher 3 bits (x) (Write Only)
// Frequency  =  4194304/(64*(2048-x)) Hz  =  65536/(2048-x) Hz

// WriteNR34 handles writes to sound register NR34
func (a *Audio) WriteNR34(value uint8) {
	if !a.control.on {
		return
	}
	a.ch3.frequency = (a.ch3.frequency & 0x00ff) | uint16(value&0x07)<<8 // Update high byte
	trigger := (value>>7)&0x01 > 0
	lengthEnable := (value>>6)&0x01 > 0
	if !a.ch3.lengthEnable && lengthEnable && a.ch3.length > 0 && a.frameSeqTicks%2 == 1 {
		a.ch3.length--
		if a.ch3.length == 0 && !trigger {
			a.ch3.enabled = false
		}
	}
	if trigger {
		a.ch3.trigger()
		if lengthEnable && a.ch3.length == 256 && a.frameSeqTicks%2 == 1 {
			a.ch3.length--
		}
	}
	a.ch3.lengthEnable = lengthEnable
}

// ReadNR34 handles reads from sound register NR34
func (a *Audio) ReadNR34() uint8 {
	if a.ch3.lengthEnable {
		return 0xff
	}
	return 0xbf
}

// FF20 - NR41 - Channel 4 Sound Length (R/W)
//   Bit 5-0 - Sound length data (t1: 0-63)
// Sound Length = (64-t1)*(1/256) seconds
// The Length value is used only if Bit 6 in NR44 is set.

// WriteNR41 handles writes to sound register NR41
func (a *Audio) WriteNR41(value uint8) {
	a.ch4.length = 64 - (value & 0x3f)
}

// ReadNR41 handles reads from sound register NR41
func (a *Audio) ReadNR41() uint8 {
	return 0xff
}

// FF21 - NR42 - Channel 4 Volume Envelope (R/W)
//   Bit 7-4 - Initial Volume of envelope (0-0Fh) (0=No Sound)
//   Bit 3   - Envelope Direction (0=Decrease, 1=Increase)
//   Bit 2-0 - Number of envelope sweep (n: 0-7)
//             (If zero, stop envelope operation.)
// Length of 1 step = n*(1/64) seconds

// WriteNR42 handles writes to sound register NR42
func (a *Audio) WriteNR42(value uint8) {
	if !a.control.on {
		return
	}
	a.ch4.initialVolume = value >> 4
	a.ch4.envelopeIncrease = (value>>3)&0x01 > 0
	a.ch4.envelopeSweep = value & 0x07
	a.ch4.dacEnabled = a.ch4.initialVolume > 0 || a.ch4.envelopeIncrease
	if !a.ch4.dacEnabled {
		a.ch4.enabled = false
	}
}

// ReadNR42 handles reads from sound register NR42
func (a *Audio) ReadNR42() uint8 {
	nr42 := a.ch4.initialVolume<<4 | a.ch4.envelopeSweep
	if a.ch4.envelopeIncrease {
		nr42 += 0x08
	}
	return nr42
}

// FF22 - NR43 - Channel 4 Polynomial Counter (R/W)
// The amplitude is randomly switched between high and low at the given
// frequency. A higher frequency will make the noise to appear 'softer'.
// When Bit 3 is set, the output will become more regular, and some frequencies
// will sound more like Tone than Noise.
//   Bit 7-4 - Shift Clock Frequency (s)
//   Bit 3   - Counter Step/Width (0=15 bits, 1=7 bits)
//   Bit 2-0 - Dividing Ratio of Frequencies (r)
// Frequency = 524288 Hz / r / 2^(s+1)     ;For r=0 assume r=0.5 instead

// WriteNR43 handles writes to sound register NR43
func (a *Audio) WriteNR43(value uint8) {
	if !a.control.on {
		return
	}
	a.ch4.shift = value >> 4
	a.ch4.lfsrWidth = (value >> 3) & 0x01
	a.ch4.divisor = value & 0x07
}

// ReadNR43 handles reads from sound register NR43
func (a *Audio) ReadNR43() uint8 {
	return a.ch4.shift<<4 | a.ch4.lfsrWidth<<3 | a.ch4.divisor
}

// FF23 - NR44 - Channel 4 Counter/consecutive; Inital (R/W)
//   Bit 7   - Initial (1=Restart Sound)     (Write Only)
//   Bit 6   - Counter/consecutive selection (Read/Write)
//             (1=Stop output when length in NR41 expires)

// WriteNR44 handles writes to sound register NR44
func (a *Audio) WriteNR44(value uint8) {
	if !a.control.on {
		return
	}
	trigger := (value>>7)&0x01 > 0
	lengthEnable := (value>>6)&0x01 > 0
	if !a.ch4.lengthEnable && lengthEnable && a.ch4.length > 0 && a.frameSeqTicks%2 == 1 {
		a.ch4.length--
		if a.ch4.length == 0 && !trigger {
			a.ch4.enabled = false
		}
	}
	if trigger {
		a.ch4.trigger()
		if lengthEnable && a.ch4.length == 64 && a.frameSeqTicks%2 == 1 {
			a.ch4.length--
		}
	}
	a.ch4.lengthEnable = lengthEnable
}

// ReadNR44 handles reads from sound register NR44
func (a *Audio) ReadNR44() uint8 {
	if a.ch4.lengthEnable {
		return 0xff
	}
	return 0xbf
}

// FF24 - NR50 - Channel control / ON-OFF / Volume (R/W)
// The volume bits specify the "Master Volume" for Left/Right sound output.
//   Bit 7   - Output Vin to SO2 terminal (1=Enable)
//   Bit 6-4 - SO2 output level (volume)  (0-7)
//   Bit 3   - Output Vin to SO1 terminal (1=Enable)
//   Bit 2-0 - SO1 output level (volume)  (0-7)
// The Vin signal is received from the game cartridge bus, allowing external
// hardware in the cartridge to supply a fifth sound channel, additionally to the
// gameboys internal four channels. As far as I know this feature isn't used by
// any existing games.

// WriteNR50 handles writes to sound register NR50
func (a *Audio) WriteNR50(value uint8) {
	if !a.control.on {
		return
	}
	a.control.vinLeftEnable = (value>>7)&0x01 > 0
	a.control.volumeLeft = (value >> 4) & 0x07
	a.control.vinRightEnable = (value>>3)&0x01 > 0
	a.control.volumeRight = value & 0x07
}

// ReadNR50 handles reads from sound register NR50
func (a *Audio) ReadNR50() uint8 {
	nr50 := a.control.volumeLeft<<4 | a.control.volumeRight
	if a.control.vinLeftEnable {
		nr50 += 0x80
	}
	if a.control.vinRightEnable {
		nr50 += 0x08
	}
	return nr50
}

// FF25 - NR51 - Selection of Sound output terminal (R/W)
//   Bit 7 - Output sound 4 to SO2 terminal
//   Bit 6 - Output sound 3 to SO2 terminal
//   Bit 5 - Output sound 2 to SO2 terminal
//   Bit 4 - Output sound 1 to SO2 terminal
//   Bit 3 - Output sound 4 to SO1 terminal
//   Bit 2 - Output sound 3 to SO1 terminal
//   Bit 1 - Output sound 2 to SO1 terminal
//   Bit 0 - Output sound 1 to SO1 terminal

// WriteNR51 handles writes to sound register NR51
func (a *Audio) WriteNR51(value uint8) {
	if !a.control.on {
		return
	}
	a.control.ch4Left = value&0x80 > 0
	a.control.ch3Left = value&0x40 > 0
	a.control.ch2Left = value&0x20 > 0
	a.control.ch1Left = value&0x10 > 0
	a.control.ch4Right = value&0x08 > 0
	a.control.ch3Right = value&0x04 > 0
	a.control.ch2Right = value&0x02 > 0
	a.control.ch1Right = value&0x01 > 0
}

// ReadNR51 handles reads from sound register NR51
func (a *Audio) ReadNR51() uint8 {
	var nr51 uint8
	if a.control.ch4Left {
		nr51 += 0x80
	}
	if a.control.ch3Left {
		nr51 += 0x40
	}
	if a.control.ch2Left {
		nr51 += 0x20
	}
	if a.control.ch1Left {
		nr51 += 0x10
	}
	if a.control.ch4Right {
		nr51 += 0x08
	}
	if a.control.ch3Right {
		nr51 += 0x04
	}
	if a.control.ch2Right {
		nr51 += 0x02
	}
	if a.control.ch1Right {
		nr51 += 0x01
	}
	return nr51
}

// FF26 - NR52 - ch Enable/off
// If your GB programs don't use sound then Write 00h to this register to save
// 16% or more on GB power consumption. Disabeling the sound controller by
// clearing Bit 7 destroys the contents of all sound registers. Also, it is not
// possible to access any sound registers (execpt FF26) while the sound
// controller is disabled.
//   Bit 7 - All ch Enable/off  (0: stop all sound circuits) (Read/Write)
//   Bit 3 - Sound 4 ON flag (Read Only)
//   Bit 2 - Sound 3 ON flag (Read Only)
//   Bit 1 - Sound 2 ON flag (Read Only)
//   Bit 0 - Sound 1 ON flag (Read Only)
// Bits 0-3 of this register are Read only status bits, writing to these bits
// does NOT enable/disable sound. The flags get set when sound output is
// restarted by setting the Initial flag (Bit 7 in NR14-NR44), the flag remains
// set until the sound length has expired (if enabled). A volume envelopes which
// has decreased to zero volume will NOT cause the sound flag to go off.

// WriteNR52 handles writes to sound register NR52
func (a *Audio) WriteNR52(value uint8) {
	// fmt.Printf("> NR52 - 0x%02x\n", value)
	if (value >> 7) == 0 {
		a.control.on = true
		a.WriteNR10(0x00)
		a.WriteNR12(0x00)
		a.WriteNR13(0x00)
		a.WriteNR14(0x00)
		a.WriteNR22(0x00)
		a.WriteNR23(0x00)
		a.WriteNR24(0x00)
		a.WriteNR30(0x00)
		a.WriteNR32(0x00)
		a.WriteNR33(0x00)
		a.WriteNR34(0x00)
		a.WriteNR42(0x00)
		a.WriteNR43(0x00)
		a.WriteNR44(0x00)
		a.WriteNR50(0x00)
		a.WriteNR51(0x00)
		a.ch1.duty = 0
		a.ch2.duty = 0
		a.control.on = false
	} else {
		if !a.control.on {
			a.frameSeqTicks = 0
		}
		a.control.on = true
	}
}

// ReadNR52 handles reads from sound register NR52
func (a *Audio) ReadNR52() uint8 {
	nr52 := uint8(0x70)
	if a.control.on {
		nr52 += 0x80
	}
	if a.ch4.enabled {
		nr52 += 0x08
	}
	if a.ch3.enabled {
		nr52 += 0x04
	}
	if a.ch2.enabled {
		nr52 += 0x02
	}
	if a.ch1.enabled {
		nr52 += 0x01
	}
	// fmt.Printf("< NR52 - 0x%02x\n", nr52)
	return nr52
}

// WriteWaveRAM updates the audio channel 3 wave RAM
func (a *Audio) WriteWaveRAM(addr uint16, value uint8) {
	if a.ch3.enabled {
		if a.ch3.sampleTimer < 4 {
			a.ch3.waveram[a.ch3.lastAccessed] = value
		}
	}
	a.ch3.waveram[addr-0xff30] = value
}

// ReadWaveRAM reads the audio channel 3 wave RAM
func (a *Audio) ReadWaveRAM(addr uint16) uint8 {
	if a.ch3.enabled {
		if a.ch3.sampleTimer < 4 {
			fmt.Printf("READ LA : %v - %v - %02x - %04x\n", a.ch3.enabled, a.ch3.sampleTimer, a.ch3.lastAccessed, a.ch3.waveram[a.ch3.lastAccessed])
			return a.ch3.waveram[a.ch3.lastAccessed]
		}
		fmt.Printf("READ FF : %v - %v - %02x - %04x\n", a.ch3.enabled, a.ch3.sampleTimer, a.ch3.lastAccessed, 0xff)
		return 0xff
	}
	fmt.Printf("READ NO : %v - %v - %02x - %04x\n", a.ch3.enabled, a.ch3.sampleTimer, a.ch3.lastAccessed, a.ch3.waveram[addr-0xff30])
	return a.ch3.waveram[addr-0xff30]
}
