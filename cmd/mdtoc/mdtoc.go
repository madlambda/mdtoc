package main

import (
	"fmt"
	"log"
	"os"

	"github.com/madlambda/mdtoc"
)

func main() {
	if len(os.Args) == 1 {
		exiterr(mdtoc.Generate(os.Stdin, os.Stdout))
		return
	}
	if len(os.Args) == 2 {
		exiterr(mdtoc.GenerateFromFile(os.Args[1], os.Stdout))
		return
	}

	if os.Args[1] != "-w" {
		fmt.Printf("usage: %s -w <file>\n", os.Args[0])
		return
	}

	exiterr(mdtoc.GenerateInPlace(os.Args[2]))
}

func exiterr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
