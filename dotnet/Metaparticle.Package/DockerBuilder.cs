using System.IO;
using static Metaparticle.Package.Util;

namespace Metaparticle.Package {
    public class DockerBuilder : ImageBuilder
    {
        public bool Build(string configFile, string imageName, TextWriter o, TextWriter e)
        {
            System.Console.WriteLine(configFile);
            var info = Directory.GetParent(configFile);
            System.Console.WriteLine($"build -t {imageName} {info.FullName}");
            var err = new StringWriter();
            var proc = Exec("docker", $"build -t {imageName} {info.FullName}", o, err);
            System.Console.WriteLine(err.ToString());
            return proc.ExitCode == 0;
        }

        public bool Push(string imageName, TextWriter stdout, TextWriter stderr)
        {
            var proc = Exec("docker", $"push {imageName}", stdout, stderr);
            return proc.ExitCode == 0;
        }
    }
}