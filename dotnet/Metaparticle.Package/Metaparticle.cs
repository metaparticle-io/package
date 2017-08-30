using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Threading;
using System.Threading.Tasks;
using Mono.Unix;
using dockerfile;
using System.Text;

namespace MetaparticlePackage
{
    public class Metaparticle
    {
        private MetaparticleConfig config;

        public Metaparticle(MetaparticleConfig config) {
            this.config = config;
        }


        public void Exec(String file, String args, TextWriter stdout=null, TextWriter stderr=null)
        {
            var task = ExecAsync(file, args, stdout, stderr);
            task.Wait();
        }

        public string ExecOld(String file, String args, bool verboseOverride=false)
        {
            var o = new StringWriter();
            var e = new StringWriter();
            var task = ExecAsync(file, args, o, e);
            task.Wait();
            var output = o.ToString();
            var err = e.ToString();
            if ((config != null && config.Verbose) || verboseOverride)
            {
                Console.Write(output);
                Console.Write(err);
            }

            return output;
        }

        public static async Task Copy(StreamReader reader, TextWriter writer) {
            if (reader == null || writer == null) {
                return;
            }
            var read = await reader.ReadLineAsync();
            while (read != null) {
                await writer.WriteLineAsync(read);
                await writer.FlushAsync();
                read = await reader.ReadLineAsync();
            }
        }

        public async Task ExecAsync(String file, String args, TextWriter stdout=null, TextWriter stderr=null)
        {
            var startInfo = new ProcessStartInfo(file, args);
            startInfo.RedirectStandardOutput = true;
            startInfo.RedirectStandardError = true;
            var proc = Process.Start(startInfo);
            var runTask = Task.Run(() => { proc.WaitForExit(); });
            var outputTask = Copy(proc.StandardOutput, stdout);
            var errTask = Copy(proc.StandardError, stderr);

            await Task.Run(() => {
                Task.WaitAll(new []{
                    runTask,
                    outputTask,
                    errTask,
                });
            });
            return;
        }

        private static string getArgs(string[] args) {
            if (args == null || args.Length == 0) {
                return "";
            }
            var b = new StringBuilder();
            foreach (var arg in args) {
                b.Append(arg);
                b.Append(" ");
            }
            return b.ToString().Trim();
        }

        public void Build(string[] args)
        {
            var proc = Process.GetCurrentProcess();
            var procName = proc.ProcessName;
            string exe = null;
            string dir = null;
            TextWriter o = config.Verbose ? Console.Out : null;
            TextWriter e = config.Verbose ? Console.Error : null;
            if (procName == "dotnet")
            {
                dir = "bin/release/netcoreapp2.0/debian.8-x64/publish";
                Exec("/opt/dotnet/dotnet", "publish -r debian.8-x64 -c release", stdout: o, stderr: e);
                var dirInfo = new UnixDirectoryInfo(dir);
                foreach (var file in dirInfo.GetFileSystemEntries())
                {
                    if (file.Name.EndsWith(".pdb"))
                    {
                        exe = file.Name.Substring(0, file.Name.Length - 4);
                    }
                }
            }
            else
            {
                exe = procName;
                var prog = proc.MainModule.FileName;
                dir = Directory.GetParent(prog).FullName;
            }
            var instructions = new List<Instruction>();
            instructions.Add(new Instruction("FROM", "debian:9"));
            instructions.Add(new Instruction("RUN", " apt-get update && apt-get install libunwind8 libicu57"));
            instructions.Add(new Instruction("COPY", string.Format("* /exe/", dir)));
            instructions.Add(new Instruction("CMD", string.Format("/exe/{0} {1}", exe, getArgs(args))));

            var df = new Dockerfile(instructions.ToArray(), new Comment[0]);
            File.WriteAllText(dir + "/Dockerfile", df.Contents());

            string imgName = (string.IsNullOrEmpty(config.Repository) ? exe : config.Repository + "/" + exe);
            if (!string.IsNullOrEmpty(config.Version)) {
                imgName += ":" + config.Version;
            }
            var info = new UnixDirectoryInfo(dir);
            Exec("docker", string.Format("build -t {0} {1}", imgName, info.FullName), stdout: o, stderr: e);

            var idWriter = new StringWriter();
            Exec("docker", string.Format("run -d {0}", imgName), stdout: idWriter);
            var id = idWriter.ToString().Trim();
            Exec("docker", string.Format("logs -f {0}", id), Console.Out, Console.Error);
        }

        public static bool InDockerContainer()
        {
            var info = File.ReadAllText("/proc/1/cgroup");
            // This is a little approximate...
            return info.IndexOf("docker") != -1;
        }

        public static void Run(string[] args, Action<string[]> main)
        {
            if (InDockerContainer())
            {
                main(args);
                return;
            }
            MetaparticleConfig config = new MetaparticleConfig();
            var trace = new StackTrace();
            foreach (object attribute in trace.GetFrame(1).GetMethod().GetCustomAttributes(true))
            {
                if (attribute is MetaparticleConfig)
                {
                    config = (MetaparticleConfig)attribute;
                }
            }
            var mp = new Metaparticle(config);
            mp.Build(args);
        }
    }
}
