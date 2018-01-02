package metaparticle

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
)

// DockerImpl is a docker implementation of the Builder and Executor interfaces
type DockerImpl struct {
	imageClient     dockerImageClient
	containerRunner dockerContainerRunner
}

// Docker's client.ImageAPIClient is too big, so I created this interface with the only methods needed from it.
// This will make mocking code for tests much smaller
type dockerImageClient interface {
	ImageBuild(ctx context.Context, context io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error)
	ImagePush(ctx context.Context, ref string, options types.ImagePushOptions) (io.ReadCloser, error)
}

// Docker's client.ContainerAPIClient is too big so I created this interface with the only methods needed from it
// This will make mocking code for tests much smaller
type dockerContainerRunner interface {
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (container.ContainerCreateCreatedBody, error)
	ContainerStart(ctx context.Context, container string, options types.ContainerStartOptions) error
	ContainerStop(ctx context.Context, container string, timeout *time.Duration) error
	ContainerRemove(ctx context.Context, container string, options types.ContainerRemoveOptions) error
	ContainerLogs(ctx context.Context, container string, options types.ContainerLogsOptions) (io.ReadCloser, error)
}

func newDockerImpl(imageClient dockerImageClient, containerRunner dockerContainerRunner) (*DockerImpl, error) {
	if imageClient == nil && containerRunner == nil {
		dockerClient, err := client.NewEnvClient()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create docker client")
		}

		// dockerClient implements both APIs, but having separate struct members with specific interfaces
		// will make mocking for testing easier
		return &DockerImpl{dockerClient, dockerClient}, nil
	}

	return &DockerImpl{imageClient, containerRunner}, nil
}

// NewDockerImpl returns a singleton struct that uses docker to implement metaparticle.Builder and metaparticle.Executor.
//
// It uses the environment variables DOCKER_CERT_PATH, DOCKER_HOST, DOCKER_API_VERSION and DOCKER_TLS_VERIFY
// to instantiate instantiate a docker API client.
// When these variables are not specified, it defaults to the client running on the local machine.
func NewDockerImpl() (*DockerImpl, error) {
	return newDockerImpl(nil, nil)
}

// createTarGz creates a tarball of the directory in order to send it to the docker server. It returns
// the full path of the file created and an error
func createTarGz(dir string) (string, error) {
	tarball, err := ioutil.TempFile("", "context.tar.gz")
	if err != nil {
		return "", errors.Wrap(err, "Failed to create temporary project tarball")
	}
	defer tarball.Close()

	gw := gzip.NewWriter(tarball)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	err = filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			var filePath string
			if path != "." {
				filePath = path
			} else {
				filePath = info.Name()
			}
			header, err := tar.FileInfoHeader(info, filePath)
			if err != nil {
				return errors.Wrap(err, "Error creating a file header")
			}

			if err := tw.WriteHeader(header); err != nil {
				return errors.Wrap(err, "Error writing file header")
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return errors.Wrap(err, "Error opening temporary tarball file")
			}
			defer file.Close()
			_, err = io.Copy(tw, file)
			return err
		})

	if err != nil {
		return "", errors.Wrap(err, "Error creating temporary tarball with project's contents")
	}

	return tarball.Name(), nil
}

// printStreamResponse decodes a json stream response from the docker server
func printStreamResponse(body io.ReadCloser, out io.Writer) error {
	var line struct {
		Stream string `json:"stream"`
	}
	decoder := json.NewDecoder(body)
	for {
		err := decoder.Decode(&line)
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		out.Write([]byte(line.Stream))
	}
	return nil
}

// Build creates a tarball with the directory's contents and sends it to docker to be build the image
func (d *DockerImpl) Build(dir string, image string, stdout io.Writer, stderr io.Writer) error {
	if len(dir) == 0 {
		return errEmptyContextDirectory
	}
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return errInexistentContextDirectory
		}
		return err
	}

	if len(image) == 0 {
		return errEmptyImageName
	}

	name, err := createTarGz(dir)
	if err != nil {
		return errors.Wrap(err, "Error creating temporaty tarball")
	}
	tarfile, err := os.Open(name)
	if err != nil {
		return errors.Wrap(err, "Error opening the temporary tarball")
	}
	defer os.Remove(name)
	defer tarfile.Close()
	ctx := context.Background()
	res, err := d.imageClient.ImageBuild(ctx, tarfile, types.ImageBuildOptions{Tags: []string{image}})
	if err != nil {
		return errors.Wrap(err, "Error sending build request to docker")
	}

	if err = printStreamResponse(res.Body, stdout); err != nil {
		return errors.Wrap(err, "Error reading build output from docker")
	}

	return nil
}

