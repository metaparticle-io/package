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

If you want your image pushed to a registry, two things are needed.

First you need to add `Publish: true` to the `metaparticle.Package` annotation.

Then you need to set two environment variables, `MP_REGISTRY_USERNAME` with your registry user, and `MP_REGISTRY_PASSWORD` with
your registry password.

These are used to generate the authorization string passed to the docker server. When you execute `go run main.go`
with both set, they'll be used to authenticate you to the remote registry.

Important:exclamation: The name of the repository must be a canonical name (e.g. docker.io/brendanburns). So, if you're using docker hub
you can't omit the docker.io part of the repository name.

### Docker Hub

Changing the annotation on the previous example we get:

```go
        &metaparticle.Package{
            Repository: "docker.io/brendanburns",
            Builder:    "docker",
            Publish:    true,
        },
```

After this change, you can run the command below and your container will be pushed to docker hub:

```bash
env MP_REGISTRY_USERNAME=youruser MP_REGISTRY_PASSWORD=yourpassword go run main.go
```

### Google Container Registry

Registries on GCR are scoped to projects, so with your `[[PROJECT_ID]]` in hand, you need to put `gcr.io/PROJECT_ID` on the `Repository` string:

```go
        &metaparticle.Package{
            Repository: "gcr.io/[[PROJECT_ID]]",
            Builder:    "docker",
            Publish:    true,
        },
```

Then, you can use google cloud's application-default token to authenticate agains GCR as such:

```bash
env MP_REGISTRY_USERNAME=oauth2accesstoken MP_REGISTRY_PASSWORD=$(gcloud auth application-default print-access-token) go run main.go
```

More details about GCR authentication can be found [here](https://cloud.google.com/container-registry/docs/advanced-authentication).

### Azure Container Registry

Assuming you have a service principal configured (see docs [here](https://docs.microsoft.com/en-us/azure/azure-stack/azure-stack-create-service-principals))
with the contributor role in your ACR you just need to change the Repository string to `yourregistry.azurecr.io`:

```go
        &metaparticle.Package{
            Repository: "yourregistry.azurecr.io",
            Builder:    "docker",
            Publish:    true,
        },
```

And then you run:

```bash
env MP_REGISTRY_USERNAME=$APPLICATION_ID MP_REGISTRY_PASSWORD=$APPLICATION_KEY go run main.go
```

Note that `$APPLICATION_ID` is not the name you chose for the service principal, but the `Application ID` that's automatically
generated when you created the service principal. It can be found by viewing the properties of your service principal on the Azure Portal
by going to Azure Active Directory -> App Registrations.

And `$APPLICATION_KEY` is an authentication key you created for the service principal.