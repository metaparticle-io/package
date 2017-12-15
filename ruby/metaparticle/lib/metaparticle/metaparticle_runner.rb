require 'fileutils'

module Metaparticle
  class MetaparticleRunner
    def initialize(config)
      @config = config
    end

    def run
      options = {
        name: @config.name,
        guid: 1234567,
        services: [
          {
            name: @config.name,
            replicas: @config.replicas,
            shardSpec: @config.shardSpec,
            containers: [
              {
                image: @config.image
              }
            ],
            ports: @config.ports.map {|p| {number: p}}
          }
        ],
        serve: {
          name: @config.name
        }
      }
      if @config.public
        options[:serve][:public] = true
      end

      create_directory

      File.open(".metaparticle/service.json", "w") do |f|
        f.write(options.compact!.to_json)
      end

      `mp-compiler -f .metaparticle/service.json`
      `mp-compiler -f .metaparticle/service.json --deploy=false --attach=true`
    end

    private
    def create_directory
      FileUtils::mkdir_p('.metaparticle')
    end
  end
end
