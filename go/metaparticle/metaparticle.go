package metaparticle

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
)

// Runtime represents a place where to run the containers (Docker, rkt, etc)
type Runtime struct {
	// must return an integer >= 1
	Replicas int32

	// must return an integer >= 0
	Shards int32

	// TODO: find out what this does
	URLShardPattern string

	// The name of the executor (e.g.: docker)
	Executor string

	// Returns the ports that the service exposes
	Ports []int32

	// Returns whether the service wants a public IP when deployed (usually involves the creation of a load balancer)
	PublicAddress bool

	// Pointer to extra configuration needed by specific executors (e.g. AciRuntimeConfig)
	ExtraConfig interface{}
}

// Package encapsulates the metadata needed to build the image
type Package struct {
	// Image name
	Name string
	// Image repository (e.g. quay.io/user)
	Repository string
	Verbose    bool
	Quiet      bool
	// Image builder (e.g. docker)
	Builder string
	// Whether to publish the built image to the remote repository
	Publish bool
}

// Executor implementors are container platforms where the containers can be deployed (e.g. Azure, GCP)
type Executor interface {
	Run(image string, name string, config *Runtime, stdout io.Writer, stderr io.Writer) error
	Logs(name string, stdout io.Writer, stderr io.Writer) error
	Cancel(name string) error
}

// Builder is the interface that wraps the container runtime's build and push methods
type Builder interface {
	Build(dir string, image string, stdout io.Writer, stderr io.Writer) error
	Push(image string, stdout io.Writer, stderr io.Writer) error
}

func inDockerContainer() bool {
	if os.Getenv("METAPARTICLE_IN_CONTAINER") == "true" {
		return true
	}
	b, err := ioutil.ReadFile("/proc/1/cgroup")

	if err != nil {
		return false
	}
	s := string(b)
	return strings.Contains(s, "docker") || strings.Contains(s, "kubepods")
}

func executorFromRuntime(r *Runtime) (Executor, error) {
	if r == nil {
		return NewDockerImpl()
	}
	switch r.Executor {
	case "docker":
		return NewDockerImpl()
	case "aci":
		return NewACIExecutor()
	case "metaparticle":
		return &MetaparticleExecutor{}, nil
	default:
		return nil, fmt.Errorf("Unknown executor: %s", r.Executor)
	}
}

func builderFromPackage(pkg *Package) (Builder, error) {
	if pkg == nil {
		return NewDockerImpl()
	}
	switch pkg.Builder {
	case "docker":
		return NewDockerImpl()
	default:
		return nil, fmt.Errorf("Unknown builder: %s", pkg.Builder)
	}
}

func writeDockerfile(name string) error {
	contents := `FROM golang:1.9 as builder
WORKDIR /go/src/app
COPY . .

RUN go get -u github.com/golang/dep/cmd/dep
RUN dep init
RUN go-wrapper install


FROM ubuntu

COPY --from=builder /go/bin/app .

CMD ["./app"]
`
	return ioutil.WriteFile("Dockerfile", []byte(contents), 0644)
}

// Containerize receives a description of the runtime and metadata needed to build a container image,
// and run it.
//
// When called inside the container, it runs the function f. When called outside the container, it builds
// the container image and runs it in the specified runtime environment.
func Containerize(r *Runtime, p *Package, f func()) {
	if inDockerContainer() {
		f()
	} else {
		exec, err := executorFromRuntime(r)

		if err != nil {
			panic(fmt.Sprintf("Could not get an executor for this runtime: %v", err))
		}

		builder, err := builderFromPackage(p)
		if err != nil {
			panic(fmt.Sprintf("Could not get a builder for this package: %v", err))
		}

		image := p.Name

		if len(p.Repository) != 0 {
			image = p.Repository + "/" + image
		}
		err = writeDockerfile("metaparticle-package")
		if err != nil {
			panic(fmt.Sprintf("Could not write Dockerfile: %v", err))
		}

		err = builder.Build(".", image, os.Stdout, os.Stderr)
		if err != nil {
			panic(fmt.Sprintf("Could not build the container: %v", err))
		}

		if p.Publish {
			err = builder.Push(image, os.Stdout, os.Stderr)
			if err != nil {
				panic(fmt.Sprintf("Could not push image \"%v\" to the repository: %v", image, err))
			}
		}

		err = exec.Run(image, p.Name, r, os.Stdout, os.Stderr)
		if err != nil {
			panic(fmt.Sprintf("Error running the container: %v", err))
		}

		go func() {
			exec.Logs(p.Name, os.Stdout, os.Stderr)
		}()

		signalChan := make(chan os.Signal, 1)
		cleanupDone := make(chan bool)
		signal.Notify(signalChan, os.Interrupt)
		go func() {
			for _ = range signalChan {
				fmt.Println("Received interrupt, stopping container...")
				exec.Cancel(p.Name)
				cleanupDone <- true
			}
		}()
		<-cleanupDone
	}
}
