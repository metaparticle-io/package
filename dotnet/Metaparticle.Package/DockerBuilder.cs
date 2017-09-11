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
            var info = Directory.GetParent(configFile);
            var proc = Exec("docker", string.Format("build -t {0} {1}", imageName, info.FullName), stdout: o, stderr: e);
            return proc.ExitCode == 0;
        }

        public bool Push(string imageName, TextWriter stdout, TextWriter stderr)
        {
            var proc = Exec("docker", string.Format("push {0}", imageName), stdout: stdout, stderr: stderr);
            return proc.ExitCode == 0;
        }
    }
}