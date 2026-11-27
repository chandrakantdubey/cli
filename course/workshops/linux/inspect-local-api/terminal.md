## Try it in the terminal

### Locate the failure before writing a report

Read the top of `app.py` and find the setting it looks up for readiness. In one practice shell, compare the status before and after changing that setting. Identify the listener address it chose from the output and port file. Use `curl -i` on its local `/health` route to see response headers and the status as well as the body.

For the real report capture body output only, without headers. Verify the report with `cat` and leave the lab service running until the check finishes. If the expected file is correct and the score still reports a failed live check, confirm the background job and its current port. Use `jobs -l` from the same shell and stop your own job when you finish.
