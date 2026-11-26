An image contains application files and runtime dependencies. A container is the running instance. Compose names services and manages the project together.

This project asks Docker to allocate a free loopback host port for container port 8080. `docker compose port api 8080` prints the mapping. Use that host port with curl. When finished, run `docker compose down` from the lab to clean up the project's container and network.

## Commands and operations

Inspect the files in the lab before launching the service. The `Dockerfile` defines the image build, while `compose.yaml` describes how to run the service, publish its port, and check its health.

- `docker compose up --build -d` builds images where needed and starts the project in the background. The `-d` flag returns the prompt. A started container can still be initializing when the command returns.
- `docker compose ps` displays the project’s container state and health. Use `docker compose logs api` if the service exits or remains unhealthy.
- `docker compose port api 8080` asks Compose which host address and port map to container port 8080. Use that address with `curl` from your machine. The inside port names the application listener. The outside port is where your host connects.

Save just the response body to the requested file while the service is healthy. Check the lab before cleanup, then use `docker compose down` inside the same workspace to remove this project’s runtime resources.
