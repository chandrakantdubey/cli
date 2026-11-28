Build a tiny Node.js command that counts non-empty lines in a text file. It works as a pipeline component because it keeps data, diagnostics, and completion status separate. The starter uses Node's built-in modules and works without npm packages.

## Commands and operations

`process.argv.slice(2)` gives the user arguments without the Node executable and script path. Handle `--help` first. Write its usage line to stdout and let Node exit with status 0. For bad arguments write the usage line to stderr and set `process.exitCode = 2`.

Use `readFile(file, "utf8")` from `node:fs/promises` to read a string. Split on line endings and keep lines whose `trim()` result is nonempty. Test a file with blank lines and one without a final newline. A path containing spaces is one argument when quoted in Bash.

On a read failure print `error: cannot read input` on stderr and set `process.exitCode = 1`. Try `node line-count.mjs missing.txt > count.out 2> count.err` and immediately inspect `$?`. The counter belongs on stdout only. **Grading runs your Node.js file with fresh fixtures in a temporary directory under your user account.**

## Learning goals

By the end of this module, you should be able to:

- design a CLI contract around arguments, stdin, stdout, stderr, and exit codes;
- parse input safely instead of relying on fragile string handling;
- structure a Node CLI so core logic is testable;
- handle filesystem paths, signals, and errors predictably;
- package a command for local or team use.

## Project shape

A small CLI can start with npm and native ESM:

```bash
mkdir my-cli
cd my-cli
npm init -y
mkdir src
```

A practical `package.json` includes scripts, an entry point, and a `bin` command name:

```json
{
  "type": "module",
  "bin": {
    "my-cli": "./src/index.js"
  },
  "scripts": {
    "test": "node --test",
    "start": "node src/index.js"
  }
}
```

Make the entry file executable when it will be run directly:

```bash
chmod +x src/index.js
```

The first line should identify the interpreter:

```js
#!/usr/bin/env node
```

## Arguments and help

Node exposes raw arguments through `process.argv`, but raw parsing becomes error-prone quickly:

```js
const args = process.argv.slice(2);
```

For real tools, use a parser such as Commander, Yargs, or a small local parser for narrow needs. A good CLI documents required arguments, optional flags, defaults, and examples through `--help`.

## stdin, stdout, and stderr

A Unix-friendly CLI separates data from diagnostics:

```js
process.stdout.write(JSON.stringify(result) + "\n");
process.stderr.write("warning: skipped invalid row\n");
```

Read stdin when no file is provided:

```js
async function readStdin() {
  process.stdin.setEncoding("utf8");
  let input = "";
  for await (const chunk of process.stdin) input += chunk;
  return input;
}
```

Use stdout for output that another command might consume. Use stderr for progress, warnings, and error messages.

## Exit codes

Use zero for success and non-zero for failure:

```js
process.exitCode = 1;
```

Prefer setting `process.exitCode` over calling `process.exit()` so pending writes and cleanup can finish. Common conventions are `1` for general failure, `2` for usage errors, and `130` for Ctrl-C interruption.

## Files and paths

Use standard modules instead of manual path concatenation:

```js
import { readFile } from "node:fs/promises";
import path from "node:path";

const fullPath = path.resolve(inputPath);
const text = await readFile(fullPath, "utf8");
```

Do not assume the current working directory is the script directory. `process.cwd()` means where the user ran the command; `import.meta.url` identifies the module location.

## Errors users can act on

Catch expected errors near the CLI boundary:

```js
try {
  const text = await readFile(file, "utf8");
  process.stdout.write(summarize(text));
} catch (error) {
  process.stderr.write(`error: ${error.message}\n`);
  process.exitCode = 1;
}
```

Keep stack traces for debugging modes, not normal user mistakes like missing files or malformed JSON.

## Signals and cleanup

Long-running tools should handle interruption when they own temporary files, child processes, or network connections:

```js
process.on("SIGINT", () => {
  process.stderr.write("interrupted\n");
  process.exitCode = 130;
});
```

If cleanup is asynchronous, design it explicitly. Signal handlers are not a substitute for careful resource ownership.

## Testing a CLI

Test the core function directly, then test the command boundary with child processes:

```js
import { spawnSync } from "node:child_process";
import test from "node:test";
import assert from "node:assert/strict";

test("prints help", () => {
  const result = spawnSync("node", ["src/index.js", "--help"], { encoding: "utf8" });
  assert.equal(result.status, 0);
  assert.match(result.stdout, /Usage/);
});
```

Cover success, invalid arguments, missing files, malformed input, stdout, stderr, and exit status.

## Distribution checklist

Before sharing a Node CLI, verify:

- `bin` points to an executable file with a shebang;
- dependency versions are locked;
- `--help` is accurate;
- output is stable enough for scripts;
- errors include enough context without leaking secrets;
- tests run in a clean checkout.

## Before you practice

**Outcome:** Build a Node command with testable help, output, and failure behavior.

Understand arguments, stdout, stderr, and exit status. Node.js is required.

## Common mistakes and recovery

Do not mix help or errors into a successful machine-readable result. Test missing files, empty inputs, and filenames with spaces.

Inspect the relevant input and output before repeating a command. Use `cli-silly hint` for a starting point and `cli-silly check` for criterion-by-criterion feedback. If the starter is damaged, `cli-silly reset` archives this attempt before recreating it.

## Use it in development

Build a small repository utility used locally and in automation.

After passing, explain which files or state changed and how you verified them. Open related practice from **Do** to apply the idea to different data.
