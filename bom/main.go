package main

import (
	"fmt"
	"io"
	"os"
)

type BOM int

const (
	NONE BOM = iota
	UTF8
	UTF16_LE
	UTF16_BE
	UTF32_LE
	UTF32_BE
)

func main() {
	bom, err := bom(os.Stdin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if bom != NONE {
		fmt.Println("BOM")
	}
}

func bom(r io.Reader) (BOM, error) {
	bytes := make([]byte, 4)
	n, err := r.Read(bytes)
	if err != nil && err != io.EOF {
		return NONE, err
	}
	switch {
	case isUTF32LE(bytes, n):
		return UTF32_LE, nil
	case isUTF32BE(bytes, n):
		return UTF32_BE, nil
	case isUTF8(bytes, n):
		return UTF8, nil
	case isUTF16LE(bytes, n):
		return UTF16_LE, nil
	case isUTF16BE(bytes, n):
		return UTF16_BE, nil
	default:
		return NONE, nil
	}
}

func isUTF32LE(p []byte, len int) bool {
	return len >= 4 && p[0] == 0xff && p[1] == 0xfe && p[2] == 0x00 && p[3] == 0x00
}

func isUTF32BE(p []byte, len int) bool {
	return len >= 4 && p[0] == 0x00 && p[1] == 0x00 && p[2] == 0xfe && p[3] == 0xff
}

func isUTF8(p []byte, len int) bool {
	return len >= 3 && p[0] == 0xef && p[1] == 0xbb && p[2] == 0xbf
}

func isUTF16LE(p []byte, len int) bool {
	return len >= 2 && p[0] == 0xff && p[1] == 0xfe
}

func isUTF16BE(p []byte, len int) bool {
	return len >= 2 && p[0] == 0xfe && p[1] == 0xff
}
