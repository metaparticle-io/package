using System;

namespace Metaparticle.Tests
{
    public class Config : Attribute
    {
        public string[] Names { get;set; }
    }
}