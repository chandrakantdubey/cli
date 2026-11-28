import { useState } from "react";
import { Link, Outlet, useLocation, useNavigation } from "react-router";
import { topics, workshops } from "../course/catalog";
import Brand from "../components/Brand";

const referenceTopics = topics.filter((topic) => !topic.lessonId);

export default function CourseLayout() {
  const location = useLocation();
  const { state } = useNavigation();
  const currentTopicId = location.pathname.split("/").filter(Boolean).pop();
  const linkedLesson = topics.find((topic) => location.pathname === `/${topic.id}`)?.lessonId;
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const [query, setQuery] = useState("");
  const search = query.trim().toLowerCase();

  const closeSidebar = () => setIsSidebarOpen(false);

  return (
    <div className="min-h-screen bg-slate-50 text-slate-900">
      <a href="#lesson-content" className="sr-only focus:not-sr-only focus:fixed focus:left-4 focus:top-4 focus:z-50 focus:rounded focus:bg-white focus:p-3">Skip to lesson content</a>
      <div className="md:hidden fixed inset-x-0 top-0 z-40 flex h-14 items-center justify-between border-b border-slate-200 bg-white/95 px-4 backdrop-blur">
        <Link to="/" onClick={closeSidebar} className="text-sm font-semibold tracking-tight">
          <Brand />
        </Link>
        <button
          onClick={() => setIsSidebarOpen((open) => !open)}
          className="rounded-md p-2 text-slate-500 transition hover:bg-slate-100 hover:text-slate-900"
          aria-label="Toggle navigation"
          aria-expanded={isSidebarOpen}
        >
          <span className="text-lg leading-none">{isSidebarOpen ? "×" : "☰"}</span>
        </button>
      </div>

      {isSidebarOpen && (
        <button
          aria-label="Close navigation"
          className="fixed inset-0 z-30 bg-slate-950/15 md:hidden"
          onClick={closeSidebar}
        />
      )}

      <aside
        onKeyDown={(event) => { if (event.key === "Escape") { closeSidebar(); document.querySelector<HTMLButtonElement>('button[aria-label="Toggle navigation"]')?.focus(); } }}
        className={`fixed inset-y-0 left-0 z-40 w-72 border-r border-slate-200 bg-white transition-transform duration-200 md:translate-x-0 ${
          isSidebarOpen ? "translate-x-0" : "invisible -translate-x-full md:visible"
        }`}
      >
        <div className="flex h-full flex-col">
          <div className="border-b border-slate-100 px-5 py-5">
            <Link to="/" onClick={closeSidebar} className="block">
              <Brand />
              <div className="mt-1 text-xs text-slate-500">Terminal learning platform</div>
            </Link>
          </div>

          <nav className="flex-1 overflow-y-auto px-3 py-5">
            <label className="mb-4 block px-3 text-sm font-medium text-slate-700">Find a lesson or tool
              <input type="search" value={query} onChange={(event) => setQuery(event.target.value)} className="mt-2 w-full rounded-md border border-slate-300 px-2 py-2 font-normal focus:outline-2 focus:outline-violet-600" placeholder="Git, paths, Docker…" />
            </label>
            <div className="mb-2 px-3 text-[11px] font-semibold uppercase tracking-[0.16em] text-slate-400">
              Learn
            </div>

            <Link
              to="/"
              onClick={closeSidebar}
              className={`mb-4 block rounded-lg px-3 py-2 text-sm transition ${
                !currentTopicId
                  ? "bg-slate-100 font-medium text-slate-900"
                  : "text-slate-600 hover:bg-slate-50 hover:text-slate-900"
              }`}
            >
              Overview
            </Link>

            <Link to="/install" onClick={closeSidebar}
              aria-current={location.pathname === "/install" ? "page" : undefined}
              className="mb-5 block rounded-lg px-3 py-2 text-sm font-medium text-violet-800 hover:bg-violet-50">
              Install for Linux
            </Link>

            <div className="mb-6 space-y-4">
              {workshops.filter((workshop) => !search || workshop.title.toLowerCase().includes(search) || workshop.lessons.some((lesson) => `${lesson.title} ${lesson.id}`.toLowerCase().includes(search))).map((workshop) => (
                <section key={workshop.id}>
                  <h2 className="px-3 text-xs font-semibold text-violet-800">{workshop.title}</h2>
                  {workshop.lessons.filter((lesson) => !search || `${workshop.title} ${lesson.title} ${lesson.id}`.toLowerCase().includes(search)).map((lesson) => {
                    const href = `/workshop/${lesson.id}`;
                    const active = location.pathname === href || linkedLesson === lesson.id;
                    return (
                      <Link key={lesson.id} to={href} onClick={closeSidebar}
                        aria-current={active ? "page" : undefined}
                        className={`block rounded-lg px-3 py-2 text-sm leading-5 ${active ? "bg-violet-100 font-medium text-violet-950" : "text-slate-600 hover:bg-violet-50"}`}>
                        {lesson.title}
                      </Link>
                    );
                  })}
                </section>
              ))}
            </div>

            <section>
              <h2 className="px-3 text-xs font-semibold text-slate-600">Systems &amp; operations references</h2>
              <div className="mt-2 space-y-0.5">
                {referenceTopics.filter((topic) => !search || topic.title.toLowerCase().includes(search)).map((topic) => {
                  const active = currentTopicId === topic.id;
                  return (
                    <Link
                      key={topic.id}
                      to={`/${topic.id}`}
                      onClick={closeSidebar}
                      aria-current={active ? "page" : undefined}
                      className={`block rounded-lg px-3 py-2 text-sm leading-5 ${active ? "bg-slate-900 font-medium text-white" : "text-slate-600 hover:bg-slate-50 hover:text-slate-900"}`}
                    >
                      {topic.title.replace(/^\d+\.\s*/, "")}
                    </Link>
                  );
                })}
              </div>
            </section>
          </nav>
        </div>
      </aside>

      <main id="lesson-content" tabIndex={-1} className="min-h-screen md:ml-72" aria-busy={state !== "idle"}>
        {state !== "idle" ? <div className="fixed inset-x-0 top-0 z-50 h-0.5 bg-violet-500" role="status" aria-label="Loading lesson" /> : null}
        <div className="mx-auto min-h-screen max-w-6xl px-4 py-6 pt-20 md:px-10 md:py-10 md:pt-10">
          <Outlet />
        </div>
      </main>
    </div>
  );
}
