import os
from http.server import BaseHTTPRequestHandler, HTTPServer
from urllib.request import urlopen
class Proxy(BaseHTTPRequestHandler):
    def do_GET(self):
        try:
            with urlopen(os.environ["WORKER_URL"], timeout=2) as response:
                body = response.read()
        except OSError:
            self.send_error(503)
            return
        self.send_response(200)
        self.end_headers()
        self.wfile.write(body)
HTTPServer(("0.0.0.0", 8080), Proxy).serve_forever()
