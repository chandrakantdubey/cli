A health failure can come from configuration, process startup, or the request path. Start with the incident evidence. Compare the release environment with the expected configuration. Restart the service after changing the environment and check the actual response.

This exercise provides a disposable API, a release environment file, and an incident log. Run commands in the lab shell. Leave the app running until it has been graded and stop it afterward.

## Commands and operations

Start from the symptom in the incident log. Identify whether the process is reachable and which release mode it started in. Then compare the release file with the runbook.

- `grep '^APP_MODE=' release.env` reads the exact setting. Update the release configuration so the next start will use the corrected mode. Changing a variable after the process has started cannot update the running process.
- Load the provided release environment in the lab shell before starting the API, or explicitly pass the corrected value on startup. If you source a configuration as shell code, use only the provided file after inspecting it and export the variable for the child process.
- Read the current listener's port from `server.port`. Use `curl -i` to see the HTTP status while diagnosing, and use `curl -fsS` when saving only the response body.

A valid recovery has two pieces of evidence. The configuration records the intended mode for future starts, and the running service responds with HTTP success. Leave the corrected API running during `cli-silly check`, then stop the background process from the lab shell.
