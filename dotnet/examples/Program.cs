using System;
using System.Threading;
using Metaparticle.Package;


namespace examples {
	public class Program {
        [MetaparticleConfig(Verbose = false)]
        public static void Main(string[] args) => Metaparticle.Package.Metaparticle.Run(args, (a) =>
        {
			int i = 0;
            while (true) {
				Console.WriteLine("Hello world " + (i++));
				Thread.Sleep(10 * 1000);
			}
        });
    }
}
