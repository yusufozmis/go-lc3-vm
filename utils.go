package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"

	"github.com/eiannone/keyboard"
)

func sign_extent(x uint16, bit_count int) uint16 {
	if x>>(bit_count-1)&0x1 == 1 {
		x |= 0xFFFF << bit_count
	}
	return x
}

func update_flags(r uint16) {
	if register[r] == 0 {
		register[COND] = FL_ZRO
	} else if (register[r] >> 15) == 1 {
		register[COND] = FL_NEG
	} else {
		register[COND] = FL_POS
	}
}

var NativeEndian binary.ByteOrder

func swap16(word uint16) uint16 {
	if NativeEndian == binary.BigEndian {
		return (word << 8) | (word >> 8)
	} else {
		return word
	}
}
func check_key() uint16 {
	if char, key, err := keyboard.GetKey(); err == nil {
		if key == keyboard.KeyEsc || key == keyboard.KeyCtrlC {
			log.Println("Pressed escape")
			os.Exit(0)
		}
		return uint16(char)
	} else {
		log.Printf("Error: %s", err)
	}
	return 0
}

func readObjFile(path string) ([]byte, int64) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Error while opening file", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		log.Fatal("file.Stat failed", err)
	}

	var size int64 = info.Size()

	data := make([]byte, size)

	_, err = file.Read(data)
	if err != nil {
		log.Fatal(err)
	}
	return data, size
}

func ReadImage(path string) {
	var header uint16

	data, _ := readObjFile(path)

	buffer := bytes.NewBuffer(data)
	header = swap16(binary.BigEndian.Uint16(buffer.Next(2)))

	bufferLen := buffer.Len()
	origin := header

	for i := 0; i < bufferLen; i++ {
		b := buffer.Next(2)
		if len(b) == 0 {
			break
		}
		memory[origin] = swap16(binary.BigEndian.Uint16(b))
		origin++
	}
}
