Processes inherit exported environment variables from their parent. A variable can also be set for one command with `NAME=value command`. Changing a variable in the shell after a process starts cannot update that process.

This API binds to a free loopback port and writes it to server.port. A shell trailing `&` starts a background job, and `$!` gives its PID. Use `ps -p` to inspect the process and `curl -i` to inspect its HTTP status. Keep the workspace shell open while grading and stop the job with `kill` when you finish.

## Commands and operations

Run the API from the lab directory so Python can find `app.py` and write `server.port` where the checker expects it. `APP_MODE=ready python3 app.py` provides the environment variable for that one process. The trailing `&` gives you back the prompt while the process keeps listening.

- `> server.log 2>&1` sends both output streams to a single file. If the process exits early or `server.port` never appears, inspect the log first.
- Save `$!` as soon as you start the background process, or run `jobs -l` from the same shell to see the job and PID. Use `ps -p PID` to check that the process is still present.
- Build the URL using the port from `server.port`. `curl -i` shows the HTTP status and headers as well as the body. When saving just the body, use `curl -fsS` and redirect standard output. `-f` reports HTTP failures, `-s` hides progress, and `-S` keeps errors visible.

`Connection refused` means you did not reach a listener at that address. A `503` means you did reach the application, but it is in maintenance mode. Keep the process running for `cli-silly check`, then stop the job you started.

## Interfaces and addresses

```bash
ip addr
ip -br addr
ip link
```

## Routing

```bash
ip route
ip route get 1.1.1.1
```

A correct IP address does not guarantee a correct route.

## Connectivity

ICMP:

```bash
ping -c 4 1.1.1.1
```

TCP:

```bash
nc -vz example.com 443
```

Some environments block ICMP, so a failed ping does not prove that the service is unreachable.

## DNS

```bash
dig example.com
getent hosts example.com
```

`getent` follows the system's configured name-service path, while `dig` queries DNS directly.

## HTTP

```bash
curl -I https://example.com
curl -v https://example.com
curl -fsS https://example.com/health
```

Use `-v` when you need connection/TLS details; use `-f` in scripts when HTTP failures should become non-zero exit status.

## Listening sockets

```bash
ss -lntup
```

This answers an important question: is anything actually listening on the expected address and port?

## Ports and addresses

A service bound to `127.0.0.1:8080` is normally reachable only locally. A service bound to `0.0.0.0:8080` can listen on all IPv4 interfaces, subject to firewall and network policy.

## Troubleshooting sequence

```text
DNS -> route -> TCP connection -> TLS -> HTTP/application
```

Do not jump directly to "the network is down."
