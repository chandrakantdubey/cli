## Try it in the terminal

Run `env -u APP_MODE -u LOG_LEVEL bash bin/show-config.sh` from a shell where either variable might already exist. Compare with a one-off temporary value. There is no need to edit .bashrc.

Use `diff -u default.out preview.out` after creating the two reports. An exit status of 1 means the files differ. Compare release and default next and identify the single expected change before running the grader.
