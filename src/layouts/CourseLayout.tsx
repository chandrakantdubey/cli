import { useState } from "react";
import { Link, Outlet, useLocation } from "react-router";
import { courseContent } from "../data/moduleRegistry";

export default function CourseLayout() {
  const location = useLocation();
  const currentTopicId = location.pathname.split("/").filter(Boolean).pop();
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);

  const closeSidebar = () => setIsSidebarOpen(false);

  const groups = [
    {
      label: "Foundations",
      topics: courseContent.slice(0, 9),
    },
    {
      label: "Core tools",
      topics: courseContent.slice(9, 20),
    },
    {
      label: "Systems & operations",
      topics: courseContent.slice(20),
    },
  ];

  return (
    <div className="min-h-screen bg-slate-50 text-slate-900">
      <div className="md:hidden fixed inset-x-0 top-0 z-40 flex h-14 items-center justify-between border-b border-slate-200 bg-white/95 px-4 backdrop-blur">
        <Link to="/" onClick={closeSidebar} className="text-sm font-semibold tracking-tight">
          CLI Mastery
        </Link>
        <button
          onClick={() => setIsSidebarOpen((open) => !open)}
          className="rounded-md p-2 text-slate-500 transition hover:bg-slate-100 hover:text-slate-900"
          aria-label="Toggle navigation"
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
        className={`fixed inset-y-0 left-0 z-40 w-72 border-r border-slate-200 bg-white transition-transform duration-200 md:translate-x-0 ${
          isSidebarOpen ? "translate-x-0" : "-translate-x-full"
        }`}
      >
        <div className="flex h-full flex-col">
          <div className="border-b border-slate-100 px-5 py-5">
            <Link to="/" onClick={closeSidebar} className="block">
              <div className="text-base font-semibold tracking-tight text-slate-950">CLI Mastery</div>
              <div className="mt-1 text-xs text-slate-500">Linux, shell & operations</div>
            </Link>
          </div>

          <nav className="flex-1 overflow-y-auto px-3 py-5">
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

            <div className="space-y-5">
              {groups.map((group) => (
                <section key={group.label}>
                  <h2 className="px-3 text-[11px] font-semibold uppercase tracking-[0.14em] text-slate-400">
                    {group.label}
                  </h2>
                  <div className="mt-2 space-y-0.5">
                    {group.topics.map((topic) => {
                      const active = currentTopicId === topic.id;
                      return (
                        <Link
                          key={topic.id}
                          to={`/${topic.id}`}
                          onClick={closeSidebar}
                          className={`block rounded-lg px-3 py-2 text-sm leading-5 transition ${
                            active
                              ? "bg-slate-900 font-medium text-white"
                              : "text-slate-600 hover:bg-slate-50 hover:text-slate-900"
                          }`}
                        >
                          {topic.title.replace(/^\d+\.\s*/, "")}
                        </Link>
                      );
                    })}
                  </div>
                </section>
              ))}
            </div>
          </nav>
        </div>
      </aside>

      <main className="min-h-screen md:ml-72">
        <div className="mx-auto min-h-screen max-w-6xl px-4 py-6 pt-20 md:px-10 md:py-10 md:pt-10">
          <Outlet />
        </div>
      </main>
    </div>
  );
}
