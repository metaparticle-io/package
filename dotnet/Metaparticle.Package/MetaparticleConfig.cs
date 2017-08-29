using System;

namespace MetaparticlePackage {
    public class MetaparticleConfig : Attribute {
        public bool Verbose { get; set; }

        public string Repository { get; set; }

        public string Version { get; set; }

        public MetaparticleConfig() {
        }
    }
}
