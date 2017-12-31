package metaparticle

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/pkg/errors"
)

type mockReadCloser struct {
	*bytes.Buffer
}

func (m *mockReadCloser) Close() error {
	return nil
}

type mockImageClient struct {
	serverResponse string
	serverError    error
}

func (m *mockImageClient) ImageBuild(ctx context.Context, context io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	if m.serverError != nil {
		return types.ImageBuildResponse{}, m.serverError
	}
	body := &mockReadCloser{bytes.NewBufferString(m.serverResponse)}
	return types.ImageBuildResponse{Body: body}, nil
}

func (m *mockImageClient) ImagePush(ctx context.Context, ref string, options types.ImagePushOptions) (io.ReadCloser, error) {
	if m.serverError != nil {
		return nil, m.serverError
	}

	body := &mockReadCloser{bytes.NewBufferString(m.serverResponse)}
	return body, nil
}

func TestDockerBuild(t *testing.T) {
	serverErr := fmt.Errorf("Error")
	cases := []struct {
		name              string
		dir               string
		image             string
		serverResponse    string
		serverError       error
		expectedError     error
		expectedErrorType reflect.Type
	}{
		{
			name:          "Test docker build - empty context directory",
			image:         "test",
			expectedError: errEmptyContextDirectory,
		},
		{
			name:          "Test docker build - empty image name",
			dir:           ".",
			expectedError: errEmptyImageName,
		},
		{
			name:           "Test docker build - project build",
			dir:            "./test-data/test-project",
			image:          "test",
			serverResponse: `{"stream":"Build successful"}`,
		},
		{
			name:          "Test docker build - inexistent project directory",
			dir:           "./test-data/aknADKASj",
			image:         "test",
			expectedError: errInexistentContextDirectory,
		},
		{
			name:              "Test garbled server response",
			dir:               "./test-data/test-project",
			image:             "test",
			serverResponse:    `f9q*sa9das}`,
			expectedErrorType: reflect.TypeOf(&json.SyntaxError{}),
		},
		{
			name:          "Test server error",
			dir:           "./test-data/test-project",
			image:         "test",
			serverError:   serverErr,
			expectedError: serverErr,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDockerImpl, _ := newDockerImpl(&mockImageClient{c.serverResponse, c.serverError}, nil)
			err := mockDockerImpl.Build(c.dir, c.image, os.Stdout, os.Stderr)
			cause := errors.Cause(err)
			if c.expectedErrorType != nil {
				if reflect.TypeOf(cause) != c.expectedErrorType {
					t.Errorf("Expected error of type %v, got %v", reflect.TypeOf(cause), c.expectedErrorType)
				}
			} else if cause != c.expectedError {
				t.Errorf("Expected %v error, got %v", c.expectedError, err)
			}
		})
	}
}

func TestDockerPush(t *testing.T) {
	serverErr := fmt.Errorf("Error")
	cases := []struct {
		name              string
		image             string
		serverResponse    string
		serverError       error
		expectedError     error
		expectedErrorType reflect.Type
	}{
		{
			name:          "Test empty image name",
			expectedError: errEmptyImageName,
		},
		{
			name:           "Test image push",
			image:          "test",
			serverResponse: `{"stream":"Push successful"}`,
		},
		{
			name:              "Test garbled server response",
			image:             "test",
			serverResponse:    `f9q*sa9das}`,
			expectedErrorType: reflect.TypeOf(&json.SyntaxError{}),
		},
		{
			name:          "Test server error",
			image:         "test",
			serverError:   serverErr,
			expectedError: serverErr,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDockerImpl, _ := newDockerImpl(&mockImageClient{c.serverResponse, c.serverError}, nil)
			err := mockDockerImpl.Push(c.image, os.Stdout, os.Stderr)
			cause := errors.Cause(err)
			if c.expectedErrorType != nil {
				if reflect.TypeOf(cause) != c.expectedErrorType {
					t.Errorf("Expected error of type %v, got %v", reflect.TypeOf(cause), c.expectedErrorType)
				}
			} else if cause != c.expectedError {
				t.Errorf("Expected %v error, got %v", c.expectedError, err)
			}
		})
	}
}

type mockContainerRunner struct {
	createError error
	startError  error
	stopError   error
	removeError error
	logsError   error
}

