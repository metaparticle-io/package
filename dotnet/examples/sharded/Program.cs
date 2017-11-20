using System.IO;
using System.Net;
using Microsoft.AspNetCore;
using Microsoft.AspNetCore.Hosting;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Logging;
using static Metaparticle.Package.Driver;


namespace web
{
    public class Program
    {
        const int port = 8080;
		[Metaparticle.Runtime.Config(Ports = new int[] {port},
                                     Executor = "metaparticle",
                                     Shards = 3,
                                     ShardExpression = "^\\/users\\/([^\\/]*)\\/.*",
                                     Public = true)]
        [Metaparticle.Package.Config(Repository = "brendanburns/dotnet-web",
                                     Publish = true,
                                     Verbose = true)]
        public static void Main(string[] args) => Containerize(args, () =>
       	{
            WebHost.CreateDefaultBuilder(args)
                .UseStartup<Startup>()
				.UseKestrel(options => { options.Listen(IPAddress.Any, port); })
                .Build()
                .Run();
    	});
    }
}
