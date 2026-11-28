import { Link } from "react-router";
import projects from "../../course/projects.json";
import { workshops } from "../course/catalog";

const lesson = workshops[0].lessons[0];
const allLessons = workshops.flatMap((workshop) => workshop.lessons);
const lessonTitles = new Map(allLessons.map((item) => [item.id, item.title]));
const coreLabs = allLessons.filter((item) => !item.reviewOf).length;
const challenges = allLessons.filter((item) => item.reviewOf).length;

export default function Home() {
  return (
    <section className="mx-auto max-w-4xl py-8 md:py-14">
      <h1 className="max-w-3xl text-4xl font-semibold tracking-tight text-slate-950 md:text-5xl">
        Learn developer tools in your terminal.
      </h1>
      <p className="mt-4 font-mono text-sm font-medium text-violet-800">Native Linux · amd64 and arm64 · Offline lessons · Local progress</p>
      <p className="mt-5 max-w-2xl text-lg leading-8 text-slate-700">
        Practice Linux, Git, Bash, and Docker through {coreLabs} scored local labs and {challenges} follow-up challenges.
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
          <pre className="mt-3 overflow-x-auto rounded bg-slate-100 p-3 text-sm"><code>cli-silly do {lesson.id}</code></pre>
        </section>
      </div>
      <section className="mt-10">
        <h2 className="text-2xl font-semibold">Connected developer projects</h2>
        <p className="mt-3 text-slate-700">Apply the lessons to a complete workflow. Each step has its own workspace, feedback, and saved score.</p>
        <div className="mt-5 grid gap-4 sm:grid-cols-2">
          {projects.map((project) => (
            <article key={project.id} className="rounded-xl border border-slate-200 bg-white p-5">
              <h3 className="font-semibold">{project.title}</h3>
              <p className="mt-2 text-sm leading-6 text-slate-700">{project.outcome}</p>
              <ol className="mt-3 list-inside list-decimal space-y-2 text-sm">
                {project.lessons.map((id) => <li key={id}><Link className="text-violet-800 underline" to={`/workshop/${id}`}>{lessonTitles.get(id) ?? id}</Link></li>)}
              </ol>
            </article>
          ))}
        </div>
      </section>
      <section className="mt-10 rounded-xl border border-slate-200 bg-white p-6">
        <h2 className="text-2xl font-semibold">Built for developers learning to work in Linux</h2>
        <p className="mt-3 leading-7 text-slate-700">Start with files, paths, permissions, and streams. Apply those skills to Git, Bash automation, local APIs, and Docker. Every core lesson explains its commands, common mistakes, and a practical development use.</p>
        <ol className="mt-5 list-inside list-decimal space-y-3 text-slate-700">
          <li><strong>Read.</strong> Follow the explanation and worked example at your own pace.</li>
          <li><strong>Do.</strong> Inspect the starter files and complete an observable task in Bash.</li>
          <li><strong>Check and retry.</strong> Get feedback on each result, then try related practice with different inputs.</li>
        </ol>
        <p className="mt-5 text-sm leading-6 text-slate-600">Bash and standard Linux utilities cover the first lessons. Git, jq, Python, Node.js, and Docker are needed only by their relevant workshops. Run <code>cli-silly doctor</code> to see what is ready. Docker practice also requires its image to be downloaded before going offline.</p>
        <p className="mt-3 text-sm leading-6 text-slate-600">Your practice runs with your Linux user permissions. Archived attempts and local backups help you recover your work. No account is needed.</p>
        <a className="mt-4 inline-block font-medium text-violet-800 underline" href="https://github.com/chandrakantdubey/cli/releases">Release notes and Linux downloads</a>
      </section>
    </section>
  );
}
