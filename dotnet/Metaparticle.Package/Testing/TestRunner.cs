using System.IO;

namespace Metaparticle.Package.Testing
{
    public interface TestRunner 
    {
        bool Run(string[] tests);
    }
}