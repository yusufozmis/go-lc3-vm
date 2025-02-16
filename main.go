package main

import (
	"fmt"
	"log"
	"os"

	"github.com/eiannone/keyboard"
)

func mem_write(address, val uint16) {
	memory[address] = val
}

func mem_read(address uint16) uint16 {
	if address == MR_KBSR {
		if check_key() != 0 {
			memory[MR_KBSR] = (1 << 15)
			memory[MR_KBDR] = check_key()
		} else {
			memory[MR_KBSR] = 0
		}
	}
	return memory[address]
}

func instr(running int) int {
	for running == 1 {
		var instruction uint16 = mem_read(register[PC])
		var opcode uint16 = instruction >> 12
		register[PC]++
		switch opcode {
		case OP_ADD:
			var r0 uint16 = (instruction >> 9) & 0x7
			var r1 uint16 = (instruction >> 6) & 0x7
			var imm_flag uint16 = (instruction >> 5) & 0x1

			if imm_flag == 1 {
				var imm5 uint16 = sign_extent(instruction&0x1f, 5)
				register[r0] = register[r1] + imm5
			} else {
				var r2 uint16 = instruction & 0x7
				register[r0] = register[r1] + register[r2]
			}
			update_flags(r0)
		case OP_AND:
			var r0 uint16 = (instruction >> 9) & 0x7
			var r1 uint16 = (instruction >> 6) & 0x7
			var imm_flag uint16 = (instruction >> 5) & 0x1
			if imm_flag == 1 {
				var imm5 uint16 = sign_extent(instruction&0x1F, 5)
				register[r0] = imm5 & register[r1]
			} else {
				var r2 uint16 = (instruction & 0x7)
				register[r0] = register[r1] & register[r2]
			}
			update_flags(r0)
		case OP_NOT:
			var r0 uint16 = (instruction >> 9) & 0x7
			var r1 uint16 = (instruction >> 6) & 0x7
			register[r0] = ^register[r1]
			update_flags(r0)
		case OP_BR:
			condition_flag := (instruction >> 9) & 0x7
			pcoffset := sign_extent(instruction&0x1FF, 9)
			if (condition_flag & register[COND]) != 0 {
				register[PC] = register[PC] + pcoffset
			}
		case OP_JMP:
			var baseR uint16 = (instruction >> 6) & 0x7
			register[PC] = register[baseR]
		case OP_JSR:
			var long_flag uint16 = (instruction >> 11) & 0x1
			register[R7] = register[PC]
			if long_flag == 1 {
				var pcoffset uint16 = sign_extent(instruction&0x7FF, 11)
				register[PC] += pcoffset
			} else {
				var baseR uint16 = (instruction >> 6) & 0x7
				register[PC] = baseR
			}
		case OP_LD:
			var r0 uint16 = (instruction >> 9) & 0x7
			var pcoffset uint16 = sign_extent(instruction&0x1FF, 9)
			register[r0] = mem_read(register[PC] + pcoffset)
			update_flags(r0)
		case OP_LDI:
			var r0 uint16 = (instruction >> 9) & 0x7
			var pcoffset uint16 = sign_extent(instruction&0xFF, 9)
			register[r0] = mem_read(mem_read(register[PC] + pcoffset))
			update_flags(r0)
		case OP_LDR:
			var offset uint16 = sign_extent(instruction&0x3F, 6)
			var baseR uint16 = (instruction >> 6) & 0x7
			var r0 uint16 = (instruction >> 9) & 0x7
			register[r0] = mem_read(register[baseR] + offset)
			update_flags(r0)
		case OP_LEA:
			var r0 uint16 = (instruction >> 9) & 0x7
			var pcoffset uint16 = sign_extent(instruction&0x1FF, 9)
			register[r0] = register[PC] + pcoffset
			update_flags(r0)
		case OP_ST:
			var r0 uint16 = (instruction >> 9) & 0x7
			var pcoffset uint16 = sign_extent(instruction&0x1FF, 9)
			mem_write(register[PC]+pcoffset, register[r0])
		case OP_STI:
			var pcoffset uint16 = sign_extent(instruction&0x1FF, 9)
			var r0 uint16 = (instruction >> 9) & 0x7
			mem_write(mem_read(register[PC]+pcoffset), register[r0])
		case OP_STR:
			var pcoffset uint16 = sign_extent(instruction&0x3F, 6)
			var baseR uint16 = (instruction >> 6) & 0x7
			var r0 uint16 = (instruction >> 9) & 0x7
			mem_write(register[baseR]+pcoffset, register[r0])
		case OP_TRAP:
			register[R7] = register[PC]

			switch instruction & 0xFF {

			case TRAP_GETC:
				char, _, err := keyboard.GetKey()
				if err == nil {
					register[R0] = uint16(char)
					update_flags(R0)
				}
			case TRAP_OUT:
				fmt.Printf("%c", register[R0])
			case TRAP_PUTS:
				for address := register[R0]; memory[address] != 0x00; address++ {
					fmt.Printf("%c", memory[address])
				}
			case TRAP_IN:
				fmt.Print("Enter a character: ")
				char, _, err := keyboard.GetKey()
				if err == nil {
					fmt.Printf("%c", char)
					register[R0] = uint16(char)
					update_flags(R0)
				}
			case TRAP_PUTSP:
				for address := register[R0]; memory[address] != 0x00; address++ {
					value := memory[address]

					fmt.Printf("%c", value&0xff)

					symb := value & 0xff >> 8
					if symb != 0 {
						fmt.Printf("%c", symb)
					}
				}
			case TRAP_HALT:
				fmt.Println("HALT")
			}
		case OP_RES:
			running = 0
			return running
		case OP_RTI:
			running = 0
			return running
		default:
			running = 0
			return running
		}
	}
	return running
}

func main() {

	if len(os.Args) < 2 {
		os.Exit(1)
	}

	err := keyboard.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

	register[PC] = 0x3000

	programFile := os.Args[1]
	ReadImage(programFile)

	running := 1
	instr(running)
}
