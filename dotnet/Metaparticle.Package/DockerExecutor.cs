using System.IO;
using static Metaparticle.Package.Util;

namespace Metaparticle.Package {
    public class DockerExecutor : ContainerExecutor
    {
        public void Cancel(string id)
        {
            Exec("docker", string.Format("kill {0}", id));
        }

        public string Run(string image, Metaparticle.Runtime.Config config)
        {
            
            var idWriter = new StringWriter();
            Exec("docker", string.Format("run {0} -d {1}", portString(config.Ports), image), stdout: idWriter);
            return idWriter.ToString().Trim();
        }

        public void Logs(string id, TextWriter stdout, TextWriter stderr) {
            Exec("docker", string.Format("logs -f {0}", id), stdout, stderr);
        }

        public bool PublishRequired() {
            return false;
        }

        private string portString(int[] ports) {
            if (ports == null || ports.Length == 0) {
                return "";
            }
            var pieces = new string[ports.Length];
            for (int i = 0; i < ports.Length; i++) {
                pieces[i] = string.Format("-p {0}:{1}", ports[i], ports[i]);
            }
            return string.Join(",", pieces);
        }
    }
}