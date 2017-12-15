require 'rubygems'
require 'sinatra/base'

class App < Sinatra::Base
  set :bind, '0.0.0.0'

  get '/' do
    'Hello World'
  end
  run!
end

