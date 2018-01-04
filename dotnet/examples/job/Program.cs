using System;
using System.Threading;
using static Metaparticle.Package.Driver;


namespace simple {
	public class Program {
        [Metaparticle.Package.Config(Publish = true, Repository = "brendanburns/dotnet-simple")]
        [Metaparticle.Runtime.Config(JobCount = 4)]
        public static void Main(string[] args) => Containerize (args, () =>
        {
			Console.Out.WriteLine(args);
			Console.WriteLine("Hello world batch job!");
        });
    }
}