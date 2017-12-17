module Metaparticle
  class DockerBuilder
    def initialize(config)
      @config = config
    end

    def build
      `docker build -t #{@config.image} .`
    end

    def push
      `docker push #{@config.image}`
    end
  end
end
