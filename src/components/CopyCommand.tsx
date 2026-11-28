import { useState } from "react";

export default function CopyCommand({ command }: { command: string }) {
  const [message, setMessage] = useState("");
  async function copy() {
    try {
      await navigator.clipboard.writeText(command);
      setMessage("Copied. Paste into your terminal when ready.");
    } catch {
      setMessage("Select and copy the command below.");
    }
  }
  return (
    <div className="mt-3 rounded-lg border border-slate-200 bg-slate-100 p-3">
      <button type="button" onClick={copy} className="rounded px-2 py-1 text-sm font-medium text-violet-800 hover:bg-violet-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-600">Copy command</button>
      <pre className="mt-2 overflow-x-auto text-sm"><code>{command}</code></pre>
      <p className="mt-1 text-sm text-slate-700" role="status">{message}</p>
    </div>
  );
}
