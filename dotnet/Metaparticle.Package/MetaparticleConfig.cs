using System;

namespace Metaparticle.Package {
    public class MetaparticleConfig : Attribute {
        public bool Verbose { get; set; }

        public string Repository { get; set; }

        public string Version { get; set; }

        public string Executor { get; set; }

        public MetaparticleConfig() {
            Executor = "docker";
        }
    }
}
