# Metaparticle/Package for Python Tutorial
This is an in-depth tutorial for using Metaparticle/Package for Python

For a quick summary, please see the [README](README.md).

## Initial Setup

### Check the tools
The `docker` command line tool needs to be installed and working. Try:
`docker ps` to verify this.  Go to the [install page](https://get.docker.io) if you need
to install Docker.

The `mp-compiler` command line tool needs to be installed and working.
Try `mp-compiler --help` to verify this. Go to the [releases page](https://github.com/metaparticle-io/metaparticle-ast/releases) if you need to install
the Metaparticle compiler.


# Install Metaparticle/Package

## Install the library
```sh
pip install metaparticle_pkg
```

## Get the tutorial code
```sh
git clone https://github.com/metaparticle-io/package
cd package/tutorials/python
```

## Initial Program
Inside of the `tutorials/python` directory, you will find a simple python project.

```python
from six.moves import SimpleHTTPServer, socketserver
import socket

OK = 200

port = 8080

class MyHandler(SimpleHTTPServer.SimpleHTTPRequestHandler):
    def do_GET(self):
        self.send_response(OK)
        self.send_header("Content-type", "text/plain")
        self.end_headers()
        self.wfile.write("Hello Metparticle [{}] @ {}\n".format(self.path, socket.gethostname()).encode('UTF-8'))
        print("request for {}".format(self.path))
    def do_HEAD(self):
        self.send_response(OK)
        self.send_header("Content-type", "text/plain")
        self.end_headers()

def main():
    Handler = MyHandler
    httpd = socketserver.TCPServer(("", port), Handler)
    httpd.serve_forever()

if __name__ == '__main__':
    main()
```

## Step One: Containerize the Application
To build a container from our simple application we need to add a dependency to our
code and then update the code to read as follows:

```python
from six.moves import SimpleHTTPServer, socketserver
import socket
from metaparticle_pkg.containerize import Containerize

OK = 200

port = 8080

class MyHandler(SimpleHTTPServer.SimpleHTTPRequestHandler):
    def do_GET(self):
        self.send_response(OK)
        self.send_header("Content-type", "text/plain")
        self.end_headers()
        self.wfile.write("Hello Metparticle [{}] @ {}\n".format(self.path, socket.gethostname()).encode('UTF-8'))
        print("request for {}".format(self.path))
    def do_HEAD(self):
        self.send_response(OK)
        self.send_header("Content-type", "text/plain")
        self.end_headers()

@Containerize(
    package={'name': 'web', 'repository': 'docker.io/brendanburns', 'publish': True},
)
def main():
    Handler = MyHandler
    httpd = socketserver.TCPServer(("", port), Handler)
    httpd.serve_forever()

if __name__ == '__main__':
    main()
```

You will notice that we added a `@Containerize` annotation that describes how
to package the application. You will need to replace `your-docker-user-goes-here`
with an actual Docker repository path.

You can run this new program with:

```sh
python web.py
```

This code will start your web server again. But this time, it is running
inside a container. You can see this by running:

```sh
docker ps
```

## Step Two: Exposing the ports
If you try to access the web server on [http://localhost:8080](http://localhost:8080) you
will see that you can not actually access the server. Despite it running, the service
is not exposed. To do this, you need to add am annotation to supply the
port(s) to expose.

The code snippet to add is:

```python
...
@Containerize(
    package={'name': 'web', 'repository': 'docker.io/brendanburns', 'publish': True},
    runtime={'ports': [8080]}
)
...
```

This tells the runtime the port(s) to expose. The complete code looks like:

```python
import SimpleHTTPServer
import SocketServer
import socket
from metaparticle_pkg.containerize import Containerize

OK = 200

port = 8080

class MyHandler(SimpleHTTPServer.SimpleHTTPRequestHandler):
    def do_GET(self):
        self.send_response(OK)
        self.send_header("Content-type", "text/plain")
        self.end_headers()
        self.wfile.write("Hello Metparticle [{}] @ {}\n".format(self.path, socket.gethostname()))
        print("request for {}".format(self.path))
    def do_HEAD(self):
        self.send_response(OK)
        self.send_header("Content-type", "text/plain")
        self.end_headers()

@Containerize(
    package={'name': 'web', 'repository': 'docker.io/brendanburns', 'publish', True},
    runtime={'ports': [8080]}
)
def main():
    Handler = MyHandler
    httpd = SocketServer.TCPServer(("", port), Handler)
    httpd.serve_forever()

if __name__ == '__main__':
    main()
```

Now if you run this with `python web.py` your webserver will be successfully exposed on port 8080.

## Replicating and exposing on the web.
As a final step, consider the task of exposing a replicated service on the internet.
To do this, we're going to expand our usage of the `@containerize` tag. First we will
add a `replicas` field, which will specify the number of replicas. Second we will
set our execution environment to `metaparticle` which will launch the service
into the currently configured Kubernetes environment.

Here's what the snippet looks like:

```python
...
@Containerize(
    package={'name': 'web', 'repository': 'docker.io/brendanburns'},
    runtime={'ports': [8080], 'executor': 'metaparticle', 'replicas': 4}
)
...
```

And the complete code looks like:
```python
import SimpleHTTPServer
import SocketServer
import socket
from metaparticle_pkg.containerize import Containerize

OK = 200

port = 8080

class MyHandler(SimpleHTTPServer.SimpleHTTPRequestHandler):
    def do_GET(self):
        self.send_response(OK)
        self.send_header("Content-type", "text/plain")
        self.end_headers()
        self.wfile.write("Hello Metparticle [{}] @ {}\n".format(self.path, socket.gethostname()))
        print("request for {}".format(self.path))
    def do_HEAD(self):
        self.send_response(OK)
        self.send_header("Content-type", "text/plain")
        self.end_headers()

@Containerize(
    package={'name': 'web', 'repository': 'docker.io/brendanburns'},
    runtime={'ports': [8080], 'executor': 'metaparticle', 'replicas': 4}
)
def main():
    Handler = MyHandler
    httpd = SocketServer.TCPServer(("", port), Handler)
    httpd.serve_forever()

if __name__ == '__main__':
    main()
```

After you run this, you can see that there are four replicas running behind a
Kubernetes Service Load balancer:

```sh
$ kubectl get pods
...
$ kubectl get services
...
```

