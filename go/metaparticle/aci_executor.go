package metaparticle

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/to"

	"github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2017-08-01-preview/containerinstance"
)

var (
	defaultLocation = "westeurope"

	defaultActiveDirectoryEndpoint = azure.PublicCloud.ActiveDirectoryEndpoint
	defaultResourceManagerEndpoint = azure.PublicCloud.ResourceManagerEndpoint
)

// ACIRuntimeConfig - runtime config for ACI executor.
type ACIRuntimeConfig struct {
	SubscriptionID string
	TenantID       string
	ClientID       string
	ClientSecret   string

	ResourceGroupName string
}

// ACIExecutor implements and executor on Azure Container Engine
type ACIExecutor struct {
	config *ACIRuntimeConfig

	groupsClient containerinstance.ContainerGroupsClient
	logsClient   containerinstance.ContainerLogsClient
}

// NewACIExecutor returns a new instance of an ACI Executor
func NewACIExecutor() (*ACIExecutor, error) {
	cfg := &ACIRuntimeConfig{
		SubscriptionID: getEnvVarOrExit("AZURE_SUBSCRIPTION_ID"),
		TenantID:       getEnvVarOrExit("AZURE_TENANT_ID"),
		ClientID:       getEnvVarOrExit("AZURE_CLIENT_ID"),
		ClientSecret:   getEnvVarOrExit("AZURE_CLIENT_SECRET"),

		ResourceGroupName: getEnvVarOrExit("RESOURCE_GROUP_NAME"),
	}

	g, l, err := getACIClients(cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot get ACI clients: %v", err)
	}
	return &ACIExecutor{
		config:       cfg,
		groupsClient: g,
		logsClient:   l,
	}, nil
}

// Run creates a container on the azure executor's resource group using the given image and name
func (a *ACIExecutor) Run(image string, name string, cfg *Runtime, stdout io.Writer, stderr io.Writer) error {

	ports := []containerinstance.Port{}
	cPorts := []containerinstance.ContainerPort{}
	for _, p := range cfg.Ports {
		ports = append(ports, containerinstance.Port{
			Port:     to.Int32Ptr(p),
			Protocol: containerinstance.TCP,
		})

		cPorts = append(cPorts, containerinstance.ContainerPort{
			Port: to.Int32Ptr(p),
		})
	}

	var ipType string
	if cfg.PublicAddress {
		ipType = "Public"
	} else {
		ipType = "Private"
	}

	// for now there are default values for resource requests and limits
	// TODO - add to config struct

	parameters := containerinstance.ContainerGroup{
		Name:     to.StringPtr(name),
		Location: &defaultLocation,
		ContainerGroupProperties: &containerinstance.ContainerGroupProperties{
			IPAddress: &containerinstance.IPAddress{
				Type:  &ipType,
				Ports: &ports,
			},
			OsType: containerinstance.Linux,
			Containers: &[]containerinstance.Container{
				{
					Name: to.StringPtr(name),
					ContainerProperties: &containerinstance.ContainerProperties{
						Ports: &cPorts,
						Image: &image,
						Resources: &containerinstance.ResourceRequirements{
							Limits: &containerinstance.ResourceLimits{
								MemoryInGB: to.Float64Ptr(1),
								CPU:        to.Float64Ptr(1),
							},
							Requests: &containerinstance.ResourceRequests{
								MemoryInGB: to.Float64Ptr(1),
								CPU:        to.Float64Ptr(1),
							},
						},
					},
				},
			},
		},
	}

	// TODO - add context support
	c, err := a.groupsClient.CreateOrUpdate(context.Background(), a.config.ResourceGroupName, name, parameters)
	if err != nil {
		return fmt.Errorf("cannot create container group: %v", err)
	}

	fmt.Printf("started container group - to access, wait for it to be running, then go to: %v\n", *c.IPAddress.IP)

	return nil
}

// Logs shows the logs of the container with the given name on the azure executor's resource group
func (a *ACIExecutor) Logs(name string, stdout io.Writer, stderr io.Writer) error {
	// TODO - add context support
	// TODO - add log tail support
	<-time.After(10 * time.Second)

	for {
		c, err := a.groupsClient.Get(context.Background(), a.config.ResourceGroupName, name)
		if err != nil {
			fmt.Printf("cannot get container group: %v", err)
		}

		if *c.State == "Running" {
			break
		} else {
			fmt.Printf("Waiting for container to be in running state; current state: %v\n", *c.State)
			<-time.After(10 * time.Second)
		}
	}

	go func() {
		for {
			logs, err := a.logsClient.List(context.Background(), a.config.ResourceGroupName, name, name, to.Int32Ptr(10))
			if err != nil {
				<-time.After(10 * time.Second)
			}

			if logs.Content != nil {
				fmt.Printf("logs\n %v", *logs.Content)
				<-time.After(10 * time.Second)
			}
		}
	}()

	return nil
}

// Cancel deletes the container with the given name on the executor's resource group
func (a *ACIExecutor) Cancel(name string) error {
	// TODO - add context support
	_, err := a.groupsClient.Delete(context.Background(), a.config.ResourceGroupName, name)
	if err != nil {
		return fmt.Errorf("cannot delete container group: %v", err)
	}

	return nil
}

func getACIClients(cfg *ACIRuntimeConfig) (containerinstance.ContainerGroupsClient, containerinstance.ContainerLogsClient, error) {
	var containerGroupsClient containerinstance.ContainerGroupsClient
	var containerLogsClient containerinstance.ContainerLogsClient

	oAuthConfig, err := adal.NewOAuthConfig(defaultActiveDirectoryEndpoint, cfg.TenantID)
	if err != nil {
		return containerGroupsClient, containerLogsClient, fmt.Errorf("cannot get oAuth configuration: %v", err)
	}

	token, err := adal.NewServicePrincipalToken(*oAuthConfig, cfg.ClientID, cfg.ClientSecret, defaultResourceManagerEndpoint)
	if err != nil {
		return containerGroupsClient, containerLogsClient, fmt.Errorf("cannot get service principal token: %v", err)
	}
	containerGroupsClient = containerinstance.NewContainerGroupsClient(cfg.SubscriptionID)
	containerGroupsClient.Authorizer = autorest.NewBearerAuthorizer(token)

	containerLogsClient = containerinstance.NewContainerLogsClient(cfg.SubscriptionID)
	containerLogsClient.Authorizer = autorest.NewBearerAuthorizer(token)

	return containerGroupsClient, containerLogsClient, nil
}

func getEnvVarOrExit(varName string) string {
	value := os.Getenv(varName)
	if value == "" {
		log.Fatalf("missing environment variable %s\n", varName)
	}

	return value
}
