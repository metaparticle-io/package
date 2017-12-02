# Metaparticle/Package for Javascript Tutorial
This is an in-depth tutorial for using Metaparticle/Package for Javascript

For a quick summary, please see the [README](README.md).

## Initial Setup

### Check the tools
The `docker` command line tool needs to be installed and working. Try:
`docker ps` to verify this.  Go to the [install page](https://get.docker.io) if you need
to install Docker.

The `mp-compiler` command line tool needs to be installed and working.
Try `mp-compiler --help` to verify this. Go to the [releases page](https://github.com/metaparticle-io/metaparticle-ast/releases) if you need to install
the Metaparticle compiler.

### Get the code
```sh
$ git clone https://github.com/metaparticle-io/package
$ cd package/tutorials/javascript/
# [optional, substitute your favorite editor here...]
$ code .
```

## Initial Program
Inside of the `tutorials/javascript` directory, you will find a simple node project.

You can build this project with `npm run`.

The initial code is a very simple "Hello World"

```javascript
const http = require('http');
const os = require('os');
const port = 8080;

const server = http.createServer((request, response) => {
	console.log(request.url);
	response.end(`Hello World: hostname: ${os.hostname()}\n`);
});

server.listen(port, (err) => {
	if (err) {
		return console.log('server startup error: ', err);
	}
	console.log(`server up on ${port}`);
);
```

You can run this with `npm start`.

## Step One: Containerize the Application
To build a container from our simple application we need to add a dependency to our
build file, and then update the code.

Run:
```sh
npm install -s @metaparticle/package
```

Then update the code to read as follows:

```javascript
const http = require('http');
const os = require('os');
const mp = require('@metaparticle/package');

const port = 8080;

const server = http.createServer((request, response) => {
	console.log(request.url);
	response.end(`Hello World: hostname: ${os.hostname()}\n`);
});

mp.containerize(
	{
		repository: 'docker.io/docker-user-goes-here',
	},
	() => {
		server.listen(port, (err) => {
			if (err) {
				return console.log('server startup error: ', err);
			}
			console.log(`server up on ${port}`);
		});
	}
);
```

You will notice that we added a `mp.containizerize(...)` function.
This function takes two arguments, an object that describes how
to package the application. You will need to replace `your-docker-user-goes-here`
with an actual Docker repository path.

The `containerize` also takes a function to execute inside the container, in this case
the web server.

Once you have this, you can run the program with:

```sh
npm start
```

This code will start your web server again. But this time, it is running
inside a container. You can see this by running:

```sh
docker ps
```

## Step Two: Exposing the ports
If you try to access the web server on [http://localhost:8080](http://localhost:8080) you
will see that you can not actually access the server. Despite it running, the service
is not exposed. To do this, you need to add a `public: true` and supply the ports to expose.

The code snippet to add is:

```javascript
...
mp.containerize(
	{
                ...
		ports: [8080],
        	public: true
	},...
```

This tells the runtime the port(s) to expose. The complete code looks like:

```javascript
const http = require('http');
const os = require('os');
const mp = require('@metaparticle/package');

const port = 8080;

const server = http.createServer((request, response) => {
	console.log(request.url);
	response.end(`Hello World: hostname: ${os.hostname()}\n`);
});

mp.containerize(
	{
		ports: [8080],
		repository: 'docker.io/docker-user-goes-here',
		publish: true,
		public: true
	},
	() => {
		server.listen(port, (err) => {
			if (err) {
				return console.log('server startup error: ', err);
			}
			console.log(`server up on ${port}`);
		});
	}
);
```

Now if you run this with `npm run` your webserver will be successfully exposed on port 8080.

## Replicating and exposing on the web.
As a final step, consider the task of exposing a replicated service on the internet.
To do this, we're going to expand our usage of the options object. First we will
add a `replicas` field, which will specify the number of replicas. Second we will
set our execution environment to `metaparticle` which will launch the service
into the currently configured Kubernetes environment.

Here's what the snippet looks like:

```javascript
mp.containerize(
        {
        	ports: [8080],
        	replicas: 4,
                runner: 'metaparticle',
                repository: 'docker.io/docker-user-goes-here',
                publish: true,
                public: true
        },
        ...);
...
```

And the complete code looks like:


```javascript
const http = require('http');
const os = require('os');
const mp = require('@metaparticle/package');

const port = 8080;

const server = http.createServer((request, response) => {
	console.log(request.url);
	response.end(`Hello World: hostname: ${os.hostname()}\n`);
});

mp.containerize(
	{
        	ports: [8080],
        	replicas: 4,
		runner: 'metaparticle',
		repository: 'docker.io/docker-user-goes-here',
		publish: true,
		public: true
	},
	() => {
		server.listen(port, (err) => {
			if (err) {
				return console.log('server startup error: ', err);
			}
			console.log(`server up on ${port}`);
		});
	}
);
...
```

After you compile and run this, you can see that there are four replicas running behind a
Kubernetes Service Load balancer:

```sh
$ kubectl get pods
...
$ kubectl get services
...
```

Still looking for more? Continue on to the more advanced [sharding tutorial](sharding-tutorial.md)