func (m *mockContainerRunner) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (container.ContainerCreateCreatedBody, error) {
	return container.ContainerCreateCreatedBody{}, m.createError
}
func (m *mockContainerRunner) ContainerStart(ctx context.Context, container string, options types.ContainerStartOptions) error {
	return m.startError
}
func (m *mockContainerRunner) ContainerStop(ctx context.Context, container string, timeout *time.Duration) error {
	return m.stopError
}
func (m *mockContainerRunner) ContainerRemove(ctx context.Context, container string, options types.ContainerRemoveOptions) error {
	return m.removeError
}
func (m *mockContainerRunner) ContainerLogs(ctx context.Context, container string, options types.ContainerLogsOptions) (io.ReadCloser, error) {
	// Log uses StdCopy.stdcopy, which tests the first byte to see what is the type of the stream (stdcopy.Stdin = 0)
	return &mockReadCloser{bytes.NewBufferString("\u0000" + "Hello there")}, m.logsError
}
func TestDockerRun(t *testing.T) {
	serverError := fmt.Errorf("Error")
	cases := []struct {
		name          string
		image         string
		containerName string
		config        *Runtime
		expectedError error
		createError   error
		startError    error
	}{
		{
			name:          "Test docker run - container start",
			image:         "test",
			containerName: "test",
			config:        &Runtime{},
		},
		{
			name:          "Test docker run - container create error",
			image:         "test",
			containerName: "test",
			config:        &Runtime{},
			createError:   serverError,
			expectedError: serverError,
		},
		{
			name:          "Test docker run - container start error",
			image:         "test",
			containerName: "test",
			config:        &Runtime{},
			startError:    serverError,
			expectedError: serverError,
		},
		{
			name:          "Test docker run - container with exposed port",
			image:         "test",
			containerName: "test",
			config:        &Runtime{Ports: []int32{80}},
		},
		{
			name:          "Test docker run - nil runtime",
			image:         "test",
			containerName: "test",
			config:        nil,
			expectedError: errNilRuntimeConfig,
		},
		{
			name:          "Test docker run - empty image name",
			containerName: "test",
			config:        &Runtime{},
			expectedError: errEmptyImageName,
		},
		{
			name:          "Test docker run - empty container name",
			image:         "test",
			config:        &Runtime{},
			expectedError: errEmptyContainerName,
		},
		{
			name:          "Test docker run - port below range",
			image:         "test",
			containerName: "test",
			config:        &Runtime{Ports: []int32{-1}},
			expectedError: errInvalidPort(-1),
		},
		{
			name:          "Test docker run - port above range",
			image:         "test",
			containerName: "test",
			config:        &Runtime{Ports: []int32{70000}},
			expectedError: errInvalidPort(70000),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDockerImpl, _ := newDockerImpl(nil, &mockContainerRunner{createError: c.createError, startError: c.startError})
			err := mockDockerImpl.Run(c.image, c.containerName, c.config, os.Stdout, os.Stderr)
			if errors.Cause(err) != c.expectedError {
				t.Errorf("Expected %v error, got %v", c.expectedError, err)
			}
		})
	}
}

func TestDockerLogs(t *testing.T) {
	serverError := fmt.Errorf("Error")
	cases := []struct {
		name          string
		containerName string
		expectedError error
		logsError     error
	}{
		{
			name:          "Test docker logs - success",
			containerName: "test",
		},
		{
			name:          "Test docker logs - empty container name",
			expectedError: errEmptyContainerName,
		},
		{
			name:          "Test docker logs - server error",
			containerName: "test",
			logsError:     serverError,
			expectedError: serverError,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDockerImpl, _ := newDockerImpl(nil, &mockContainerRunner{logsError: c.logsError})
			err := mockDockerImpl.Logs(c.containerName, os.Stderr, os.Stdout)
			if errors.Cause(err) != c.expectedError {
				t.Errorf("Expected %v error, got %v", c.expectedError, err)
			}
		})
	}
}

func TestDockerCancel(t *testing.T) {
	serverError := fmt.Errorf("Error")
	cases := []struct {
		name          string
		containerName string
		expectedError error
		stopError     error
		removeError   error
	}{
		{
			name:          "Test docker cancel - success",
			containerName: "test",
		},
		{
			name:          "Test docker cancel - empty container name",
			expectedError: errEmptyContainerName,
		},
		{
			name:          "Test docker cancel - stop error",
			containerName: "test",
			stopError:     serverError,
			expectedError: serverError,
		},
		{
			name:          "Test docker cancel - remove error",
			containerName: "test",
			removeError:   serverError,
			expectedError: serverError,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDockerImpl, _ := newDockerImpl(nil, &mockContainerRunner{stopError: c.stopError, removeError: c.removeError})
			err := mockDockerImpl.Cancel(c.containerName)
			if errors.Cause(err) != c.expectedError {
				t.Errorf("Expected %v error, got %v", c.expectedError, err)
			}
		})
	}
}
