package cpu

// bits := map[string]string{
// 	"0": "01",
// 	"1": "02",
// 	"2": "04",
// 	"3": "08",
// 	"4": "10",
// 	"5": "20",
// 	"6": "40",
// 	"7": "80"}

// nobits := map[string]string{
// 	"0": "fe",
// 	"1": "fd",
// 	"2": "fb",
// 	"3": "f7",
// 	"4": "ef",
// 	"5": "df",
// 	"6": "bf",
// 	"7": "7f"}

// if !unprefixed && opcode.Mnemonic == "BIT" {
// 	fmt.Println(opcode.Mnemonic, opcode.Operand1, opcode.Operand2, opcode.Flags)
// 	fmt.Printf("{0x%2x, CPU{%s: 0x%s}, CPU{%s: 0x%s, zf: false, nf: false, hf: true}, nil},\n",
// 		addr,
// 		strings.ToLower(opcode.Operand2),
// 		bits[opcode.Operand1],
// 		strings.ToLower(opcode.Operand2),
// 		bits[opcode.Operand1])
// 	fmt.Printf("{0x%2x, CPU{%s: 0x%s, zf: true, nf: true, hf: true, cf: true}, CPU{%s: 0x%s, zf: true, nf: false, hf: true, cf: true}, nil},\n",
// 		addr,
// 		strings.ToLower(opcode.Operand2),
// 		nobits[opcode.Operand1],
// 		strings.ToLower(opcode.Operand2),
// 		nobits[opcode.Operand1])
// 	fmt.Println()
// 	fmt.Println()
// }
