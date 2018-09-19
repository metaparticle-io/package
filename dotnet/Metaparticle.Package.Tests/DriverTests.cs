using System;
using NSubstitute;
using Xunit;
using static Metaparticle.Package.Driver;
using RuntimeConfig = Metaparticle.Runtime.Config;

namespace Metaparticle.Package.Tests
{
    public class DriverTests
    {
        static int testOutput = 0;

        [Metaparticle.Runtime.Config]
        [Metaparticle.Package.Config(Verbose = true,
		Publish = false, Repository = "testrepo")] 
        private static void TestActionInContainer(string[] args) => Containerize (args, () =>
        {
            testOutput = 1;
        });

        [Fact]
        public void Containerize_Executes_Action_When_In_Container_Equals_True()
        {
            testOutput = 0;
            Environment.SetEnvironmentVariable("METAPARTICLE_IN_CONTAINER", "true");
            TestActionInContainer(null);
            Assert.True(1 == testOutput);
        }

        [Fact]
        public void Containerize_Executes_Action_When_In_Container_Equals_1()
        {
            testOutput = 0;
            Environment.SetEnvironmentVariable("METAPARTICLE_IN_CONTAINER", "1");
            TestActionInContainer(null);
            Assert.True(1 == testOutput);
        }
    }
}
