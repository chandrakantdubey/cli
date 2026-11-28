// Verify the production companion using an installed Chromium-family browser.
// Run: node scripts/test-browser.mjs http://127.0.0.1:4193
import assert from "node:assert/strict";
import { spawn } from "node:child_process";
import { mkdtemp, readFile, readdir, mkdir } from "node:fs/promises";
import { setTimeout as delay } from "node:timers/promises";

const base = process.argv[2] ?? "http://127.0.0.1:4193";
await mkdir(".local/browser-checks", { recursive: true });
const profile = await mkdtemp(`${process.cwd()}/.local/browser-checks/profile-`);
const { workshops } = JSON.parse(await readFile("course/catalog.json", "utf8"));
const lessons = await Promise.all(workshops.flatMap((group) => group.lessons).map(async (id) =>
  JSON.parse(await readFile(`course/workshops/${id}/lesson.json`, "utf8"))));
const topics = await Promise.all((await readdir("src/data/modules")).map(async (name) =>
  JSON.parse(await readFile(`src/data/modules/${name}`, "utf8"))));
const browser = spawn(process.env.CLI_SILLY_TEST_BROWSER ?? "brave-browser", [
  "--headless", "--disable-gpu", "--no-first-run", "--disable-background-networking",
  "--remote-debugging-port=0", `--user-data-dir=${profile}`, "about:blank",
], { stdio: "ignore" });
let launchError;
browser.on("error", (error) => { launchError = error; });
let socket;
try {
  let port;
  for (let i = 0; i < 100 && !port; i++) {
    if (launchError) throw launchError;
    port = (await readFile(`${profile}/DevToolsActivePort`, "utf8").catch(() => "")).split("\n")[0];
    if (!port) await delay(100);
  }
  assert(port, "Browser did not start");
  const targets = await fetch(`http://127.0.0.1:${port}/json/list`).then((response) => response.json());
  socket = new WebSocket(targets.find((target) => target.type === "page").webSocketDebuggerUrl);
  await new Promise((resolve, reject) => {
    socket.addEventListener("open", resolve);
    socket.addEventListener("error", reject);
  });
  const pending = new Map();
  const errors = [];
  let id = 0;
  socket.onmessage = ({ data }) => {
    const message = JSON.parse(data);
    if (message.id) pending.get(message.id)?.(message);
    if (message.method === "Runtime.exceptionThrown") errors.push(message.params.exceptionDetails.text);
  };
  async function send(method, params = {}) {
    const current = ++id;
    const response = new Promise((resolve) => pending.set(current, resolve));
    socket.send(JSON.stringify({ id: current, method, params }));
    const result = await response;
    pending.delete(current);
    assert(!result.error, JSON.stringify(result.error));
    return result.result;
  }
  async function evaluate(expression) {
    return (await send("Runtime.evaluate", { expression, returnByValue: true, awaitPromise: true, userGesture: true })).result.value;
  }
  async function waitFor(expression) {
    for (let i = 0; i < 200; i++) {
      if (await evaluate(expression)) return;
      await delay(50);
    }
    assert.fail(`Timed out: ${expression}`);
  }
  async function open(path, expression) {
    await send("Page.navigate", { url: base + path });
    await waitFor(`location.origin === ${JSON.stringify(base)} && location.pathname === ${JSON.stringify(path)} && (${expression})`);
  }
  await send("Runtime.enable");
  await open("/", "document.querySelector('main h1')?.textContent.includes('in your terminal')");
  const followUps = lessons.filter((lesson) => lesson.reviewOf);
  assert(await evaluate(`document.querySelector('main').textContent.includes('${followUps.length} follow-up challenges')`));
  assert.equal(await evaluate("document.querySelectorAll('main article').length"), 5);
  for (const lesson of lessons) {
    await open(`/workshop/${lesson.id}`, `document.querySelector('main h1')?.textContent === ${JSON.stringify(lesson.title)} && document.querySelector('#commands')?.textContent === 'Commands and operations' && document.querySelectorAll('#example li').length === ${lesson.demo.length}`);
    assert(await evaluate(`[...document.querySelectorAll('main code')].some((code) => code.textContent === ${JSON.stringify(`cli-silly do ${lesson.id}`)})`));
    if (lesson.reviewOf) assert(await evaluate(`!!document.querySelector('#practice a[href="/workshop/${lesson.reviewOf}"]')`));
    for (const child of followUps.filter((candidate) => candidate.reviewOf === lesson.id)) {
      assert(await evaluate(`!!document.querySelector('#practice a[href="/workshop/${child.id}"]')`));
    }
  }
  for (const topic of topics) {
    await open(`/${topic.id}`, `!!document.querySelector('main h1') && [...document.querySelectorAll('main span')].some((span) => span.textContent === ${JSON.stringify(`#${topic.id}`)})`);
    const diagrams = (topic.content?.match(/```mermaid/g) ?? []).length;
    if (diagrams) await waitFor(`document.querySelectorAll('.mermaid svg').length === ${diagrams}`);
  }
  await send("Page.navigate", { url: base + "/intro" });
  await waitFor("location.pathname === '/1-intro' && !!document.querySelector('main h1')");
  await open("/workshop/foundations/cli-intro", "!!document.querySelector('#terminal-drills')");
  await evaluate("document.querySelector('#terminal-drills summary').click()");
  await waitFor("document.querySelector('#terminal-drills').open && document.querySelector('#terminal-drills').innerText.includes('Try it in Bash')");
  await evaluate("[...document.querySelectorAll('main button')].find((button) => button.textContent === 'Copy command').click()");
  await waitFor("document.querySelector('main [role=status]').textContent.length > 0");
  await evaluate("document.querySelector('input[type=search]').focus()");
  await send("Input.insertText", { text: "ssh" });
  await waitFor("document.querySelectorAll('aside nav section a').length === 1");
  // A full navigation resets the local search before checking the mobile drawer.
  await open("/install", "document.querySelector('main h1')?.textContent === 'Install cli-silly on Linux'");
  assert(await evaluate("document.querySelector('main').textContent.includes('Ctrl+k')"));
  await send("Emulation.setDeviceMetricsOverride", { width: 375, height: 812, deviceScaleFactor: 1, mobile: true });
  await open("/workshop/core-tools/node-cli", "document.querySelector('main h1')?.textContent === 'Node.js CLI Development' && document.documentElement.scrollWidth <= 375");
  assert.equal(await evaluate("getComputedStyle(document.querySelector('aside')).visibility"), "hidden");
  await evaluate("document.querySelector('button[aria-label=\"Toggle navigation\"]').click()");
  await waitFor("document.querySelector('button[aria-label=\"Toggle navigation\"]').getAttribute('aria-expanded') === 'true'");
  assert.deepEqual(errors, []);
  console.log(`Verified ${lessons.length} workshops, ${topics.length} legacy readers, project links, search, drills, copy feedback, install, and mobile navigation at ${base}.`);
} finally {
  socket?.close();
  browser.kill();
}
