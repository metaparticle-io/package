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
   * [dotnet core](dotnet)
   * [javascript](javascript) (NodeJS)

But it's fairly straightforward to add other languages, we would love to see contributions.

## Details

For more details see the more complete walkthroughs for each language:
   * [java tutorial](java/tutorial.md)
   * [dotnet tutorial](dotnet/tutorial.md)
   * [javascript tutorial](javascript/tutorial.md)

