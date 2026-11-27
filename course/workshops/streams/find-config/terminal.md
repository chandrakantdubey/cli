## Try it in the terminal

### Separate which file from which line

Try `find config -type f -name '*.conf'`. If you remove the quotes and Bash expands a matching glob, `find` may receive an unexpected operand. Keep the pattern quoted and use the output to identify the intended search scope.

Compare `grep -rn 'debug' config` with a whole-line search. Look for comments or similar values in the broad search. Use `-l` when the report needs just a path. Use `grep` on that one file without `-l` for the other report. The example outside `config/` is deliberately irrelevant to this search.
