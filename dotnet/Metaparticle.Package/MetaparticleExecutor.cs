using System;
using System.Collections.Generic;
using System.IO;
using System.Threading;
using Metaparticle.Runtime;
using static Metaparticle.Package.Util;

namespace Metaparticle.Package 
{
    public class MetaparticleExecutor : ContainerExecutor
    {
        public void Cancel(string id)
        {
            var specFileName = Path.Combine(".metaparticle", "spec.json");

            HandleErrorExec("mp-compiler", string.Format("-f {0} --delete", specFileName));
        }

        private void HandleErrorExec(string cmd, string args, TextWriter stdout=null)
        {
            var err = new StringWriter();
            var proc = Exec(cmd, args, stdout: stdout, stderr: err);
            if (proc.ExitCode != 0) {
                Console.WriteLine(err.ToString());
            }
        }

        public void Logs(string id, TextWriter stdout, TextWriter stderr)
        {
            var specFileName = Path.Combine(".metaparticle", "spec.json");

            var args = string.Format("-f {0} --deploy=false --attach=true", specFileName);
            Exec("mp-compiler", args, stdout: stdout, stderr: stderr);
        }

        public string Run(string image, Runtime.Config config)
        {
            // TODO: make this name better...
            var name = "server";

            string replicaSpec = null;

            if (config.Shards > 0) {
                replicaSpec = string.Format(@"""shardSpec"": {{
                    ""shards"": {0},
                    ""urlPattern"": ""{1}""
                }}", config.Shards, config.ShardExpression);
            } else {
                replicaSpec = string.Format(@"""replicas"": {0}", config.Replicas);
            }

            string servicesSpec;
            if (config.Shards > 0 || config.Replicas > 0) {
                servicesSpec = @",
            ""services"": [ 
        {{
               ""name"": ""{0}"",
            {1},
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
    }}";
            } else {
                servicesSpec = "";
            }

    string jobSpec;

            if (config.JobCount > 0) {
                jobSpec = @",
                ""jobs"": [{{
                    ""name"": {0},
                    ""containers"": [
                    {{
                        ""image"": ""{1}"",
                        ""env"": [{{
                            ""name"": ""METAPARTICLE_IN_CONTAINER"",
                            ""value"": ""true""
                        }}]
                    }}
                }}]";
            } else {
                jobSpec = "";
            }

    var spec = @"{{
    ""name"": ""{0}"",
    ""guid"": 1234567
    {}
    {}
}}";
            Directory.CreateDirectory(".metaparticle");
            var specFileName = Path.Combine(".metaparticle", "spec.json");
            var fullServiceSpec = string.Format(servicesSpec, replicaSpec, image, config.Ports[0]);
            var fullJobSpec = string.Format(jobSpec, name, image);
            File.WriteAllText(specFileName, string.Format(spec, name, fullServiceSpec, fullJobSpec));

            HandleErrorExec("mp-compiler", string.Format("-f {0}", specFileName));
            return name;
        }
    
        public bool PublishRequired() {
            return true;
        }
    }
}
