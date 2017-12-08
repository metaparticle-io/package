# Metaparticle/Package for .NET Core
Metaparticle/Package is a collection of libraries intended to
make building and deploying containers a seamless and idiomatic
experience for developers.

This is the implementation for .NET Core (C#).

## Introduction
Metaparticle/Package simplifies and centralizes the task of
building and deploying a container image.

Here is a quick example.

Consider this simple .NET Core application:

```cs
namespace simple {
    public class Program {
        public static void main(string[] args) {
            Console.WriteLine("Hello world!");
        }
    }
}
```

To containerize this application, you need to add the `[Metaparticle.Package.Config ...]` annotation and the `Containerize` wrapper function
like this:

```cs
using static Metaparticle.Package.Driver;

namespace simple {
	public class Program {
        [Metaparticle.Package.Config(Repository = "brendanburns/dotnet-simple")]
        public static void Main(string[] args) => Containerize (args, () =>
        {
			Console.WriteLine("Hello world!");
        });
    }
}
```

You can then compile this application just as you have before.
But now, when you run the application, instead of printing "Hello world", it first packages itself as a container, and
then (optionally) deploys itself inside that container.

## Tutorial
For a more complete exploration of the Metaparticle/Package for .NET Core, please see the [in-depth tutorial](../tutorials/dotnet/tutorial.md).
