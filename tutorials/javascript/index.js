const http = require('http');
const os = require('os');
const port = 8080;

const server = http.createServer((request, response) => {
	console.log(request.url);
	response.end(`Hello World: hostname: ${os.hostname()}\n`);
});

server.listen(port, (err) => {
	if (err) {
		return console.log('server startup error: ', err);
	}
	console.log(`server up on ${port}`);
});
