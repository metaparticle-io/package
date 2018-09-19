using System.Collections.Generic;
using System.IO;
using System.Text;
using dockerfile;

namespace Metaparticle.Package
{
	public static class DockerfileWriter
	{
		/// <summary>
		/// Creates the dockerfle and returns the name of the file
		/// </summary>
		/// <param name="dir">Directory to create the file</param>
		/// <param name="exe">Executable used to run command/s</param>
		/// <param name="args">args for the executable</param>
		/// <param name="config"></param>
		/// <returns></returns>
		public static string Write(string dir, string exe, string[] args, Config config)
		{
			var dockerfilename = dir + "/Dockerfile";
			if (!string.IsNullOrEmpty(config.Dockerfile)) {
				File.Copy(config.Dockerfile, dockerfilename);
				return dockerfilename;
			}
			var instructions = new List<Instruction>();
			instructions.Add(new Instruction("FROM", "debian:9"));
			instructions.Add(new Instruction("RUN", " apt-get update && apt-get -qq -y install libunwind8 libicu57 libssl1.0 liblttng-ust0 libcurl3 libuuid1 libkrb5-3 zlib1g"));
			// TODO: lots of things here are constant, figure out how to cache for perf?
			instructions.Add(new Instruction("COPY", "* /exe/"));
			instructions.Add(new Instruction("CMD", $"/exe/{exe} {GetArgs(args)}"));

			var df = new Dockerfile(instructions.ToArray(), new Comment[0]);
			File.WriteAllText(dockerfilename, df.Contents());

			return dockerfilename;
		}

		private static string GetArgs(string[] args)
		{
			if (args == null || args.Length == 0)
			{
				return "";
			}
			var b = new StringBuilder();
			foreach (var arg in args)
			{
				b.Append(arg);
				b.Append(" ");
			}
			return b.ToString().Trim();
		}
	}
}