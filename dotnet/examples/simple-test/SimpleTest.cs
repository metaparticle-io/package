using System;
using Xunit;
using simple;

namespace simpletest
{
    public class SimpleTests
    {
        [Fact]
        public void GivenTwoStrings_WhenConcatenateIsCalled_TwoStringsAreConcatenated()
        {
            // arrange
            var stringOne = "Hello";
            var stringTwo = "World";
            
            // act
            var result = Program.Concatenate(stringOne, stringTwo);

            // assert
            Assert.Equal("HelloWorld", result);
        }
    }
}
