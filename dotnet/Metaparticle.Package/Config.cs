using System;

namespace Metaparticle.Package
{
    public class Config : Attribute
    {
        public bool Verbose { get; set; }

        public bool Quiet { get; set; }

        public string Repository { get; set; }

        public string Version { get; set; }

        public string Builder { get; set; }

        public bool Publish { get; set; }

        public string Dockerfile { get; set; }

        public Config()
        {
            Builder = "docker";
            LoadConfigVariablesFromEnvironment();
        }

        private void LoadConfigVariablesFromEnvironment()
        {
            Repository = TryGetEnvironmentVariable("METAPARTICLE_CONFIG_REPOSITORY");
            Publish = TryGetEnvironmentVariable("METAPARTICLE_CONFIG_PUBLISH").ToLower() == "true";
        }

        private string TryGetEnvironmentVariable(string name)
        {
            return Environment.GetEnvironmentVariable(name) ?? "";
        }
    }
}
