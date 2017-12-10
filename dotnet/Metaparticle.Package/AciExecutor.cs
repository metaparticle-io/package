using System;
using System.IO;
using System.Threading;
using Newtonsoft.Json.Linq;
using static Metaparticle.Package.Util;
using RuntimeConfig = Metaparticle.Runtime.Config;

namespace Metaparticle.Package
{
    public class AciExecutor : ContainerExecutor
    {
        public AciExecutor() {
        }

        public void Cancel(string id)
        {
            var err = new StringWriter();
            var proc = Exec("az", String.Format("container delete -g test -n {0} --yes", id), stderr: err);
            if (proc.ExitCode != 0) {
                Console.WriteLine(err.ToString());
            }
        }

        public void Logs(string id, TextWriter stdout, TextWriter stderr)
        {
            while (true) {
                var proc = Exec("az", String.Format("container logs -g test -n {0}", id), stdout: stdout, stderr: stderr);
                Thread.Sleep(5 * 1000);
            }
        }

        public string Run(string image, RuntimeConfig config)
        {
            // TODO: handle multiple ports
            var portStr = "";
            if (config.Ports != null && config.Ports.Length > 0) {
                portStr = string.Format("--port {0}", config.Ports[0]);
            }

            var pubStr = "";
            if (config.Public) {
                pubStr = "--ip-address Public";
            }

            var name = "test";
            var data = new StringWriter();
            var proc = Exec("az", string.Format("container create --image={0} -g test -n {1} {2} {3} --env=METAPARTICLE_IN_CONTAINER=true", image, name, portStr, pubStr), stdout: data);
            var json = data.ToString();
            var obj = JObject.Parse(json);
            var id = obj["id"];
            Console.WriteLine(id.ToString());
            return name;
        }
    }
}