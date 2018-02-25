package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/metaparticle-io/package/go/metaparticle"
)

var port int32 = 80

func handler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	fmt.Fprintf(w, "Hello metaparticle in Azure Container Instances from %s %s!\n", r.RequestURI, hostname)
	fmt.Printf("Request received: %s\n", r.RequestURI)
}

func main() {
	metaparticle.Containerize(
		&metaparticle.Runtime{
			Executor:      "aci",
			Ports:         []int32{port},
			PublicAddress: true,
		},
		&metaparticle.Package{
			Name:       "metaparticle-aci-demo",
			Repository: "docker.io/radumatei",
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
