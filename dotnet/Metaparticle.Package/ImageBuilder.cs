using System.IO;

namespace Metaparticle.Package 
{
    public interface ImageBuilder 
    {
        bool Build(string configFile, string imageName, TextWriter stdout = null, TextWriter stderr = null);

        bool Push(string imageName, TextWriter stdout = null, TextWriter stderr = null);
    }
}