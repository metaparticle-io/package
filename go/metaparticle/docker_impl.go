package metaparticle

import (
	"fmt"
	"io"
	"os/exec"
)

type DockerImpl struct{}

func (d *DockerImpl) Build(dir string, image string, stdout io.Writer, stderr io.Writer) error {
	if len(dir) == 0 {
		return fmt.Errorf("A context directory must be specified")
	}
	if len(image) == 0 {
		return fmt.Errorf("An image name must be specified")
	}
	cmd := exec.Command("docker", "build", "-t", image, dir)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd.Run()
}

func (d *DockerImpl) Push(image string, stdout io.Writer, stderr io.Writer) error {
	if len(image) == 0 {
		return fmt.Errorf("An image name must be specified")
	}
	cmd := exec.Command("docker", "push", image)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd.Run()
}

func (d *DockerImpl) Run(image string, name string, config *Runtime, stdout io.Writer, stderr io.Writer) error {
	if len(image) == 0 {
		return fmt.Errorf("An image must be specified")
	}

	if len(name) == 0 {
		return fmt.Errorf("The container's name must be specified")
	}

	if config == nil {
		return fmt.Errorf("config must not be nil")
	}
	cmdName := "docker"
	cmdParams := []string{"run", "-d"}
	for _, port := range config.Ports {
		cmdParams = append(cmdParams, "-p", fmt.Sprintf("%d:%d", port, port))
	}
	cmdParams = append(cmdParams, "--name", name, image)

	cmd := exec.Command(cmdName, cmdParams...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd.Run()
}

func (d *DockerImpl) Logs(name string, stdout io.Writer, stderr io.Writer) error {
	if len(name) == 0 {
		return fmt.Errorf("A container name must be specified")
	}
	cmd := exec.Command("docker", "logs", "-f", name)
	cmd.Stderr = stderr
	cmd.Stdout = stdout

	return cmd.Run()
}

func (d *DockerImpl) Cancel(name string) error {
	cmd := exec.Command("docker", "stop", name)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Could not stop the container: %v", err)
	}
	cmd = exec.Command("docker", "rm", "-f", name)
	return cmd.Run()
}
