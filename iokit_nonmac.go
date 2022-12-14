//go:build !iokit
// +build !iokit

// just define the symbol for all other non-native platforms
package main

import "fmt"

func GetSystemInfo() {
	fmt.Println("ERROR: not implemented")
}
