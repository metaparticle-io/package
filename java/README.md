# Metaparticle/Package for Java
Metaparticle/Package is a collection of libraries intended to 
make building and deploying containers a seamless and idiomatic
experience for developers.

This is the implementation for Java.

## Introduction
Metaparticle/Package simplifies and centralizes the task of
building and deploying a container image.

Here is a quick example.

Consider this simple Java application:

```java
public class Main {
    public static void main(String[] args) {
        System.out.println("Hello world!");
    }
}
```

To containerize this application, you need to add the `@Package` annotation and the `Containerize` wrapper function
like this:

```java
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

You can then compile this application just as you have before.
But now, when you run the application, instead of printing "Hello world", it first packages itself as a container, and
then (optionally) deploys itself inside that container.

## Tutorial
For a more complete exploration of the Metaparticle/Package for Java, please see the [in-depth tutorial](tutorial.md).