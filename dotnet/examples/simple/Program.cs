using System;
using System.Threading;
using static Metaparticle.Package.Driver;


namespace simple 
{
	public class Program 
	{
        [Metaparticle.Runtime.Config]
        [Metaparticle.Package.Config(Verbose = true,
			Publish = false, Repository = "docker.io/docker-user-goes-here/dotnet-simple")] 
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
