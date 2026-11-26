Compose creates a network for the project. Services on that network can address each other by service name and container port. Inside the api container, localhost points back to the api container itself. The worker accepts traffic on port 8081 within the project.

## Commands and operations

Think about where each request starts. The API container's `localhost` belongs to that container. The worker is a separate process in a separate service, so the API must use the worker's Compose name and the worker's internal port.

- Read `WORKER_URL` in the API service's environment. Use the service name listed in `compose.yaml`. A host-published port matters to calls from the host, not to this service-to-service request.
- Use `docker compose up -d --build` after editing configuration so the running services use the updated settings. `docker compose ps` lists the project state and the API health check.
- `docker compose logs api` shows application output and recent failures. If the API remains unhealthy, check whether its URL points inside the wrong container, and verify the worker is running.

Once the API is healthy, use its published host mapping with `curl` and save the response body. The worker may be reachable on the Compose network without publishing its port to the host. Run the lab check while both services are up, then use `docker compose down` for cleanup.
