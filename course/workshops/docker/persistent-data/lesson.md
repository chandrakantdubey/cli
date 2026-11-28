A container's writable layer follows that container. A named volume has a separate lifecycle. Mount a named volume to store data you want to keep when replacing the container.

`docker compose up -d --force-recreate` replaces existing containers while reattaching their named volumes. `docker compose down` removes containers and the network and keeps the named volume. Use `docker compose down -v` only when you are finished with this disposable lab data.

## Commands and operations

The named volume is mounted at `/data` inside the API container. A write under `/data` goes to that volume, where the next container for this service can read it. A write into an ordinary unmounted path stays in the current container's filesystem.

- `docker compose exec -T api sh -c '...'` starts a command inside the running API service. The quoted command belongs to the inner shell. If you leave the redirect outside the quotes, your host shell will try to open the file instead.
- Save the old container ID using `docker compose ps -q api`. The quiet output gives the ID with no status table. Use the lab's requested filename so the checker can compare it with the replacement.
- `docker compose up -d --force-recreate` replaces the container and reattaches the volume. Wait for health, then read the new host mapping with `docker compose port api 8080`.

Verify both the changed container ID and the preserved response. The checker uses the recorded host port and live service. Keep the service running until grading finishes. Use `docker compose down -v` only when you intend to discard the volume data as well.

## Before you practice

**Outcome:** Prove that stored data survives container replacement.

Complete Compose health first and understand containers and volumes.

## Common mistakes and recovery

Recreating a container is different from deleting its volume. docker compose down -v removes the exercise data.

Inspect the relevant input and output before repeating a command. Use `cli-silly hint` for a starting point and `cli-silly check` for criterion-by-criterion feedback. If the starter is damaged, `cli-silly reset` archives this attempt before recreating it.

## Use it in development

Verify persistence before relying on a volume for application state.

After passing, explain which files or state changed and how you verified them. Open related practice from **Do** to apply the idea to different data.
