import { useState } from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { Link, useParams } from "react-router";
import { workshops } from "../course/catalog";

const bodies = import.meta.glob<string>("../../course/workshops/*/*/lesson.md", {
  eager: true,
  import: "default",
  query: "?raw",
});

export default function WorkshopLesson() {
  const { workshopId, lessonId } = useParams();
  const workshop = workshops.find((item) => item.id === workshopId);
  const lesson = workshop?.lessons.find((item) => item.id === `${workshopId}/${lessonId}`);
  if (!workshop || !lesson) {
    return <p>Workshop lesson unavailable. <Link to="/">Browse the course</Link>.</p>;
  }
  return <LessonContent key={lesson.id} lesson={lesson} workshopTitle={workshop.title} />;
}

function LessonContent({ lesson, workshopTitle }: {
  lesson: (typeof workshops)[number]["lessons"][number];
  workshopTitle: string;
}) {
  const body = bodies[`../../course/workshops/${lesson.id}/lesson.md`];
  const [frame, setFrame] = useState(0);
  const step = lesson.demo[frame];

  return (
    <article className="mx-auto max-w-3xl">
      <p className="font-mono text-sm font-semibold text-violet-700">{workshopTitle}</p>
      <h1 className="mt-3 text-4xl font-semibold tracking-tight text-slate-950">{lesson.title}</h1>

      <div className="prose prose-slate mt-6 max-w-none">
        <ReactMarkdown remarkPlugins={[remarkGfm]}>{body}</ReactMarkdown>
      </div>

      <section className="mt-8 rounded-xl bg-slate-950 p-5 text-slate-100 sm:p-7" aria-label="Worked example">
        <div aria-live="polite">
          <p className="text-sm font-medium text-violet-300">Worked example · {frame + 1} of {lesson.demo.length}</p>
          <pre className="mt-5 overflow-x-auto text-sm text-amber-200"><code>$ {step.command}</code></pre>
          <pre className="my-4 overflow-x-auto text-sm text-cyan-200">{step.visual}</pre>
          <p className="text-sm leading-relaxed">{step.explanation}</p>
        </div>
        <div className="mt-5 flex flex-wrap gap-2">
          <button className="rounded border border-slate-500 px-3 py-2 text-sm disabled:opacity-40"
            onClick={() => setFrame(frame - 1)} disabled={frame === 0}>Previous</button>
          <button className="rounded bg-violet-300 px-3 py-2 text-sm font-medium text-slate-950 disabled:opacity-40"
            onClick={() => setFrame(frame + 1)} disabled={frame === lesson.demo.length - 1}>Next step</button>
        </div>
      </section>

      <section className="mt-8 rounded-xl border border-slate-200 bg-white p-5 sm:p-7">
        <h2 className="text-xl font-semibold">Your task</h2>
        <p className="mt-3 text-slate-700">{lesson.task}</p>
        {lesson.tools?.length ? <p className="mt-2 text-sm">Tools needed: {lesson.tools.join(", ")}.</p> : null}
        <ul className="my-4 list-inside list-disc text-sm leading-7">
          {lesson.checks.map((check) => <li key={check.id}>{check.label} · {check.points} points</li>)}
        </ul>
        <p className="text-sm">
          <Link className="font-medium text-violet-800 underline" to="/install">Install cli-silly</Link> and open this lesson in your terminal.
        </p>
        <pre className="mt-3 overflow-x-auto rounded bg-slate-100 p-3 text-sm"><code>cli-silly learn {lesson.id}</code></pre>
      </section>
    </article>
  );
}
