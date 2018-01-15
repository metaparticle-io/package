## Metaparticle API 

> see https://aka.ms/autorest 



## Getting Started 

To build the SDKs for My API, simply install AutoRest via `npm` (`npm install -g autorest`) and then run:

> `autorest readme.md`



To see additional help and options, run:

> `autorest --help`



For other options on installation see [Installing AutoRest](https://aka.ms/autorest/install) on the AutoRest github page.



---



## Configuration 

The following are the settings for this using this API with AutoRest.



``` yaml
input-file:
- /mnt/c/Users/bburns/gopath/src/github.com/metaparticle-io/metaparticle-ast/api.yaml

output-folder: ./Metaparticle.NET
csharp: # just having a 'csharp' node enables the use of the csharp generator.
  namespace: Microsoft.MyApp.MyNameSpace #override the namespace 
  output-folder : generated/csharp # relative to the global value.

log-file: ./logs.txt

```
