## Try it in the terminal

### Keep the process identity in view

Start the starter worker with Bash in the background and capture `$!` immediately. Run `jobs -l` in the same practice shell. Compare its job number with the numeric PID accepted by `ps -p`. If you start a second background job, `$!` will change, so keep the first value in a named variable.

Read `worker.sh` and find the trap before sending TERM. Once its readiness file exists, inspect the command line with `ps`, then ask your worker to stop. Wait for that exact child before opening `stopped.txt`. A second `ps` should show that the worker exited. Target only the PID you started in this lab.
