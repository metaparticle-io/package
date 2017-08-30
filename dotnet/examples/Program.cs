using System;
using System.Threading;
using static Metaparticle.Package.Metaparticle;


namespace examples {
	public class Program {
		//[Metaparticle.Runtime.Config(Replicas = 3)]
        [Metaparticle.Package.Config(Verbose = false)] 
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
