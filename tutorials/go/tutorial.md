# Metaparticle/Package for Golang Tutorial
This is an in-depth tutorial for using Metaparticle/Package for Go.

For a quick summary, please see the [README](README.md).

## Initial Setup

### Check the tools
The `docker` command line tool needs to be installed and working. Try:
`docker ps` to verify this.  Go to the [install page](https://get.docker.io) if you need
to install Docker.

The `mp-compiler` command line tool needs to be installed and working.
Try `mp-compiler --help` to verify this. Go to the [releases page](https://github.com/metaparticle-io/metaparticle-ast/releases) if you need to install
the Metaparticle compiler.

### Get the code
```sh
$ git clone https://github.com/metaparticle-io/package
$ cd tutorials/go/
# [optional, substitute your favorite editor here...]
$ code .
```

## Initial Program
Inside of the `tutorials/go` directory, you will find a simple go web server.

You can build this project with `go build main.go`.

The initial code is a very simple "Hello World"

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var port int32 = 8080

func handler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	fmt.Fprintf(w, "Hello metaparticle from %s %s!\n", r.RequestURI, hostname)
	fmt.Printf("Request received: %s\n", r.RequestURI)
}

func main() {
     log.Println("Starting server on :8080")
     http.HandleFunc("/", handler)
     if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
     	log.Fatal("Couldn't start the server: ", err)
     }
}
```

## Step One: Containerize the Application
To build a container from our simple application we need to add a dependency to our
build file, and then update the code.

Run:
```sh
$ go get https://github.com/metaparticle-io/package/go/metaparticle
```

Then update the code to read as follows:

```go
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
		&metaparticle.Package{
			Name:       "metaparticle-web-demo",
			Repository: "your-docker-user-goes-here",
			Builder:    "docker",
			Verbose:    true,
			Publish:    true,
		},
		func() {
			log.Println("Starting server on :8080")
			http.HandleFunc("/", handler)
			err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
			if err != nil {
				log.Fatal("Couldn't start the server: ", err)
			}
		})
}
```

You will notice that we added a `&metaparticle.Package` struct that describes how
to package the application. You will need to replace `your-docker-user-goes-here`
with an actual Docker repository path.

You will also notice that we wrapped the main function in the `metaparticle.Containerize`
function which kicks off the Metaparticle code.

```sh
go run main.go
```

This code will start your web server again. But this time, it is running
inside a container. You can see this by running:

```sh
docker ps
```

## Step Two: Exposing the ports
If you try to access the web server on [http://localhost:8080](http://localhost:8080) you
will see that you can not actually access the server. Despite it running, the service
is not exposed. To do this, you need to add a `&metaparticle.Runtime` struct to supply the
port(s) to expose.

The code snippet to add is:

```go
...
		&metaparticle.Runtime{
			Ports:    []int32{port},
			Executor: "docker",
		},
...
```

This tells the runtime the port(s) to expose. The complete code looks like:

```go
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
			Executor: "docker",
		},
		&metaparticle.Package{
			Name:       "metaparticle-web-demo",
			Repository: "your-docker-user-goes-here",
			Builder:    "docker",
			Verbose:    true,
			Publish:    true,
		},
		func() {
			log.Println("Starting server on :8080")
			http.HandleFunc("/", handler)
			err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
			if err != nil {
				log.Fatal("Couldn't start the server: ", err)
			}
		})
}
```

Now if you run this with `dotnet run` your webserver will be successfully exposed on port 8080.

You can verify that it works by running `curl localhost:8080`

## Replicating and exposing on the web.
As a final step, consider the task of exposing a replicated service on the internet.
To do this, we're going to expand our usage of the `&metaparticle.Runtime` tag. First we will
add a `replicas` field, which will specify the number of replicas. Second we will
set our execution environment to `metaparticle` which will launch the service
into the currently configured Kubernetes environment.

Here's what the snippet looks like:

```go
...
		&metaparticle.Runtime{
			Ports:    []int32{port},
			Executor: "metaparticle",
			Replicas: 3,
		},
...
```

And the complete code looks like:
```go
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
			Repository: "brendanburns",
			Builder:    "docker",
			Verbose:    true,
			Publish:    true,
		},
		func() {
			log.Println("Starting server on :8080")
			http.HandleFunc("/", handler)
			err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
			if err != nil {
				log.Fatal("Couldn't start the server: ", err)
			}
		})
}
```

You can run this using:
```sh
go run main.go
```

After you compile and run this, you can see that there are four replicas running behind a
Kubernetes Service Load balancer:

```sh
$ kubectl get pods
...
$ kubectl get services
...
```

// Still looking for more? Continue on to the more advanced [sharding tutorial](sharding-tutorial.md)
