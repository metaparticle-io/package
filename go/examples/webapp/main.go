package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/metaparticle-io/package/go/metaparticle"
)

var port int32 = 8080

func handler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	fmt.Fprintf(w, "Hello containers [%s] from %s", r.RequestURI, hostname)
}

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
			log.Println("Starting server on :8080")
			http.HandleFunc("/", handler)
			err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
			if err != nil {
				log.Fatal("Couldn't start the server: ", err)
			}
		})
}
