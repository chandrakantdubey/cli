import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { Link, useLoaderData } from "react-router";
import type { loadWorkshopLesson } from "../course/content";

export default function WorkshopLesson() {
  const content = useLoaderData<typeof loadWorkshopLesson>();
  if (!content) {
    return <p>Workshop lesson unavailable. <Link to="/">Browse the course</Link>.</p>;
  }
  const { lesson, body, workshopTitle } = content;

  return (
    <article className="mx-auto max-w-3xl min-w-0">
      <p className="font-mono text-sm font-semibold text-violet-700">{workshopTitle}</p>
      <h1 className="mt-3 text-4xl font-semibold tracking-tight text-slate-950">{lesson.title}</h1>
      <nav aria-label="On this lesson" className="mt-5 flex flex-wrap gap-x-5 gap-y-2 text-sm font-medium text-violet-800">
        <a className="underline decoration-violet-300 underline-offset-4" href="#commands">Commands and operations</a>
        <a className="underline decoration-violet-300 underline-offset-4" href="#example">Worked example</a>
        <a className="underline decoration-violet-300 underline-offset-4" href="#practice">Practice task</a>
      </nav>

      <div className="prose prose-slate mt-7 max-w-none prose-pre:overflow-x-auto prose-h2:scroll-mt-20">
        <ReactMarkdown
          remarkPlugins={[remarkGfm]}
          components={{
            h2: ({ children }) => (
              <h2 id={children === "Commands and operations" ? "commands" : undefined}>{children}</h2>
            ),
          }}
        >
          {body}
        </ReactMarkdown>
      </div>

      <section id="example" className="mt-10 scroll-mt-20" aria-labelledby="example-heading">
        <h2 id="example-heading" className="text-2xl font-semibold">Worked example</h2>
        <ol className="mt-4 space-y-3">
          {lesson.demo.map((step, index) => (
            <li key={`${index}-${step.command}`} className="min-w-0 rounded-lg bg-slate-950 p-4 text-slate-100 sm:p-5">
              <pre className="overflow-x-auto text-sm text-amber-200"><code><span aria-hidden="true">$ </span>{step.command}</code></pre>
              {step.visual ? <pre className="mt-3 overflow-x-auto text-sm text-cyan-200">{step.visual}</pre> : null}
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
        <pre className="mt-3 overflow-x-auto rounded bg-slate-100 p-3 text-sm"><code>cli-silly do {lesson.id}</code></pre>
      </section>
    </article>
  );
}
