A health failure can come from configuration, process startup, or the request path. Start with the incident evidence. Compare the release environment with the expected configuration. Restart the service after changing the environment and check the actual response.

This exercise provides a disposable API, a release environment file, and an incident log. Run commands in the lab shell. Leave the app running until it has been graded and stop it afterward.

## Commands and operations

Start from the symptom in the incident log. Identify whether the process is reachable and which release mode it started in. Then compare the release file with the runbook.

- `grep '^APP_MODE=' release.env` reads the exact setting. Update the release configuration so the next start will use the corrected mode. Changing a variable after the process has started cannot update the running process.
- Load the provided release environment in the lab shell before starting the API, or explicitly pass the corrected value on startup. If you source a configuration as shell code, use only the provided file after inspecting it and export the variable for the child process.
- Read the current listener's port from `server.port`. Use `curl -i` to see the HTTP status while diagnosing, and use `curl -fsS` when saving only the response body.

A valid recovery has two pieces of evidence. The configuration records the intended mode for future starts, and the running service responds with HTTP success. Leave the corrected API running during `cli-silly check`, then stop the background process from the lab shell.

## Before you practice

**Outcome:** Recover a local release using incident evidence, configuration, and a live health check.

Complete log filtering, environment variables, and the local API lab. Python 3 and curl are required.

## Common mistakes and recovery

A healthy saved file is not proof that the current service is healthy. Verify the live endpoint and its current configuration.

Inspect the relevant input and output before repeating a command. Use `cli-silly hint` for a starting point and `cli-silly check` for criterion-by-criterion feedback. If the starter is damaged, `cli-silly reset` archives this attempt before recreating it.

## Use it in development

Practice an incident response that leaves a reproducible configuration and supporting evidence.

After passing, explain which files or state changed and how you verified them. Open related practice from **Do** to apply the idea to different data.
