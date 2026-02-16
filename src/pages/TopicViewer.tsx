import { useParams, Navigate } from "react-router";
import ReactMarkdown from "react-markdown";
import rehypeHighlight from "rehype-highlight";
import remarkGfm from "remark-gfm";
import "highlight.js/styles/github-dark.css";
import { courseContent } from "../data/moduleRegistry";
import { useEffect, useRef, useState } from "react";
import mermaid from "mermaid";

// Initialize mermaid
mermaid.initialize({
  startOnLoad: false,
  theme: "base",
  themeVariables: {
    fontFamily: "Inter, sans-serif",
    primaryColor: "#e0e7ff", // indigo-100
    edgeLabelBackground: "#ffffff",
    tertiaryColor: "#f1f5f9", // slate-100
  },
  securityLevel: "loose",
});

const Mermaid = ({ chart }: { chart: string }) => {
  const ref = useRef<HTMLDivElement>(null);
  const [svg, setSvg] = useState("");

  useEffect(() => {
    if (chart && ref.current) {
      const id = `mermaid-${Math.random().toString(36).substr(2, 9)}`;
      mermaid
        .render(id, chart)
        .then(({ svg }) => {
          setSvg(svg);
        })
        .catch((error) => {
          console.error("Mermaid render error:", error);
          setSvg(
            `<div class="text-red-500 bg-red-50 p-2 rounded text-sm">Failed to render diagram</div>`,
          );
        });
    }
  }, [chart]);

  return (
    <div
      className="mermaid my-6 flex justify-center bg-white p-4 rounded-lg shadow-sm border border-slate-100 overflow-x-auto"
      ref={ref}
      dangerouslySetInnerHTML={{ __html: svg }}
    />
  );
};

export default function TopicViewer() {
  const { topicId } = useParams();
  const topic = courseContent.find((t) => t.id === topicId);

  if (!topic) {
    return <Navigate to="/intro" replace />;
  }

  return (
    <article className="prose prose-slate prose-lg w-full max-w-none prose-p:text-slate-700 prose-headings:text-slate-900 prose-li:text-slate-700 prose-strong:text-slate-900">
      <div className="bg-white py-4 md:py-6 px-8 rounded-lg">
        <div className="flex items-center space-x-3 mb-8">
          <span
            className={`px-2.5 py-1 text-xs font-bold rounded-md uppercase tracking-wide ${
              topic.level === "Basic"
                ? "bg-teal-50 text-teal-700 border border-teal-100"
                : "bg-indigo-50 text-indigo-700 border border-indigo-100"
            }`}
          >
            {topic.level}
          </span>
          <span className="text-slate-400 text-sm font-medium">
            #{topic.id}
          </span>
        </div>

        <div className="markdown-body">
          <ReactMarkdown
            rehypePlugins={[rehypeHighlight]}
            remarkPlugins={[remarkGfm]}
            components={{
              code: ({ node, inline, className, children, ...props }: any) => {
                const match = /language-(\w+)/.exec(className || "");
                const isMermaid = match && match[1] === "mermaid";

                if (isMermaid) {
                  return (
                    <Mermaid chart={String(children).replace(/\n$/, "")} />
                  );
                }

                return inline ? (
                  <code
                    className="bg-slate-100 text-slate-800 px-1.5 py-0.5 rounded text-sm font-mono border border-slate-200"
                    {...props}
                  >
                    {children}
                  </code>
                ) : (
                  <code className={className} {...props}>
                    {children}
                  </code>
                );
              },
              pre: ({ node, ...props }) => (
                <pre
                  className="rounded-xl shadow-sm border border-slate-200 bg-[#0d1117] overflow-x-hidden p-0 my-8 px-4"
                  {...props}
                />
              ),
              h1: ({ node, ...props }) => (
                <h1
                  className="text-4xl font-extrabold tracking-tight text-slate-900 mb-8 pb-4 border-b border-slate-100"
                  {...props}
                />
              ),
              h2: ({ node, ...props }) => (
                <h2
                  className="text-2xl font-bold tracking-tight text-slate-900 mt-12 mb-6"
                  {...props}
                />
              ),
              h3: ({ node, ...props }) => (
                <h3
                  className="text-xl font-bold tracking-tight text-slate-800 mt-8 mb-4"
                  {...props}
                />
              ),
              table: ({ node, ...props }) => (
                <div className="overflow-x-auto my-8 border border-slate-200 rounded-lg shadow-sm">
                  <table
                    className="min-w-full divide-y divide-slate-200"
                    {...props}
                  />
                </div>
              ),
              thead: ({ node, ...props }) => (
                <thead className="bg-slate-50" {...props} />
              ),
              th: ({ node, ...props }) => (
                <th
                  className="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider"
                  {...props}
                />
              ),
              tbody: ({ node, ...props }) => (
                <tbody
                  className="bg-white divide-y divide-slate-100"
                  {...props}
                />
              ),
              td: ({ node, ...props }) => (
                <td
                  className="px-6 py-4 whitespace-nowrap text-sm text-slate-600 font-medium"
                  {...props}
                />
              ),
              blockquote: ({ node, ...props }) => (
                <blockquote
                  className="border-l-4 border-blue-500 bg-blue-50/50 pl-4 py-1 pr-1 my-6 rounded-r italic text-slate-700"
                  {...props}
                />
              ),
            }}
          >
            {topic.content}
          </ReactMarkdown>
        </div>
      </div>
    </article>
  );
}
