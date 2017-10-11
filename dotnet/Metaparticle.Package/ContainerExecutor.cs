using System.IO;
namespace Metaparticle.Package {
    public interface ContainerExecutor {
        string Run(string image, Metaparticle.Runtime.Config config);

        void Cancel(string id);

        void Logs(string id, TextWriter stdout, TextWriter stderr);

        bool PublishRequired();
    }
}