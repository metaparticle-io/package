#!/usr/bin/python
from metaparticle_pkg import Containerize

import logging
import socket
from six.moves import SimpleHTTPServer, socketserver

# all metaparticle output is accessible through the stdlib logger (debug level)
logging.basicConfig(level=logging.INFO)
logging.getLogger('metaparticle_pkg.runner').setLevel(logging.DEBUG)
logging.getLogger('metaparticle_pkg.builder').setLevel(logging.DEBUG)


OK = 200
port = 8080


class MyHandler(SimpleHTTPServer.SimpleHTTPRequestHandler):
    def do_GET(self):
        self.send_response(OK)
        self.send_header("Content-type", "text/plain")
        self.end_headers()
        self.wfile.write("Hello Metaparticle [{}] @ {}\n".format(self.path, socket.gethostname()).encode('UTF-8'))
        print(("request for {}".format(self.path)))

    def do_HEAD(self):
        self.send_response(OK)
        self.send_header("Content-type", "text/plain")
        self.end_headers()


@Containerize(
    package={
        # to run this example you'll need to change these values
        'name': 'web',
        'repository': 'docker.io/brendanburns',
    },
    runtime={
        'ports': [8080],
        'executor': 'metaparticle',
        'replicas': 3,
        'public': True
    }
)
def main():
    Handler = MyHandler
    httpd = socketserver.TCPServer(("", port), Handler)
    httpd.serve_forever()


if __name__ == '__main__':
    main()
