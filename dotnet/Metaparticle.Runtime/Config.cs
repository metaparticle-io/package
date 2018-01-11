namespace Metaparticle.Runtime {
    public class Config : System.Attribute {
        public int Replicas { get; set; }

        public int Shards { get; set; }

        public int JobCount { get; set; }

        public string ShardExpression { get; set; }

        public string Executor { get; set; }

        public int[] Ports { get; set; }

        public bool Public { get; set; }

        public Config() {
            Executor = "docker";
            Replicas = 1;
        }
    }
}