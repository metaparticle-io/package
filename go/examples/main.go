package main

import (
	"fmt"

	"github.com/xfernando/package/go/metaparticle"
)

var port int32 = 8080

func main() {
	metaparticle.Containerize(
		&metaparticle.Runtime{
			Ports:           []int32{port},
			Shards:          0,
			URLShardPattern: "^\\/users\\/([^\\/]*)\\/.*",
			Executor:        "docker"},
		&metaparticle.Package{Repository: "xfernando",
			Builder: "docker",
			Verbose: true},
		func() {
			fmt.Println("Hello there")
		})
}
