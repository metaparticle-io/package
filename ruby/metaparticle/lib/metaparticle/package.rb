require "ostruct"
require "metaparticle/docker"
require "metaparticle/docker_builder"
require "metaparticle/docker_runner"
require "metaparticle/metaparticle_runner"

module Metaparticle
  class Package
    def self.containerize(config, &block)
      new(config, &block)
    end

    def initialize(config, &block)
      @config = OpenStruct.new(config)
      @app = block
      run!
    end

    private
    def write_dockerfile
      dockerfile = <<-DOCKERFILE
FROM ruby:2.4-alpine

RUN mkdir -p /#{@config.name}
WORKDIR /#{@config.name}

COPY Gemfile /#{@config.name}/
COPY Gemfile.lock /#{@config.name}/
RUN bundle install

COPY . /#{@config.name}/
CMD ruby app.rb
      DOCKERFILE
      File.open("Dockerfile", "w") do |f|
        f.write(dockerfile)
      end
    end

    def docker
      @docker ||= Docker.new
    end

    def builder
      klass = @config.builder || 'docker'
      @builder ||= classify("#{klass}_builder").new(@config)
    end

    def runner
      klass = @config.runner || 'docker'
      @runner ||= classify("#{klass}_runner").new(@config)
    end

    def classify(name)
      Object.const_get("Metaparticle::"+name.split('_').collect!{ |w| w.capitalize }.join)
    end

    def run!
      if docker.in_docker_container?
        @app.call
      end
      @config.image = "#{@config.repository}/#{@config.name}"
      write_dockerfile
      builder.build
      builder.push
      runner.run
      return true
    end
  end
end
