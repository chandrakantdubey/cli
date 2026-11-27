## Try it in the terminal

### Trace two relative paths at once

From the workspace root, use `ls services` to see the siblings. Run `cd services/api` and `pwd`, then `ls ../worker`. For the input path `../worker/owner.txt`, each `..` moves the lookup up one level before the shell passes the path to the command.

The input and the output of one command can have different prefixes. From `services/api`, `env.txt` names the nearby input and `../../api-env.txt` names a report two directory levels up. Before writing the real report, run `ls ../..` to confirm where it will land. Run `cli-silly task` from the subdirectory to see that it still knows this lesson.
