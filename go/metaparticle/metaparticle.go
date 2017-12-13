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

	// The ACI runtime config
	AciConfig *AciRuntimeConfig
}

// AciRuntimeConfig - runtime config for ACI executor.
type AciRuntimeConfig struct {
	// Azure tenant id.
	AzureTenantID string

	// Azure subscription id.
	// Container instances will be created using the given subscription.
	AzureSubscriptionID string

	// Azure client id (aka. application id).
	AzureClientID string

	// Azure client secret (aka. application key).
	// If specified, then will authenticate to Azure as service principal.
	// And then no need to specify username/password.
	AzureClientSecret string

	// If client secret is not specified, use username and password.
	Username string

	// If client secret is not specified, use username and password.
	Password string

	// The resource group to create container instance.
	AciResourceGroup string
}

type Package struct {
	Repository string
	Verbose    bool
	Quiet      bool
	Builder    string
	Publish    bool
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
		return &DockerImpl{}, nil
	}
	switch r.Executor {
	case "docker":
		return &DockerImpl{}, nil
	case "aci":
		// TODO: find a way to parameterize the name of the resource group
		return &ACIExecutor{"test"}, nil
	case "metaparticle":
		return &MetaparticleExecutor{}, nil
	default:
		return nil, fmt.Errorf("Unknown executor: %s", r.Executor)
	}
}

func builderFromPackage(pkg *Package) (Builder, error) {
	if pkg == nil {
		return &DockerImpl{}, nil
	}
	switch pkg.Builder {
	case "docker":
		return &DockerImpl{}, nil
	default:
		return nil, fmt.Errorf("Unknown builder: %s", pkg.Builder)
	}
}

func writeDockerfile(name string) error {
	contents := `FROM golang:1.9
WORKDIR /go/src/app
COPY . .

RUN go-wrapper download
RUN go-wrapper install
	
CMD ["go-wrapper", "run"]
`
	return ioutil.WriteFile("Dockerfile", []byte(contents), 0644)
}

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

		name := "web"
		image := "test"

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
				panic(fmt.Sprintf("Could not push the image to the repository: %v", err))
			}
		}

		err = exec.Run(image, name, r, os.Stdout, os.Stderr)
		if err != nil {
			panic(fmt.Sprintf("Error executing the container: %v", err))
		}

		go func() {
			exec.Logs(name, os.Stdout, os.Stderr)
		}()

		signalChan := make(chan os.Signal, 1)
		cleanupDone := make(chan bool)
		signal.Notify(signalChan, os.Interrupt)
		go func() {
			for _ = range signalChan {
				fmt.Println("Received interrupt, stopping container...")
				exec.Cancel(name)
				cleanupDone <- true
			}
		}()
		<-cleanupDone
	}
}
