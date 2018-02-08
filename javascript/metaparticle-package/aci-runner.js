(function() {
var msRestAzure = require('ms-rest-azure');
var resourceManagementClient = require('azure-arm-resource');
var containerInstanceManagementClient = require('azure-arm-containerinstance');

var localOptions = null;

function deleteContainer(aciClient, name) {
  aciClient.containerGroups.deleteMethod(localOptions.resourceGroup, name)
      .then(cg => {
        console.log(`Deleted container group ${name}`);
      })
      .catch(err => {
        console.log(`Failed to deleted container group ${name}`);
      });
}

function getContainerLogs(aciClient, name) {
  aciClient.containerLogs.list(localOptions.resourceGroup, name, name)
      .then(log => {
        console.log('===== Tail logs start =====');
        console.log(log.content);
        console.log('===== Tail logs end =====');
      })
      .catch(err => {
        console.log(`Failed to fetch log for container group ${name}; logs might not be ready yet.`);
      })
}

function createContainer(rgClient, aciClient, img, name, options) {
  rgClient.resourceGroups.get(options.resourceGroup)
      .then(resourceGroup => {
        let ipAddress = null;
        let ports = null;

        if (options && options.public) {
          ipAddress = {
            type: 'Public'
          };
          if (options.ports) {
            ports = options.ports.map((p) => {
              return {port: p, protocol: 'TCP'};
            });
            ipAddress.ports = ports;
          }
        }

        let containerGroup = {
          name: name,
          location: resourceGroup.location,
          osType: 'Linux',
          restartPolicy: 'Never',
          ipAddress: ipAddress,
          containers: [{
            name: name,
            image: img,
            resources: {requests: {memoryInGB: 1.5, cpu: 1.0}},
            environmentVariables:
                [{name: 'METAPARTICLE_IN_CONTAINER', value: 'true'}],
            ports: ports
          }]
        };

        aciClient.containerGroups
            .createOrUpdate(options.resourceGroup, name, containerGroup)
            .then(cg => {
              console.log(`Created container group with id ${cg.id}`);
              setInterval(() => {
                getContainerLogs(aciClient, name);
              }, 5000);
            })
            .catch(err => {
              console.log(`Failed to create container group: ${err.message}`);
            });
      })
      .catch(err => {
        console.log(`Failed to get resource group ${options.resourceGroup}: ${err.message}`);
      })
}

module.exports.run = (img, name, options) => {
  localOptions = options;

  msRestAzure
      .loginWithServicePrincipalSecret(
          options.clientId, options.clientSecret, options.tenantId)
      .then(credential => {
        const rgClient = new resourceManagementClient.ResourceManagementClient(
            credential, options.subscriptionId);
        const aciClient = new containerInstanceManagementClient(
            credential, options.subscriptionId);
        createContainer(rgClient, aciClient, img, name, options);
      })
      .catch(err => {
        console.log(`Failed to authenticate: ${err.message}`);
      });
};

module.exports.cancel = (name) => {
  msRestAzure
      .loginWithServicePrincipalSecret(
          localOptions.clientId, localOptions.clientSecret, localOptions.tenantId)
      .then(credential => {
        const aciClient = new containerInstanceManagementClient(
            credential, localOptions.subscriptionId);
        deleteContainer(aciClient, name);
      })
      .catch(err => {
        console.log(`Failed to authenticate: ${err.message}`);
      });
};

})();
