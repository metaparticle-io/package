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
            Executor:        "docker"},
        &metaparticle.Package{Repository: "xfernando",
            Builder: "docker"},
        func() {
            fmt.Println("Hello World")
        })
}
```

Then you only have to do `go run main.go`, and a container will be built and deployed to your docker instance.