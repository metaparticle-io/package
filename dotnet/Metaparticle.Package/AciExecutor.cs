using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Threading;
using RuntimeConfig = Metaparticle.Runtime.Config;
using Metaparticle.Runtime;
using Microsoft.Azure.Management.ContainerInstance;
using Microsoft.Azure.Management.ContainerInstance.Models;
using Microsoft.Azure.Management.ResourceManager.Fluent;
using Microsoft.Azure.Management.ResourceManager.Fluent.Authentication;
using Microsoft.Azure.Management.ResourceManager.Fluent.Models;

namespace Metaparticle.Package
{
    public class AciExecutor : ContainerExecutor
    {
        private readonly AciConfig aciConfig;

        public AciExecutor(RuntimeConfig runtimeConfig)
        {
            if (!(runtimeConfig is AciConfig))
            {
                throw new ArgumentException("AciExecutor can only accept AciConfig as runtime config");
            }

            this.aciConfig = runtimeConfig as AciConfig;
            this.ValidateAciConfig(this.aciConfig);
        }

        public void Cancel(string id)
        {
            var containerInstanceClient = this.GetContainerInstanceClient();
            containerInstanceClient.ContainerGroups.Delete(
                resourceGroupName: this.aciConfig.AciResourceGroup,
                containerGroupName: id);
        }

        public void Logs(string id, TextWriter stdout, TextWriter stderr)
        {
            var containerInstanceClient = this.GetContainerInstanceClient();

            while (true)
            {
                try {
                    var logs = containerInstanceClient.ContainerLogs.List(
                        resourceGroupName: this.aciConfig.AciResourceGroup,
                        containerGroupName: id,
                        containerName: id);

                    stdout.WriteLine("===== Tail logs start =====");
                    stdout.Write(logs.Content);
                    stdout.WriteLine("===== Tail logs end =====");
                    stdout.WriteLine();
                    stdout.Flush();

                    Thread.Sleep(5 * 1000);
                } catch (Exception) {
                    // Logs won't be ready until the image pulling is done.
                    stderr.WriteLine("Logs not ready yet...");
                    stderr.Flush();

                    Thread.Sleep(10 * 1000);
                }
            }
        }

        public string Run(string image, RuntimeConfig config)
        {
            var resourceGroup = this.GetResourceGroup();
            var containerGroupName = $"metaparticle-exec-{DateTime.UtcNow.Ticks}";

            var containers = new Container[]
            {
                new Container(
                name: containerGroupName,
                image: image,
                resources: new ResourceRequirements (requests: new ResourceRequests(memoryInGB: 1.5, cpu: 1.0)),
                environmentVariables: new List<EnvironmentVariable>
                {
                    new EnvironmentVariable(name: "METAPARTICLE_IN_CONTAINER", value: "true")
                },
                ports: config.Ports?.Select(p => new ContainerPort(p)).ToList())
            };

            var containerGroup = new ContainerGroup(
                name: containerGroupName,
                osType: OperatingSystemTypes.Linux,
                location: resourceGroup.Location,
                ipAddress: config.Public
                    ? new IpAddress (config.Ports?.Select (p => new Port(p, ContainerNetworkProtocol.TCP)).ToList())
                    : null,
                containers: containers
            );

            var containerInstanceClient = this.GetContainerInstanceClient();

            var createdContainerGroup = containerInstanceClient.ContainerGroups.CreateOrUpdate(
                resourceGroupName: this.aciConfig.AciResourceGroup,
                containerGroupName: containerGroupName,
                containerGroup: containerGroup);

            Console.WriteLine(createdContainerGroup.Id);
            return containerGroupName;
        }

        private ContainerInstanceManagementClient GetContainerInstanceClient()
        {
            var azureCredentials = this.GetAzureCredentials();

            return new ContainerInstanceManagementClient(azureCredentials)
            {
                BaseUri = new Uri(AzureEnvironment.AzureGlobalCloud.ResourceManagerEndpoint),
                SubscriptionId = this.aciConfig.AzureSubscriptionId
            };
        }

        private ResourceManagementClient GetResourceClient()
        {
            var azureCredentials = this.GetAzureCredentials();
            return new ResourceManagementClient(azureCredentials)
            {
                BaseUri = new Uri(AzureEnvironment.AzureGlobalCloud.ResourceManagerEndpoint),
                SubscriptionId = this.aciConfig.AzureSubscriptionId
            };
        }

        private ResourceGroupInner GetResourceGroup()
        {
            var resourceManagementClient = this.GetResourceClient();

            var resourceGroup = resourceManagementClient.ResourceGroups.GetAsync(this.aciConfig.AciResourceGroup).Result;
            
            if (resourceGroup == null)
            {
                throw new ArgumentException($"Provided aciConfig.AciResourceGroup {aciConfig.AciResourceGroup} doesn't exist");
            }

            return resourceGroup;
        }

        private AzureCredentials GetAzureCredentials()
        {
            var azureCredentialsFactory = new AzureCredentialsFactory();
            return string.IsNullOrWhiteSpace(this.aciConfig.AzureClientSecret)
                ? azureCredentialsFactory.FromUser(
                    username: this.aciConfig.Username,
                    password: this.aciConfig.Password,
                    clientId: this.aciConfig.AzureClientId,
                    tenantId: this.aciConfig.AzureTenantId,
                    environment: AzureEnvironment.AzureGlobalCloud)
                : azureCredentialsFactory.FromServicePrincipal(
                    clientId: this.aciConfig.AzureClientId,
                    clientSecret: this.aciConfig.AzureClientSecret,
                    tenantId: this.aciConfig.AzureTenantId,
                    environment: AzureEnvironment.AzureGlobalCloud);
        }

        private void ValidateAciConfig(AciConfig aciConfig)
        {
            if (string.IsNullOrWhiteSpace(aciConfig.AzureTenantId))
            {
                throw new ArgumentNullException("aciConfig.AzureTenantId");
            }

            if (string.IsNullOrWhiteSpace(aciConfig.AzureSubscriptionId))
            {
                throw new ArgumentNullException("aciConfig.AzureSubscriptionId");
            }

            if (string.IsNullOrWhiteSpace(aciConfig.AzureClientId))
            {
                throw new ArgumentNullException("aciConfig.AzureClientId");
            }

            if (string.IsNullOrWhiteSpace(aciConfig.AzureClientSecret)
                && (string.IsNullOrWhiteSpace(aciConfig.Username)
                || string.IsNullOrWhiteSpace(aciConfig.Password)))
            {
                throw new ArgumentException("Must specify either aciConfig.AzureClientSecret or aciConfig.Username and aciConfig.Password");
            }
        }

        public bool PublishRequired() {
            return true;
        }
    }
}