using System;
using Xunit;
using Metaparticle.Package;

namespace Metaparticle.Package.Tests
{
    public class ConfigTests
    {
        [Fact]
        public void Config_Loads_Environment_Variables()
        {
            Environment.SetEnvironmentVariable("METAPARTICLE_CONFIG_REPOSITORY", "testRepo");
            Environment.SetEnvironmentVariable("METAPARTICLE_CONFIG_PUBLISH", "true");
            var config = new Config();
            Assert.True("testRepo" == config.Repository);
            Assert.True(true == config.Publish); 
        }

        [Fact]
        public void Config_Defaults_Builder_To_Docker()
        {
            var config = new Config();
            Assert.True("docker" == config.Builder);
        }
    }
}
