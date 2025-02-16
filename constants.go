package main

const MEMORY_MAX = 1 << 16

// Registers
const (
	R0 = iota
	R1
	R2
	R3
	R4
	R5
	R6
	R7
	PC //program counter
	COND
	COUNT
)

// Instruction set
const (
	OP_BR   uint16 = iota
	OP_ADD         //add
	OP_LD          //load
	OP_ST          //store
	OP_JSR         //jump register
	OP_AND         //bitwise and
	OP_LDR         //load register
	OP_STR         //store register
	OP_RTI         //unused
	OP_NOT         //bitwise not
	OP_LDI         //load indirect
	OP_STI         //store indirect
	OP_JMP         //jump
	OP_RES         //reserved(unused)
	OP_LEA         //load effective address
	OP_TRAP        //execute trap
)

// Condition Flags
const (
	FL_POS = 1 << 0 //Positive
	FL_ZRO = 1 << 1 //Zero
	FL_NEG = 1 << 2 //Negative
)

const (
	TRAP_GETC  = 0x20 + iota //get character from keyboard, not echoed onto terminal
	TRAP_OUT                 //output a character
	TRAP_PUTS                //output a word string
	TRAP_IN                  //get a charactter from keyboard, echoed onto terminal
	TRAP_PUTSP               //output a byte string
	TRAP_HALT                //halt the program
)

const (
	MR_KBSR = 0xFE00
	MR_KBDR = 0xFE02
)

var memory [MEMORY_MAX]uint16

var register [COUNT]uint16
