#!/usr/bin/python

from six.moves import SimpleHTTPServer, socketserver
import socket
from metaparticle_pkg.containerize import Containerize

OK = 200

port = 8080

class MyHandler(SimpleHTTPServer.SimpleHTTPRequestHandler):
    def do_GET(self):
        self.send_response(OK)
        self.send_header("Content-type", "text/plain")
        self.end_headers()
        self.wfile.write("Hello Metparticle [{}] @ {}\n".format(self.path, socket.gethostname()).encode('UTF-8'))
        print("request for {}".format(self.path))
    def do_HEAD(self):
        self.send_response(OK)
        self.send_header("Content-type", "text/plain")
        self.end_headers()

@Containerize(
    package={'name': 'web', 'repository': 'docker.io/brendanburns'},
    runtime={'ports': [8080], 'executor': 'metaparticle', 'replicas': 3}
)
def main():
    Handler = MyHandler
    httpd = socketserver.TCPServer(("", port), Handler)
    httpd.serve_forever()

if __name__ == '__main__':
    main()
