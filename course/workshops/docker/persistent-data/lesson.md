A container's writable layer follows that container. A named volume has a separate lifecycle. Mount a named volume to store data you want to keep when replacing the container.

`docker compose up -d --force-recreate` replaces existing containers while reattaching their named volumes. `docker compose down` removes containers and the network and keeps the named volume. Use `docker compose down -v` only when you are finished with this disposable lab data.

## Commands and operations

The named volume is mounted at `/data` inside the API container. A write under `/data` goes to that volume, where the next container for this service can read it. A write into an ordinary unmounted path stays in the current container's filesystem.

- `docker compose exec -T api sh -c '...'` starts a command inside the running API service. The quoted command belongs to the inner shell. If you leave the redirect outside the quotes, your host shell will try to open the file instead.
- Save the old container ID using `docker compose ps -q api`. The quiet output gives the ID with no status table. Use the lab's requested filename so the checker can compare it with the replacement.
- `docker compose up -d --force-recreate` replaces the container and reattaches the volume. Wait for health, then read the new host mapping with `docker compose port api 8080`.

Verify both the changed container ID and the preserved response. The checker uses the recorded host port and live service. Keep the service running until grading finishes. Use `docker compose down -v` only when you intend to discard the volume data as well.
