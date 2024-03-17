package utils

import (
	"strconv"
)

func IntToBytes(n int) []byte {
	// Create a byte array with the same size as the int
	bytes := make([]byte, 4)

	// Convert the int to bytes
	for i := uint(0); i < 4; i++ {
		// Shift the bits of the int to extract each byte
		bytes[i] = byte(n >> (i * 8))
	}

	return bytes
}

func BytesToString(bytes []byte) string {
	var strBytes string
	for _, b := range bytes {
		strBytes += strconv.Itoa(int(b))
	}

	return strBytes
}
