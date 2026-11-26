An image contains application files and runtime dependencies. A container is the running instance. Compose names services and manages the project together.

This project asks Docker to allocate a free loopback host port for container port 8080. `docker compose port api 8080` prints the mapping. Use that host port with curl. When finished, run `docker compose down` from the lab to clean up the project's container and network.