// Push pushes the image to the docker registry
func (d *DockerImpl) Push(image string, stdout io.Writer, stderr io.Writer) error {
	if len(image) == 0 {
		return errEmptyImageName
	}
	ctx := context.Background()
	res, err := d.imageClient.ImagePush(ctx, image, types.ImagePushOptions{})
	if res != nil {
		defer res.Close()
	}
	if err != nil {
		return errors.Wrap(err, "Error sending push request to docker")
	}
	if err = printStreamResponse(res, stdout); err != nil {
		return errors.Wrap(err, "Error reading push output from docker")
	}

	return nil
}

func parsePorts(ports []int32) (nat.PortMap, nat.PortSet, error) {
	portBindings := make(nat.PortMap)
	exposedPorts := make(nat.PortSet)
	for _, port := range ports {
		if port <= 0 || port > 65535 {
			return nil, nil, errInvalidPort(port)
		}
		// TODO: some people might need UDP ports, will probably need to switch the slice of int32 to something else
		natPort, err := nat.NewPort("tcp", strconv.Itoa(int(port)))
		exposedPorts[natPort] = struct{}{}
		if err != nil {
			return nil, nil, errors.Wrap(err, "Error converting the port number")
		}
		portBinding := []nat.PortBinding{{HostPort: strconv.Itoa(int(port))}}

		portBindings[natPort] = portBinding
	}

	return portBindings, exposedPorts, nil
}

// Run creates and starts a container with the given image and name, and runtime options (e.g. exposed ports) specified in the config parameter
func (d *DockerImpl) Run(image string, name string, config *Runtime, stdout io.Writer, stderr io.Writer) error {
	if len(image) == 0 {
		return errEmptyImageName
	}

	if len(name) == 0 {
		return errEmptyContainerName
	}

	if config == nil {
		return errNilRuntimeConfig
	}

	portBindings, exposedPorts, err := parsePorts(config.Ports)
	if err != nil {
		return errors.Wrap(err, "Error parsing configuration ports")
	}

	ctx := context.Background()
	if _, err := d.containerRunner.ContainerCreate(ctx, &container.Config{Image: image, ExposedPorts: exposedPorts},
		&container.HostConfig{PortBindings: portBindings}, nil, name); err != nil {
		return errors.Wrap(err, "Error creating the container")
	}

	if err := d.containerRunner.ContainerStart(ctx, name, types.ContainerStartOptions{}); err != nil {
		return errors.Wrap(err, "Error starting the container")
	}
	return nil
}

// Logs attaches to the container with the given name and prints the log to stdout
func (d *DockerImpl) Logs(name string, stdout io.Writer, stderr io.Writer) error {
	if len(name) == 0 {
		return errEmptyContainerName
	}
	ctx := context.Background()
	res, err := d.containerRunner.ContainerLogs(ctx, name, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true})
	if err != nil {
		return errors.Wrap(err, "Error getting container logs")
	}
	defer res.Close()

	_, err = stdcopy.StdCopy(stdout, stderr, res)
	return err
}

// Cancel stops and removes the container with the given name
func (d *DockerImpl) Cancel(name string) error {
	if len(name) == 0 {
		return errEmptyContainerName
	}
	ctx := context.Background()
	timeout := 60 * time.Second

	if err := d.containerRunner.ContainerStop(ctx, name, &timeout); err != nil {
		return errors.Wrap(err, "Error stopping the container")
	}
	if err := d.containerRunner.ContainerRemove(ctx, name, types.ContainerRemoveOptions{}); err != nil {
		return errors.Wrap(err, "Error removing the container")
	}

	return nil
}
