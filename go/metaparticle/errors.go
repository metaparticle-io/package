package metaparticle

import (
	"errors"
	"fmt"
)

type errInvalidPort int32

func (e errInvalidPort) Error() string {
	return fmt.Sprintf("Invalid port number %v", int32(e))
}

var (
	errEmptyContainerName         = errors.New("A container name must be specified")
	errEmptyImageName             = errors.New("An image must be specified")
	errNilRuntimeConfig           = errors.New("Runtime config must not be nil")
	errEmptyContextDirectory      = errors.New("A context directory must be specified")
	errInexistentContextDirectory = errors.New("Context directory does not exist")
)
