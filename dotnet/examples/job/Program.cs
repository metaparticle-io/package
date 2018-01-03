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
			int i = 0;
            while (true) {
				Console.WriteLine("Hello world " + (i++));
				Thread.Sleep(10 * 1000);
			}
        });
    }
}