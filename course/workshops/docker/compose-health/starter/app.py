from http.server import BaseHTTPRequestHandler, HTTPServer
from pathlib import Path

class API(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path != "/health":
            self.send_error(404)
            return
        self.send_response(200)
        self.end_headers()
        self.wfile.write(b'container: ready\n')

HTTPServer(("0.0.0.0", 8080), API).serve_forever()
