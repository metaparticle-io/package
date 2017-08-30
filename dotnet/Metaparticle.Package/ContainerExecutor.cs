using System.IO;
namespace Metaparticle.Package {
    public interface ContainerExecutor {
        string Run(string image);

        void Cancel(string id);

        void Logs(string id, TextWriter stdout, TextWriter stderr);
    }
}