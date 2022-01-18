package util

import (
	"math/rand"
)

const alphanumeric = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// use uppercase and lowercase because there will be cases where we are receiving hex strings from inputs that are mixed case (such as Ethereum addresses)
const hex = "0123456789abcdefABCDEF"

// RandStringBytes returns a random alphanumeric string of given length
func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphanumeric[rand.Intn(len(alphanumeric))]
	}
	return string(b)
}

// RandHexString returns a random hex string of given length
func RandHexString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = hex[rand.Intn(len(hex))]
	}
	return string(b)
}
