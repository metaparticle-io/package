# coding: utf-8
lib = File.expand_path("../lib", __FILE__)
$LOAD_PATH.unshift(lib) unless $LOAD_PATH.include?(lib)
require "metaparticle/version"

Gem::Specification.new do |spec|
  spec.name          = "metaparticle"
  spec.version       = Metaparticle::VERSION
  spec.authors       = ["Christopher Hein"]
  spec.email         = ["me@christopherhein.com"]

  spec.summary       = %q{Allows you to include containerization and deployment hooks into your ruby applications}
  spec.description   = %q{Allows you to include containerization and deployment hooks into your ruby applications}
  spec.homepage      = "https://metaparticle.io/"
  spec.license       = "MIT"

  spec.files = Dir['lib/**/*.rb']
  spec.bindir        = "exe"
  spec.executables   = spec.files.grep(%r{^exe/}) { |f| File.basename(f) }
  spec.require_paths = ["lib"]

  spec.add_development_dependency "bundler", "~> 1.15"
  spec.add_development_dependency "rake", "~> 10.0"
  spec.add_development_dependency "rspec", "~> 3.0"
end
