package main

import (
	"fmt"

	"github.com/metaparticle-io/package/go/metaparticle"
)

func main() {
	metaparticle.Containerize(
		&metaparticle.Runtime{
			Executor: "docker"},
		&metaparticle.Package{Repository: "xfernando",
			Builder: "docker"},
		func() {
			fmt.Println("Hello World")
		})
}
