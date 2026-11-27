A configuration reader uses two environment variables. Compare the result under default, preview, and release settings so a developer can tell which value goes into each process.

## Commands and operations

`env -u REGION bash sample.sh` runs the child command without that inherited variable. In Bash, `${REGION:-local}` expands to local if the variable is unset or empty. The same technique lets the lab show defaults regardless of your shell's environment.

An assignment directly before a command, such as `REGION=eu bash sample.sh`, applies to that command. An `export REGION=eu` instead affects future commands in this shell until the setting is changed. A subshell groups temporary exports and discards them when it exits.

Try `env -u LOG_LEVEL APP_MODE=release bash bin/show-config.sh` to see how removing one setting and providing another work together. Inspect the outputs as a pair of settings in each report.
