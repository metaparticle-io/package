# Metaparticle/Package

Language Idiomatic bindings for building Container Images.

## What's this about?
Containers are an optimal way to package and deploy your code. However, teaching developers to learn a new
configuration file format, and toolchain, just to package their application in a container is an
unnecessary barrier to entry for many programmers just starting out with containers.

Metaparticle/Package simplifies the task of building and deploying container images. Metaparticle/Package is
a collection of libraries that enable programmers to build and deploy containers using code that feels
familiar to them.

## Can you give me an example?
Here's a simple example of building a containerized Java application:

```Java
import io.metaparticle.annotations.Package;
import static io.metaparticle.Metaparticle.Containerize;

public class Main {
    @Package(repository="brendanburns",
             jarFile="path/to/my-fat-jar.jar")
    public static void main(String[] args) {
        Containerize(() -> {
            System.out.println("Hello Metaparticle/Package");
        });
    }
}
```

When you run this program via the `java` command or your IDE, rather than simply executing your code, this program
packages up the Java code in a container, and runs that container.

## What languages do you support?

Currently:
   * [java](java)
   * [.NET core](dotnet)
   * [javascript](javascript) (NodeJS)
   * [go](go)
   * [ruby](ruby)

But it's fairly straightforward to add other languages, we would love to see contributions.

## Details

For more details see the more complete walkthroughs for each language:
   * [java tutorial](tutorials/java/tutorial.md)
   * [.NET Core tutorial](tutorials/dotnet/tutorial.md)
   * [javascript tutorial](tutorials/javascript/tutorial.md)
   * [Ruby tutorial](tutorials/ruby/tutorial.md)

## Contribute
There are many ways to contribute to Metaparticle

 * [Submit bugs](https://github.com/metaparticle-io/package/issues) and help us verify fixes as they are checked in.
 * Review the source code changes.
 * Engage with other Metaparticle users and developers on [gitter](https://gitter.im/metaparticle-io/Lobby).
 * Join the #metaparticle discussion on [Twitter](https://twitter.com/MetaparticleIO).
 * [Contribute bug fixes](https://github.com/metaparticle-io/package/pulls).

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/). For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or contact [opencode@microsoft.com](mailto://opencode@microsoft.com) with any additional questions or comments.
