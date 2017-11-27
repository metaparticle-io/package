const http = require('http');
const os = require('os');
const mp = require('@metaparticle/package');

const port = 8080;

const server = http.createServer((request, response) => {
	console.log(request.url);
	response.end(`path: ${request.url}, host: ${os.hostname()}\n`);
});

mp.containerize(
	{
		ports: [8080],
		shardSpec: {
			shards: 3,
			"urlPattern": "(.*)"
		},
		runner: 'metaparticle',
		repository: 'brendanburns',
		publish: true,
		public: true
	},
	() => {
		server.listen(port, (err) => {
			if (err) {
				return console.log('server startup error: ', err);
			}
			console.log(`server up on ${port}`);
		});
	}
);
