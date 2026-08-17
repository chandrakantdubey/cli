import { Link } from "react-router";
import { courseContent } from "../data/moduleRegistry";

const basics = courseContent.filter((topic) => topic.level === "Basic");
const advanced = courseContent.filter((topic) => topic.level === "Advanced");

export default function Home() {
  const firstTopic = courseContent[0];

  return (
    <section className="mx-auto max-w-5xl px-6 py-10 md:px-10 md:py-14">
      <div className="max-w-3xl">
        <p className="text-xs font-semibold uppercase tracking-[0.18em] text-slate-400">
          CLI Mastery
        </p>
        <h1 className="mt-3 text-4xl font-semibold tracking-tight text-slate-950 md:text-5xl">
          Learn the command line by understanding the system behind it.
        </h1>
        <p className="mt-5 max-w-2xl text-lg leading-8 text-slate-600">
          A practical path through Linux, shell usage, automation, networking,
          Git, containers, troubleshooting, and production operations.
        </p>

        <div className="mt-8 flex flex-wrap gap-3">
          <Link
            to={`/${firstTopic.id}`}
            className="inline-flex items-center rounded-lg bg-slate-900 px-4 py-2.5 text-sm font-semibold text-white transition hover:bg-slate-800"
          >
            Start with the intro
          </Link>
          <Link
            to={`/${basics[1]?.id ?? firstTopic.id}`}
            className="inline-flex items-center rounded-lg border border-slate-200 bg-white px-4 py-2.5 text-sm font-semibold text-slate-700 transition hover:border-slate-300 hover:bg-slate-50"
          >
            Explore the basics
          </Link>
        </div>
      </div>

      <div className="mt-14 grid gap-5 md:grid-cols-2">
        <div className="rounded-2xl border border-slate-200 bg-white p-6">
          <p className="text-sm font-semibold text-slate-900">Foundations</p>
          <p className="mt-2 text-sm leading-6 text-slate-500">
            Build the mental model first: shell, filesystem, users, streams,
            packages, networking, processes, scripting, and Git.
          </p>
          <p className="mt-4 text-xs font-medium uppercase tracking-wider text-slate-400">
            {basics.length} topics
          </p>
        </div>

        <div className="rounded-2xl border border-slate-200 bg-white p-6">
          <p className="text-sm font-semibold text-slate-900">Advanced operations</p>
          <p className="mt-2 text-sm leading-6 text-slate-500">
            Go deeper into services, storage, containers, observability,
            security, automation, cloud-native systems, and operations.
          </p>
          <p className="mt-4 text-xs font-medium uppercase tracking-wider text-slate-400">
            {advanced.length} topics
          </p>
        </div>
      </div>
    </section>
  );
}
