package main

import (
	"fmt"
	"os"

	"github.com/katcipis/mdtoc"
)

func main() {
	if len(os.Args) == 1 {
		mdtoc.Generate(os.Stdin, os.Stdout)
		return
	}
	if len(os.Args) == 2 {
		mdtoc.GenerateFromFile(os.Args[1], os.Stdout)
		return
	}

	if os.Args[1] != "-w" {
		fmt.Printf("usage: %s -w <file>\n", os.Args[0])
		return
	}

	mdtoc.GenerateInPlace(os.Args[2])
}
