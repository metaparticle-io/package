# Metaparticle/Package for Go

Metaparticle/Package is a collection of libraries intended to
make building and deploying containers a seamless and idiomatic
experience for developers.

This is the implementation for Go.

## Introduction

Metaparticle/Package simplifies and centralizes the task of
building and deploying a container image.

Here is a quick example.

Consider this simple Go application:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello World")
}
```

To containerize this application, you need to replace the contents of your main on the main.go file (the entrypoint is still hardcoded)
function with the `Containerize` wrapper, passing the `Runtime` and `Package` structs, and wrapping your old main on a `func()`, like this:

```go
package main

import(
    "fmt"
    "github.com/metaparticle-io/package/go/metaparticle"
)

func main() {
    metaparticle.Containerize(
        &metaparticle.Runtime{
            Executor:   "docker",
        },
        &metaparticle.Package{
            Repository: "docker.io/brendanburns",
            Builder:    "docker"
        },
        func() {
            fmt.Println("Hello World")
        })
}
```

Then you only have to do `go run main.go`, and a container will be built and deployed to your docker instance.

## Registry authentication

If you want your image pushed to your registry, there are two things needed.

First you need to add `Publish: true` to the `metaparticle.Package` annotation. Changing the previous example we get:

```go
package main

import(
    "fmt"
    "github.com/metaparticle-io/package/go/metaparticle"
)

func main() {
    metaparticle.Containerize(
        &metaparticle.Runtime{
            Executor:   "docker",
        },
        &metaparticle.Package{
            Repository: "docker.io/brendanburns",
            Builder:    "docker",
            Publish:    true,
        },
        func() {
            fmt.Println("Hello World")
        })
}
```

Then you need to set two environment variables, `MP_REGISTRY_USERNAME` with your registry user, and `MP_REGISTRY_PASSWORD` with
your registry password. Or you can pass them in the shell when starting the program:

```bash
env MP_REGISTRY_USERNAME=youruser MP_REGISTRY_PASSWORD=yourpassword go run main.go
```

These are used to generate the authorization string passed to the docker server. When you execute `go run main.go`
with both set, they'll be used to authenticate you to the remote registry.

The name of the repository must be a canonical name (e.g. docker.io/brendanburns). So, if you're using docker hub
you can't omit the docker.io part of the repository name.