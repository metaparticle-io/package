using System;
using System.IO;
using static Metaparticle.Package.Util;

namespace Metaparticle.Package.Testing
{
    public class DotnetTestRunner : TestRunner 
    {
        public bool Run(string[] tests)
        {
            foreach (var testFolder in tests)
            {
                Console.WriteLine($"Running Test {testFolder}");
                var result = Exec("dotnet", $"test {testFolder}", Console.Out, Console.Error);
                
                if (result.ExitCode != 0)
                    return false;
            }
            
            return true;
        }
    }
}