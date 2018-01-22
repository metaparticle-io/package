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

### Environment variables

You can set some of the config attributes through environment variables so that it can be controlled from outside of the source code. This could be useful for CICD scenarios.

E.g. 

```
set MP_CONFIG_REPOSITORY=docker.io/myrepo/myimagename:sometag
```

This will set the `Repository` property that you would otherwise set in the attributes. See `Config.cs` for supported environment variable overrides.

### Tests
If you wish to add some test project to your metaparticle that get run as part of the build pipeline, you can add the tests projects as a CSV to an environment variable.

```
set METAPARTICLE_TESTS_CSV=../relpath/project1,../relpath2/project2
```

You can test this using the following example

```
set METAPARTICLE_TESTS_CSV=../simple-test
cd examples/simple
dotnet run
```

## Tutorial
For a more complete exploration of the Metaparticle/Package for .NET Core, please see the [in-depth tutorial](../tutorials/dotnet/tutorial.md).
