# Metaparticle/Package for Java Tutorial
This is an in-depth tutorial for using Metaparticle/Package for Java

For a quick summary, please see the [README](README.md).

## Initial Setup

### Check the tools
The `docker` command line tool needs to be installed and working. Try:
`docker ps` to verify this.  Go to the [install page](https://get.docker.io) if you need
to install Docker.

The `mp-compiler` command line tool needs to be installed and working.
Try `mp-compiler --help` to verify this. Go to the [releases page](https://github.com/metaparticle-io/metaparticle-ast/releases) if you need to install
the Metaparticle compiler.


### Install dependent libraries
For now, you need to manually install the needed libraries locally. We'll
get them up into Maven central shortly...

```sh
# Install the client library

git clone https://github.com/metaparticle-io/metaparticle-java
cd metaparticle-java
mvn install

# Install Metaparticle/Package

git clone https://github.com/metaparticle-io/package
cd package/java
mvn install
```

### Get the code
```sh
$ git clone https://github.com/metaparticle-io/package
$ cd package/tutorials/java/
# [optional, substitute your favorite editor here...]
$ code .
```

## Initial Program
Inside of the `tutorials/java` directory, you will find a simple maven project.

You can build this project with `mvn compile`.

The initial code is a very simple "Hello World"

```java
package io.metaparticle.tutorial;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpServer;

public class Main {
    private static final int port = 8080;

    public static void main(String[] args) {
        try {
            HttpServer server = HttpServer.create(new InetSocketAddress(port), 0);
            server.createContext("/", new HttpHandler() {
                @Override
                public void handle(HttpExchange t) throws IOException {
                    String msg = "Hello Containers [" + t.getRequestURI() + "] from " + System.getenv("HOSTNAME");
                    t.sendResponseHeaders(200, msg.length());
                    OutputStream os = t.getResponseBody();
                    os.write(msg.getBytes());
                    os.close();
                    System.out.println("[" + t.getRequestURI() + "]");
                }
            });
            server.start();
        } catch (IOException ex) {
            ex.printStackTrace();
        }
    }
}
```

## Step One: Containerize the Application
To build a container from our simple application we need to add a dependency to our
build file, and then update the code.

Add:
```xml
...
    <dependency>
      <groupId>io.metaparticle</groupId>
      <artifactId>metaparticle-package</artifactId>
      <version>0.1-SNAPSHOT</version>
    </dependency>
...
```

Then update the code to read as follows:

```java
package io.metaparticle.tutorial;

import io.metaparticle.annotations.Package;
import static io.metaparticle.Metaparticle.Containerize;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpServer;

public class Main {
    private static final int port = 8080;

    @Package(repository = "your-docker-user-goes-here",
             jarFile = "target/metaparticle-package-tutorial-0.1-SNAPSHOT-jar-with-dependencies.jar")
    public static void main(String[] args) {
        Containerize(() -> {
            try {
                HttpServer server = HttpServer.create(new InetSocketAddress(port), 0);
                server.createContext("/", new HttpHandler() {
                    @Override
                    public void handle(HttpExchange t) throws IOException {
                        String msg = "Hello Containers [" + t.getRequestURI() + "] from " + System.getenv("HOSTNAME");
                        t.sendResponseHeaders(200, msg.length());
                        OutputStream os = t.getResponseBody();
                        os.write(msg.getBytes());
                        os.close();
                        System.out.println("[" + t.getRequestURI() + "]");
                    }
                });
                server.start();
            } catch (IOException ex) {
                ex.printStackTrace();
            }
        });
    }
}
```

You will notice that we added a `@Package` annotation that describes how
to package the application. You will need to replace `your-docker-user-goes-here`
with an actual Docker Hub user name.

You will also notice that we wrapped the main function in the `Containerize`
function which kicks off the Metaparticle code.

Once you have this, you can build the code with `mvn compile`.  If you see an error like: 
`Could not find artifact io.metaparticle:metaparticle-package:jar:0.1-SNAPSHOT `
You may need to move into `$BASE/package/java` and run `mvn install`.

Once the code is compiled, you can run this new program with:

```sh
mvn exec:java -Dexec.mainClass=io.metaparticle.tutorial.Main
```

This code will start your web server again. But this time, it is running
inside a container. You can see this by running:

```sh
docker ps
```

## Step Two: Exposing the ports
If you try to access the web server on [http://localhost:8080](http://localhost:8080) you
will see that you can not actually access the server. Despite it running, the service
is not exposed. To do this, you need to add a `@Runtime` annotation to supply the
port(s) to expose.

The code snippet to add is:

```java
...
    @Runtime(ports={port})
...
```

This tells the runtime the port(s) to expose. The complete code looks like:

```java
package io.metaparticle.tutorial;

import io.metaparticle.annotations.Package;
import io.metaparticle.annotations.Runtime;
import static io.metaparticle.Metaparticle.Containerize;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpServer;

public class Main {
    private static final int port = 8080;

    @Runtime(ports={port})        
    @Package(repository = "your-docker-user-goes-here",
             jarFile = "target/metaparticle-package-tutorial-0.1-SNAPSHOT-jar-with-dependencies.jar")
    public static void main(String[] args) {
        Containerize(() -> {
            try {
                HttpServer server = HttpServer.create(new InetSocketAddress(port), 0);
                server.createContext("/", new HttpHandler() {
                    @Override
                    public void handle(HttpExchange t) throws IOException {
                        String msg = "Hello Containers [" + t.getRequestURI() + "] from " + System.getenv("HOSTNAME");
                        t.sendResponseHeaders(200, msg.length());
                        OutputStream os = t.getResponseBody();
                        os.write(msg.getBytes());
                        os.close();
                        System.out.println("[" + t.getRequestURI() + "]");
                    }
                });
                server.start();
            } catch (IOException ex) {
                ex.printStackTrace();
            }
        });
    }
}
```

Now if you run this with `mvn compile exec:java -Dexec.mainClass=io.metaparticle.tutorial.Main` your webserver will be successfully exposed on port 8080.

## Replicating and exposing on the web.
As a final step, consider the task of exposing a replicated service on the internet.
To do this, we're going to expand our usage of the `@Runtime` and `@Package` annotations.
First, `@Package` has its `publish` field added with a value of `true`. This is necessary in order to push the built image to the Docker repository.
Next, we add a `executor` field to the `@Runtime` annotation to set our execution environment to `metaparticle` which will launch the service
into the currently configured Kubernetes environment. Finally, we add a `replicas` field to the `@Runtime` annotation. This specifies the number of replicas to schedule.

Here's what the snippet looks like:

```java
...
    @Runtime(ports={port},
             replicas = 4,
             executor = "metaparticle")
    @Package(repository = "your-docker-user-goes-here",
            jarFile = "target/metaparticle-package-tutorial-0.1-SNAPSHOT-jar-with-dependencies.jar",
            publish = true,
            verbose = true)           
...
```

And the complete code looks like:
```java
package io.metaparticle.tutorial;

import io.metaparticle.annotations.Package;
import io.metaparticle.annotations.Runtime;
import static io.metaparticle.Metaparticle.Containerize;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpServer;

public class Main {
    private static final int port = 8080;

    @Runtime(ports={port},
             replicas=4,
             executor="metaparticle")    
    @Package(repository = "your-docker-user-goes-here",
            jarFile = "target/metaparticle-package-tutorial-0.1-SNAPSHOT-jar-with-dependencies.jar",
            publish = true,
            verbose = true)           
    public static void main(String[] args) {
        Containerize(() -> {
            try {
                HttpServer server = HttpServer.create(new InetSocketAddress(port), 0);
                server.createContext("/", new HttpHandler() {
                    @Override
                    public void handle(HttpExchange t) throws IOException {
                        String msg = "Hello Containers [" + t.getRequestURI() + "] from " + System.getenv("HOSTNAME");
                        t.sendResponseHeaders(200, msg.length());
                        OutputStream os = t.getResponseBody();
                        os.write(msg.getBytes());
                        os.close();
                        System.out.println("[" + t.getRequestURI() + "]");
                    }
                });
                server.start();
            } catch (IOException ex) {
                ex.printStackTrace();
            }
        });
    }
}
```

After you compile and run this, you can see that there are four pod replicas running behind a
Kubernetes ClusterIP service:

```sh
$ kubectl get pods
...
$ kubectl get services
...
```

Still looking for more? Continue on to the more advanced [sharding tutorial](sharding-tutorial.md)
