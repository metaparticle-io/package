using System;
using MetaparticlePackage;


namespace examples {
	public class Program {
        [MetaparticleConfig(Verbose = false)]
        public static void Main(string[] args) => Metaparticle.Run(args, (a) =>
        {
            Console.WriteLine("Hello world");
        });
    }
}
