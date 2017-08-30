using System.IO;
using static Metaparticle.Package.Util;

namespace Metaparticle.Package {
    public class DockerExecutor : ContainerExecutor
    {
        public void Cancel(string id)
        {
            Exec("docker", string.Format("kill {0}", id));
        }

        public string Run(string image)
        {
            var idWriter = new StringWriter();
            Exec("docker", string.Format("run -d {0}", image), stdout: idWriter);
            return idWriter.ToString().Trim();
        }

        public void Logs(string id, TextWriter stdout, TextWriter stderr) {
            Exec("docker", string.Format("logs -f {0}", id), stdout, stderr);
        }
    }
}