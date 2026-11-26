from http.server import BaseHTTPRequestHandler, HTTPServer
from pathlib import Path

class API(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path != "/health":
            self.send_error(404)
            return
        self.send_response(200)
        self.end_headers()
        self.wfile.write(Path("/data/status.txt").read_bytes() if Path("/data/status.txt").exists() else b"pending\n")

HTTPServer(("0.0.0.0", 8080), API).serve_forever()
