package main

import (
	"bytes"
	"testing"
)

func TestBom(t *testing.T) {
	var tests = []struct {
		input []byte
		want  BOM
	}{
		{[]byte{}, NONE},
		{[]byte{0x00}, NONE},
		{[]byte{0xfe}, NONE},
		{[]byte{0xff}, NONE},
		{[]byte{0xff, 0xfe}, UTF16_LE},
		{[]byte{0xfe, 0xff}, UTF16_BE},
		{[]byte{0xef, 0xbb, 0xbf}, UTF8},
		{[]byte{0xff, 0xfe, 0x00, 0x00}, UTF32_LE},
		{[]byte{0x00, 0x00, 0xfe, 0xff}, UTF32_BE},
	}
	for _, test := range tests {
		got, err := bom(bytes.NewBuffer(test.input))
		if err != nil {
			t.Fatal(err)
		}
		if got != test.want {
			t.Errorf("bom(%q) = %v", test.input, got)
		}
	}
}
