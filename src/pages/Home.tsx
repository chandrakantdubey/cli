import { Link } from "react-router";
import { workshops } from "../course/catalog";

const lesson = workshops[0].lessons[0];
const coreLabs = workshops.flatMap((workshop) => workshop.lessons).filter((item) => !item.reviewOf).length;

export default function Home() {
  return (
    <section className="mx-auto max-w-4xl py-8 md:py-14">
      <h1 className="max-w-3xl text-4xl font-semibold tracking-tight text-slate-950 md:text-5xl">
        Learn developer tools in your terminal.
      </h1>
      <p className="mt-5 max-w-2xl text-lg leading-8 text-slate-700">
        Practice Linux, Git, Bash, and Docker through {coreLabs} scored local labs and a review challenge.
        Read a lesson, practice in a real Bash workspace, and get scored feedback. Your files and scores stay on your machine.
      </p>
      <Link
        to={`/workshop/${lesson.id}`}
        className="mt-7 inline-flex rounded-lg bg-violet-700 px-4 py-3 text-sm font-semibold text-white hover:bg-violet-800"
      >
        Read the first lesson
      </Link>
      <Link to="/install" className="ml-4 mt-7 inline-flex rounded-lg border border-violet-200 px-4 py-3 text-sm font-semibold text-violet-800 hover:bg-violet-50">
        Install for Linux
      </Link>

      <div className="mt-10 grid gap-6 md:grid-cols-2">
        <section className="min-w-0 rounded-xl bg-slate-950 p-5 text-white">
          <h2 className="text-sm font-semibold text-violet-300">Sample lab result</h2>
          <pre className="mt-4 overflow-x-auto text-xs leading-7"><code><span className="text-cyan-200">$ cli-silly check</span>{"\n"}<span className="text-emerald-300">PASS</span>{"  40/40  Script is executable\n"}<span className="text-amber-200">TRY</span>{"    0/60  Create the status report\n\n"}{"40/100  ████░░░░░░"}</code></pre>
        </section>
        <section className="min-w-0 rounded-xl border border-slate-200 bg-white p-5">
          <h2 className="font-semibold text-slate-950">Practice in the CLI</h2>
          <p className="mt-2 text-sm leading-6 text-slate-700">Install cli-silly on Linux. Then open this lesson in your terminal.</p>
          <pre className="mt-3 overflow-x-auto rounded bg-slate-100 p-3 text-sm"><code>cli-silly learn {lesson.id}</code></pre>
        </section>
      </div>
    </section>
  );
}
