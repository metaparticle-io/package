using System.IO;
using static Metaparticle.Package.Util;

namespace Metaparticle.Package {
    public class DockerBuilder : ImageBuilder
    {
        public DockerBuilder()
        {
        }

        public bool Build(string configFile, string imageName, TextWriter o, TextWriter e)
        {
            System.Console.WriteLine(configFile);
            var info = Directory.GetParent(configFile);
            System.Console.WriteLine(string.Format("build -t {0} {1}", imageName, info.FullName));
            var err = new StringWriter();
            var proc = Exec("docker", string.Format("build -t {0} {1}", imageName, info.FullName), stdout: o, stderr: err);
            System.Console.WriteLine(err.ToString());
            return proc.ExitCode == 0;
        }

        public bool Push(string imageName, TextWriter stdout, TextWriter stderr)
        {
            var proc = Exec("docker", string.Format("push {0}", imageName), stdout: stdout, stderr: stderr);
            return proc.ExitCode == 0;
        }
    }
}