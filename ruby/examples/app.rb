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
    replicas: 3,
    runner: 'metaparticle',
    repository: 'christopherhein',
    publish: true,
    public: true,
		# shardSpec: {
		# 	shards: 3,
		# 	"urlPattern": "(.*)"
		# },
  }) do
  App.run!
end
