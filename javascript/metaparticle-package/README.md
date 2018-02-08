
# Metaparticle/Package for JavaScript
Metaparticle/Package is a collection of libraries intended to 
make building and deploying containers a seamless and idiomatic
experience for developers.

This is the implementation for JavaScript (NodeJS).

## Introduction
Metaparticle/Package simplifies and centralizes the task of
building and deploying a container image.

Here is a quick example.

Consider this simple Javascript application:

```javascript
console.log('hello world!');
```

To containerize this application, you need to use the `@metaparticle/package` library and
the `containerize` wrapper function like this:

```javascript
const mp = require('@metaparticle/package');

mp.containerize(
	{
		repository: 'brendanburns',
	},
	() => {
        console.log('hello world');
	}
);
```

You can then compile this application just as you have before.
But now, when you run the application, instead of printing "Hello world", it first packages itself as a container, and
then (optionally) deploys itself inside that container.

## Tutorial
For a more complete exploration of the Metaparticle/Package for javascript, please see the [in-depth tutorial](../tutorials/javascript/tutorial.md).

### Use ACI

```javascript
const mp = require('@metaparticle/package');

mp.containerize(
    {
		repository: 'myrepo',
		publish: 'true',
		runner: 'aci',
		resourceGroup: 'myResourceGroup',
		subscriptionId: '00000000-0000-0000-0000-000000000000',
		public: 'true',
		ports: [80],
		clientId: '00000000-0000-0000-0000-000000000000',
		clientSecret: 'mykey',
		tenantId: '00000000-0000-0000-0000-000000000000'
    },
    () => {
        console.log('hello world');
    }
);
```
