import os
from http.server import BaseHTTPRequestHandler, HTTPServer
from pathlib import Path

mode = os.environ.get("APP_MODE", "maintenance")
healthy = mode == "ready"
body = b"service: ready\n"
class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path != "/health":
            self.send_error(404)
            return
        self.send_response(200 if healthy else 503)
        self.end_headers()
        self.wfile.write(body if healthy else b"service: maintenance\n")

with HTTPServer(("127.0.0.1", 0), Handler) as server:
    Path("server.port").write_text(f"{server.server_port}\n")
    server.serve_forever()
