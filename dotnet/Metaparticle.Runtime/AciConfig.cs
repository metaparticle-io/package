namespace Metaparticle.Runtime
{
    /// <summary>
    /// Runtime config for AciExecutor.
    /// </summary>
    public class AciConfig : Config
    {
        /// <summary>
        /// Azure tenant id.
        /// </summary>
        public string AzureTenantId { get; set; }

        /// <summary>
        /// Azure subscription id.
        /// Container instances will be created using the given subscription.
        /// </summary>
        public string AzureSubscriptionId { get; set; }

        /// <summary>
        /// Azure client id (aka. application id).
        /// </summary>
        public string AzureClientId { get; set; }

        /// <summary>
        /// Azure client secret (aka. application key).
        /// If specified, then will authenticate to Azure as service principal.
        /// And then no need to specify username/password.
        /// </summary>
        public string AzureClientSecret { get; set; }

        public string Username { get; set; }

        public string Password { get; set; }

        /// <summary>
        /// The resource group to create container instance.
        /// </summary>
        public string AciResourceGroup { get; set; } = "metaparticle-execution";

        public AciConfig()
        {
            Executor = "aci";
        }
    }
}
