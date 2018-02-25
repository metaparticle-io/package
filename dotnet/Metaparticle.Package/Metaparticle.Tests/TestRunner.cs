using System.IO;

namespace Metaparticle.Tests
{
    public interface TestRunner 
    {
        bool Run(string[] tests);
    }
}