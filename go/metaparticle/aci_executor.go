package metaparticle

import "io"
import "os/exec"
import "time"
import "fmt"

// ACIExecutor implements and executor on Azure Container Engine
type ACIExecutor struct {
	ResourceGroup string
}

// Run creates a container on the azure executor's resource group using the given image and name
func (a *ACIExecutor) Run(image string, name string, cfg *Runtime, stdout io.Writer, stderr io.Writer) error {
	if len(image) == 0 {
		return fmt.Errorf("An image must be specified")
	}

	if len(name) == 0 {
		return fmt.Errorf("The container's name must be specified")
	}

	if cfg == nil {
		return fmt.Errorf("cfg must not be nil")
	}

	cmdName := "az"
	cmdParams := []string{"container", "create", "--image", image, "-g", a.ResourceGroup, "-n", name, "--env=METAPARTICLE_IN_CONTAINER=true"}

	for _, port := range cfg.Ports {
		cmdParams = append(cmdParams, "--port="+string(port))
	}

	if cfg.PublicAddress {
		cmdParams = append(cmdParams, "--ip-address=Public")
	}

	cmd := exec.Command(cmdName, cmdParams...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd.Run()
}

// Logs shows the logs of the container with the given name on the azure executor's resource group
func (a *ACIExecutor) Logs(name string, stdout io.Writer, stderr io.Writer) error {
	if len(name) == 0 {
		return fmt.Errorf("A container name must be specified")
	}
	cmd := exec.Command("az", "container", "logs", "-g", a.ResourceGroup, "-n", name)
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	// I don't know why java's implementation sleeps for 5 seconds after running logs, so I'm just doing the same here.
	// But if I were to guess, it's probably just a wait time for azure to respond with the logs
	<-time.After(5 * time.Second)
	return nil
}

// Cancel deletes the container with the given name on the executor's resource group
func (a *ACIExecutor) Cancel(name string) error {
	cmd := exec.Command("az", "container", "delete", "-g", a.ResourceGroup, "-n", name, "--yes")
	return cmd.Run()
}
