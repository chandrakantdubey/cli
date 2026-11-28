import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { Link, useLoaderData } from "react-router";
import type { loadWorkshopLesson } from "../course/content";
import CopyCommand from "../components/CopyCommand";

function sectionId(text: string) {
  return text === "Commands and operations" ? "commands" : text.toLowerCase().replace(/[^\p{L}\p{N}]+/gu, "-").replace(/^-|-$/g, "");
}

export default function WorkshopLesson() {
  const content = useLoaderData<typeof loadWorkshopLesson>();
  if (!content) {
    return <p>Workshop lesson unavailable. <Link to="/">Browse the course</Link>.</p>;
  }
  const { lesson, body, terminalNotes, workshopTitle, followUps } = content;

  return (
    <article className="mx-auto max-w-3xl min-w-0">
      <p className="font-mono text-sm font-semibold text-violet-700">{workshopTitle}</p>
      <h1 className="mt-3 text-4xl font-semibold tracking-tight text-slate-950">{lesson.title}</h1>
      <nav aria-label="On this lesson" className="mt-5 flex flex-wrap gap-x-5 gap-y-2 text-sm font-medium text-violet-800">
        <a className="underline decoration-violet-300 underline-offset-4" href="#commands">Commands and operations</a>
        <a className="underline decoration-violet-300 underline-offset-4" href="#example">Worked example</a>
        <a className="underline decoration-violet-300 underline-offset-4" href="#practice">Practice task</a>
        {terminalNotes ? <a className="underline decoration-violet-300 underline-offset-4" href="#terminal-drills">Terminal drills</a> : null}
      </nav>

      <div className="prose prose-slate mt-7 max-w-none prose-pre:overflow-x-auto prose-h2:scroll-mt-20">
        <ReactMarkdown
          remarkPlugins={[remarkGfm]}
          components={{
            h2: ({ children }) => (
              <h2 id={typeof children === "string" ? sectionId(children) : undefined}>{children}</h2>
            ),
          }}
        >
          {body}
        </ReactMarkdown>
      </div>

      {terminalNotes ? (
        <details id="terminal-drills" className="mt-8 scroll-mt-20 rounded-xl border border-violet-200 bg-violet-50 p-5">
          <summary className="cursor-pointer font-semibold text-violet-900">Terminal drills and extra practice</summary>
          <p className="mt-3 text-sm text-slate-700">These drills are also bundled in the offline CLI reader. Use the lesson workspace so the example files are available.</p>
          <div className="prose prose-slate mt-4 max-w-none prose-pre:overflow-x-auto"><ReactMarkdown remarkPlugins={[remarkGfm]}>{terminalNotes}</ReactMarkdown></div>
        </details>
      ) : null}

      <section id="example" className="mt-10 scroll-mt-20" aria-labelledby="example-heading">
        <h2 id="example-heading" className="text-2xl font-semibold">Worked example</h2>
        <ol className="mt-4 space-y-3">
          {lesson.demo.map((step, index) => (
            <li key={`${index}-${step.command}`} className="min-w-0 rounded-lg bg-slate-950 p-4 text-slate-100 sm:p-5">
              <pre className="overflow-x-auto text-sm text-amber-200"><code><span aria-hidden="true">$ </span>{step.command}</code></pre>
              {step.visual ? <pre className="mt-3 overflow-x-auto text-sm text-cyan-200">{step.visual}</pre> : null}
              {step.flow ? <p className="mt-3 text-sm font-medium text-violet-300">{step.flow.join(" → ")}</p> : null}
              <p className="mt-3 text-sm leading-relaxed">{step.explanation}</p>
            </li>
          ))}
        </ol>
      </section>

      <section id="practice" className="mt-8 scroll-mt-20 rounded-xl border border-slate-200 bg-white p-5 sm:p-7">
        <h2 className="text-xl font-semibold">Practice task</h2>
        {lesson.reviewOf ? (
          <p className="mt-3 text-sm">
            Optional review challenge. <Link className="font-medium text-violet-800 underline" to={`/workshop/${lesson.reviewOf}`}>Read the original lesson</Link>.
          </p>
        ) : null}
        <p className="mt-3 whitespace-pre-line text-slate-700">{lesson.task}</p>
        <p className="mt-2 text-sm">Lab tools: {["Bash", ...(lesson.tools ?? [])].join(", ")}.</p>
        <ul className="my-4 list-inside list-disc text-sm leading-7">
          {lesson.checks.map((check) => <li key={check.id}>{check.label} · {check.points} points</li>)}
        </ul>
        <p className="text-sm">
          <Link className="font-medium text-violet-800 underline" to="/install">Install cli-silly</Link> and open this lesson in your terminal.
          Run <code>cli-silly check</code> inside the lab to grade your work.
        </p>
        <CopyCommand key={lesson.id} command={`cli-silly do ${lesson.id}`} />
        <p className="mt-3 text-sm text-slate-600">Linux amd64 or arm64. Practice uses your normal user permissions. Scores and files remain local. Read the task, open Bash, and type <code>exit</code> to return.</p>
        {followUps.length ? (
          <div className="mt-5 border-t border-slate-200 pt-4">
            <h3 className="font-semibold">More practice</h3>
            <ul className="mt-2 space-y-2 text-sm">
              {followUps.map((item) => (
                <li key={item.id}><Link className="text-violet-800 underline" to={`/workshop/${item.id}`}>{item.title}</Link></li>
              ))}
            </ul>
          </div>
        ) : null}
      </section>
    </article>
  );
}
