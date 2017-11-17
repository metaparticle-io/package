# Metaparticle/Package for Java Tutorial
This is an in-depth tutorial for using Metaparticle/Package for Java

For a quick summary, please see the [README](README.md).

## Initial Setup

### Check the tools
The `docker` command line tool needs to be installed and working. Try:
`docker ps` to verify this.

### Get the code
```sh
$ git clone https://github.com/metaparticle-io/package
$ cd tutorials/dotnet/
# [optional, substitute your favorite editor here...]
$ code .
```

## Initial Program
Inside of the `tutorials/dotnet` directory, you will find a simple maven project.

You can build this project with `dotnet build`.

The initial code is a very simple "Hello World"

```cs
using System.IO;
using System.Linq;
using System.Net;
using System.Threading.Tasks;
using Microsoft.AspNetCore;
using Microsoft.AspNetCore.Hosting;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Logging;

namespace web
{
    public class Program
    {
        const int port = 8080;
        public static void Main(string[] args)
       	{
            WebHost.CreateDefaultBuilder(args)
                .UseStartup<Startup>()
				.UseKestrel(options => { options.Listen(IPAddress.Any, port); })
                .Build()
                .Run();
    	}
    }
}
```

## Step One: Containerize the Application
To build a container from our simple application we need to add a dependency to our
build file, and then update the code.

Run:
```sh
dotnet add package Metaparticle.Package
```

Then update the code to read as follows:

```cs
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
        [Metaparticle.Package.Config(Repository = "docker.io/docker-user-goes-here/simple-web", Publish = false)]
        public static void Main(string[] args) => Containerize(args, () =>
       	{
            WebHost.CreateDefaultBuilder(args)
                .UseStartup<Startup>()
				.UseKestrel(options => { options.Listen(IPAddress.Any, port); })
                .Build()
                .Run();
    	});
    }
}```

You will notice that we added a `Metaparticle.Package.Config` annotation that describes how
to package the application. You will need to replace `your-docker-user-goes-here`
with an actual Docker repository path.

You will also notice that we wrapped the main function in the `Containerize`
function which kicks off the Metaparticle code.

You can run this new program with:

```sh
dotnet run
```

This code will start your web server again. But this time, it is running
inside a container. You can see this by running:

```sh
docker ps
```

## Step Two: Exposing the ports
If you try to access the web server on [http://localhost:8080](http://localhost:8080) you
will see that you can not actually access the server. Despite it running, the service
is not exposed. To do this, you need to add a `[Metaparticle.Runtime.Config ...]` annotation to supply the
port(s) to expose.

The code snippet to add is:

```cs
...
    [Metaparticle.Runtime.Config Ports = int[] {8080}]
...
```

This tells the runtime the port(s) to expose. The complete code looks like:

```cs
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
        [Metaparticle.Package.Config(Repository = "brendanburns/dotnet-simple-web", Publish = true, Verbose = true)]
        [Metaparticle.Runtime.Config Ports = int[] {port}]
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
```

Now if you run this with `mvn compile exec:java -Dexec.mainClass=io.metaparticle.tutorial.Main` your webserver will be successfully exposed on port 8080.

## Replicating and exposing on the web.
As a final step, consider the task of exposing a replicated service on the internet.
To do this, we're going to expand our suage of the `@Runtime` tag. First we will
add a `replicas` field, which will specify the number of replicas. Second we will
set our execution environment to `metaparticle` which will launch the service
into the currently configured Kubernetes environment.

Here's what the snippet looks like:

```cs
...
    [Metaparticle.Runtime.Config Ports = int[] {port}, Executor = "metaparticle", Replicas = 4]
...
```

And the complete code looks like:
```cs
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
        [Metaparticle.Package.Config(Repository = "brendanburns/dotnet-simple-web", Publish = true, Verbose = true)]
        [Metaparticle.Runtime.Config Ports = int[] {port}, Executor = "metaparticle", Replicas = 4]
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
```

After you compile and run this, you can see that there are four replicas running behind a
Kubernetes Service Load balancer:

```sh
$ kubectl get pods
...
$ kubectl get services
...
```