# Metaparticle/Package for Ruby
Metaparticle/Package is a collection of libraries intended to 
make building and deploying containers a seamless and idiomatic
experience for developers.

This is the implementation for Ruby.

## Introduction
Metaparticle/Package simplifies and centralizes the task of
building and deploying a container image.

Here is a quick example.

Consider this simple Ruby application:

```ruby
class App < Sinatra::Base
  set :bind, '0.0.0.0'

  get '/' do
    'metaparticle from ruby!'
  end
end
```

To containerize this application, you need to use the `metaparticle` gem and
the `Metaparticle::Package` wrapper class like this:

```ruby
require 'metaparticle'
require 'sinatra/base'

Metaparticle::Package.containerize(
  {
    name: 'sinatra-app',
    ports: [4567],
    replicas: 3,
    runner: 'metaparticle',
    repository: 'christopherhein',
    publish: true,
    public: true
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

Then when you run the ruby file it will automatically create a docker image,
push to your registry of choice, and then use `mp-compiler` to push to your k8s
cluster.

## Tutorial
For a more complete exploration of the Metaparticle/Package for Ruby, please see the [in-depth tutorial](../tutorials/ruby/tutorial.md).

