using System;
using System.Threading;
using static Metaparticle.Package.Driver;


namespace simple {
	public class Program {
        [Metaparticle.Package.Config(Verbose = true,
			Publish = true, Repository = "brendanburns")] 
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
