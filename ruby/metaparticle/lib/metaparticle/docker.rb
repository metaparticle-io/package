module Metaparticle
  class Docker
    def in_docker_container?
      if ENV['METAPARTICLE_IN_CONTAINER'] == "true"
        return true
      end

      # Using same hack to work on macOS
      begin
        info = File.readlines('/proc/1/cgroup')

        # horribly ineffient, can do this better
        if !info.select {|line| line =~ /docker/}.empty?
          return true
        end
        if !info.select {|line| line =~ /kubepods/}.empty?
          return true
        end
      rescue
        return false
      end

      return false
    end
  end
end
