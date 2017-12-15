# Metaparticle/Package for Ruby Tutorial
This is an in-depth tutorial for using Metaparticle/Package for Ruby

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
$ cd package/tutorials/ruby/
# [optional, substitute your favorite editor here...]
$ code .
```

## Initial Program
Inside of the `tutorials/ruby` directory, you will find a simple sinatra project.

You can build this project with `ruby app.rb`.

The initial code is a very simple "Hello World"

```ruby
require 'rubygems'
require 'sinatra/base'

class App < Sinatra::Base
  set :bind, '0.0.0.0'

  get '/' do
    'Hello World'
  end
  run!
end

```

You can run this with `ruby app.rb`.

## Step One: Containerize the Application
To build a container from our simple application we need to add a dependency to our
build file, and then update the code.

Run:
```sh
gem install metaparticle
```

Then update the code to read as follows:

```ruby
require 'rubygems'
require 'metaparticle'
require 'sinatra/base'

Metaparticle::Package.containerize(
  {
    ports: [4567],
    name: 'sinatra-app',
    repository: 'christopherhein',
  }) do
  class App < Sinatra::Base
    set :bind, '0.0.0.0'

    get '/' do
      'metaparticle from ruby!'
    end
    run!
  end
end
```

You will notice that we added a `Metaparticle::Package.containerize` class.
This class takes two arguments, an object that describes how
to package the application. You will need to replace `your-docker-user-goes-here`
with an actual Docker repository path.

The `Metaparticle::Package.containerize` also takes a function to execute inside the container, in this case
the web server.

Once you have this, you can run the program with:

```sh
ruby app.rb
```

This code will start your web server again. But this time, it is running
inside a container. You can see this by running:

```sh
docker ps
```

## Replicating and exposing on the web.
As a final step, consider the task of exposing a replicated service on the internet.
To do this, we're going to expand our usage of the options object. First we will
add a `replicas` field, which will specify the number of replicas. Second we will
set our execution environment to `metaparticle` which will launch the service
into the currently configured Kubernetes environment.

Here's what the snippet looks like:

```ruby
Metaparticle::Package.containerize(
  {
    name: 'sinatra-app',
    ports: [4567],
    replicas: 4,
    runner: 'metaparticle',
    repository: 'christopherhein',
    publish: true,
    public: true
  }) do
  ...
end
```

And the complete code looks like:

```ruby
require 'rubygems'
require 'metaparticle'
require 'sinatra/base'

class App < Sinatra::Base
  set :bind, '0.0.0.0'

  get '/' do
    'metaparticle from ruby!'
  end
end

Metaparticle::Package.containerize(
  {
    name: 'sinatra-app',
    ports: [4567],
    replicas: 4,
    runner: 'metaparticle',
    repository: 'christopherhein',
    publish: true,
    public: true
  }) do
  App.run!
end
```

After you compile and run this, you can see that there are four replicas running behind a
Kubernetes Service Load balancer:

```sh
$ kubectl get pods
...
$ kubectl get services
...
```
