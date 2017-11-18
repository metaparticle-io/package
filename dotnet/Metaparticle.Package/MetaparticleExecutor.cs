using System;
using System.Collections.Generic;
using System.IO;
using System.Threading;
using Metaparticle.Runtime;
using static Metaparticle.Package.Util;

namespace Metaparticle.Package {
    public class MetaparticleExecutor : ContainerExecutor
    {
        public void Cancel(string id)
        {
            HandleErrorExec("kubectl", "delete deployments " + id);
            HandleErrorExec("kubectl", "delete services " + id);
        }

        private void HandleErrorExec(string cmd, string args, TextWriter stdout=null) {
            var err = new StringWriter();
            var proc = Exec(cmd, args, stdout: stdout, stderr: err);
            if (proc.ExitCode != 0) {
                Console.WriteLine(err.ToString());
            }
        }

        public void Logs(string id, TextWriter stdout, TextWriter stderr)
        {
            var args = string.Format("-l app={0} --template \"{{{{.Message}}}}\"", id);
            Console.WriteLine("ktail " + args);
            Exec("ktail", args, stdout: stdout, stderr: stderr);
        }

        public string Run(string image, Runtime.Config config)
        {
            // TODO: make this name better...
            var name = "server";

            string replicaSpec = null;
            string shardSpec = null;

            if (config.Shards > 0) {
                shardSpec = string.Format(@"""shardSpec"": {{
                    ""shards"": {0},
                    ""urlPattern"": ""{1}""
                }", config.Shards, config.ShardExpression);
            } else {
                replicaSpec = string.Format(@"""replicas"": {0}", config.Replicas);
            }

            var spec = @"{{
    ""name"": ""name"",
    ""guid"": 1234567, 
    ""services"": [ 
        {{
               ""name"": ""{0}"",
            {},
            ""containers"": [
                {{
                    ""image"": ""{2}"",
                    ""env"": [{{
                        ""name"": ""METAPARTICLE_IN_CONTAINER"",
                        ""value"": ""true""
                    }}]
                }}
            ],
            ""ports"": [{{
                ""number"": {3}
            }}]
        }}
    ],
    ""serve"": {{
        ""name"": ""{0}"",
        ""public"": true
    }}
}}";
            var specFileName = Path.Combine(Path.GetTempPath(), "spec.json");
            File.WriteAllText(specFileName, string.Format(spec, name, config.Replicas, image, config.Ports[0]));

            HandleErrorExec("/home/bburns/gopath/bin/compiler", string.Format("-f {0}", specFileName));
            return name;
        }
    }
}