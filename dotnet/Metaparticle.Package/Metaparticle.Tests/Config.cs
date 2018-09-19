using System;

namespace Metaparticle.Tests
{
    public class Config : Attribute
    {
		/// <summary>
		/// Directories containing the test projects. This may be absolute or relative to the current project.
		/// </summary>
        public string[] Names { get;set; }
    }
}