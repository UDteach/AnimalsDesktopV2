//go:build !windows && !darwin

package main

import "fmt"

func main() {
	fmt.Println("Animals Desktop is currently implemented for Windows.")
}
