package helpers

import (
	"crypto/rand"
	"fmt"
)

func GenerateOTP() (string, error) {
	var buf [4]byte

	if _, err := rand.Read(buf[:]); err != nil {
		return "", err
	}

	n := uint32(buf[0]) |
		uint32(buf[1])<<8 |
		uint32(buf[2])<<16 |
		uint32(buf[3])<<24

	return fmt.Sprintf("%06d", n%1000000), nil
}
