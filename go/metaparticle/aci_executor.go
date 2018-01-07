package metaparticle

import (
	"fmt"
	"io"
	"os/exec"
	"time"
)

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
	ResourceGroup string
}

// ACIExecutor implements and executor on Azure Container Engine
type ACIExecutor struct {
	config *AciRuntimeConfig
}

// Run creates a container on the azure executor's resource group using the given image and name
func (a *ACIExecutor) Run(image string, name string, cfg *Runtime, stdout io.Writer, stderr io.Writer) error {
	if len(image) == 0 {
		return errEmptyImageName
	}

	if len(name) == 0 {
		return errEmptyContainerName
	}

	if cfg == nil {
		return fmt.Errorf("cfg must not be nil")
	}

	cmdName := "az"
	cmdParams := []string{"container", "create", "--image", image, "-g", a.config.ResourceGroup, "-n", name, "--env=METAPARTICLE_IN_CONTAINER=true"}

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
	cmd := exec.Command("az", "container", "logs", "-g", a.config.ResourceGroup, "-n", name)
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
	cmd := exec.Command("az", "container", "delete", "-g", a.config.ResourceGroup, "-n", name, "--yes")
	return cmd.Run()
}
