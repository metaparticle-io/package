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
	fmt.Fprintf(w, "Hello metaparticle from %s %s!\n", r.RequestURI, hostname)
	fmt.Printf("Request received: %s\n", r.RequestURI)
}

func main() {
	metaparticle.Containerize(
		&metaparticle.Runtime{
			Ports:    []int32{port},
			Executor: "metaparticle",
			Replicas: 3,
		},
		&metaparticle.Package{
			Name:       "metaparticle-web-demo",
			Repository: "docker.io/brendanburns",
			Builder:    "docker",
			Verbose:    true,
			Publish:    true,
		},
		func() {
			log.Printf("Starting server on :%d\n", port)
			http.HandleFunc("/", handler)
			err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
			if err != nil {
				log.Fatal("Couldn't start the server: ", err)
			}
		})
}
