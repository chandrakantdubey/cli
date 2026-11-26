from http.server import BaseHTTPRequestHandler, HTTPServer
class Worker(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.end_headers()
        self.wfile.write(b"worker: online\n")
HTTPServer(("0.0.0.0", 8081), Worker).serve_forever()
