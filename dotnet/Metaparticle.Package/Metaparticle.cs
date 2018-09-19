using System;
using System.Diagnostics;
using System.IO;
using Metaparticle.Tests;
using static Metaparticle.Package.Util;
using RuntimeConfig = Metaparticle.Runtime.Config;
using TestConfig = Metaparticle.Tests.Config;

namespace Metaparticle.Package
{
	public class Driver
    {
        private readonly Config config;
        private readonly RuntimeConfig runtimeConfig;
        private readonly TestConfig testConfig;

	    public Driver(Config config, RuntimeConfig runtimeConfig, TestConfig testConfig)
        {
            this.config = config;
            this.runtimeConfig = runtimeConfig;
            this.testConfig = testConfig;
        }

	    private ImageBuilder getBuilder()
        {
            switch (config.Builder.ToLowerInvariant())
            {
                case "docker":
                    return new DockerBuilder();
                case "aci":
                    return new DockerBuilder();
                default:
                    return null;
            }
        }

        private ContainerExecutor getExecutor()
        {
            if (runtimeConfig == null) {
                return null;
            }
            switch (runtimeConfig.Executor.ToLowerInvariant())
            {
                case "docker":
                    return new DockerExecutor();
                case "aci":
                    return new AciExecutor(runtimeConfig);
                case "metaparticle":
                    return new MetaparticleExecutor();
                default:
                    return null;
            }
        }

        public void Build(string[] args)
        {
            var proc = Process.GetCurrentProcess();
            var procName = proc.ProcessName;
            string exe = null;
            string dir;
            TextWriter o = config.Verbose ? Console.Out : null;
            TextWriter e = config.Quiet ? Console.Error : null;
            if (procName == "dotnet")
            {
                RunTests();

                dir = "bin/release/netcoreapp2.0/debian.8-x64/publish";
                Exec("dotnet", "publish -r debian.8-x64 -c release", stdout: o, stderr: e);
                var files = Directory.GetFiles(dir);
                foreach (var filePath in files)
                {
                    var file = new FileInfo(filePath);
                    if (file.Name.EndsWith(".runtimeconfig.json"))
                    {
                        exe = file.Name.Substring(0, file.Name.Length - ".runtimeconfig.json".Length);
                    }
                }
            }
            else
            {
                exe = procName;
                var prog = proc.MainModule.FileName;
                dir = Directory.GetParent(prog).FullName;
            }
            var dockerfilename = DockerfileWriter.Write(dir, exe, args, config);
            var builder = getBuilder();

            string imgName = (string.IsNullOrEmpty(config.Repository) ? exe : config.Repository);
            if (!string.IsNullOrEmpty(config.Version)) {
                imgName += ":" + config.Version;
            }
            if (!builder.Build(dockerfilename, imgName, stdout: o, stderr: e))
            {
                Console.Error.WriteLine("Image build failed.");
                return;
            }

            if (config.Publish)
            {
                if (!builder.Push(imgName, stdout: o, stderr: e))
                {
                    Console.Error.WriteLine("Image push failed.");
                    return;
                }
            }

            if (runtimeConfig == null)
            {
                return;
            }

            var exec = getExecutor();
            if (exec.PublishRequired() && !config.Publish)
            {
                Console.Error.WriteLine("Image publish is required, but image was not published. Set publish to true in the package config.");
                return;
            }
            var id = exec.Run(imgName, runtimeConfig);

            Console.CancelKeyPress += delegate
            {
                exec.Cancel(id);
            };

            exec.Logs(id, Console.Out, Console.Error);
        }

        private void RunTests()
        {
            var runTestsResult = new DotnetTestRunner().Run(testConfig.Names);
            if (runTestsResult == false)
            {
                throw new Exception("Tests Failed.");   
            }
        }

	    public static bool InDockerContainer()
        {
            switch (Environment.GetEnvironmentVariable("METAPARTICLE_IN_CONTAINER"))
            {
                case "true":
                case "1":
                    return true;
                case "false":
                case "0":
                    return false;
            }
            // This only works on Linux
            const string cgroupPath = "/proc/1/cgroup";
            if (File.Exists(cgroupPath)) {
                var info = File.ReadAllText(cgroupPath);
                // This is a little approximate...
                return info.IndexOf("docker", StringComparison.Ordinal) != -1 || info.IndexOf("kubepods", StringComparison.Ordinal) != -1;
            }
            return false;
        }

        public static void Containerize(string[] args, Action main)
        {
            if (InDockerContainer())
            {
                main();
                return;
            }
            Config config = new Config();
            RuntimeConfig runtimeConfig = null;
            TestConfig testConfig = null;

            var trace = new StackTrace();
            foreach (object attribute in trace.GetFrame(1).GetMethod().GetCustomAttributes(true))
            {
                if (attribute is Config)
                {
                    config = (Config) attribute;
                }
                if (attribute is RuntimeConfig)
                {
                    runtimeConfig = (RuntimeConfig) attribute;
                }
                if (attribute is TestConfig)
                {
                    testConfig = (TestConfig) attribute;
                }
            }
            var mp = new Driver(config, runtimeConfig, testConfig);
            mp.Build(args);
        }
    }
}
