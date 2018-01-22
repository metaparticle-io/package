using System.IO;
using System.Linq;
using System.Net;
using System.Threading.Tasks;
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
        [Metaparticle.Runtime.Config(Ports = new int[] {8080})]
        [Metaparticle.Package.Config(Repository = "docker.io/docker-user-goes-here/dotnet-simple-web", Publish = false, Verbose = true)]
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
