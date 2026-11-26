A container's writable layer follows that container. A named volume has a separate lifecycle. Mount a named volume to store data you want to keep when replacing the container.

`docker compose up -d --force-recreate` replaces existing containers while reattaching their named volumes. `docker compose down` removes containers and the network and keeps the named volume. Use `docker compose down -v` only when you are finished with this disposable lab data.
