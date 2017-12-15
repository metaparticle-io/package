module Metaparticle
  class DockerRunner
    def initialize(config)
      @config = config
    end

    def run
      `docker run -it --name #{@config.name} #{port_string} #{@config.image}`
    end

    private
    def port_string
      @port_string = @config.ports.map {|port| "-p #{port}:#{port}" }.join(" ")
    end
  end
end
