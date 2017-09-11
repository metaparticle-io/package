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
		[Metaparticle.Runtime.Config(Ports = new int[] {5000}, Executor = "aci", Public = true)]
        [Metaparticle.Package.Config(Repository = "brendanburns", Publish = true)]
        public static void Main(string[] args) => Containerize(args, () =>
       	{
            WebHost.CreateDefaultBuilder(args)
                .UseStartup<Startup>()
				.UseKestrel(options => { options.Listen(IPAddress.Any, 5000); })
                .Build()
                .Run();
    	});
    }
}
