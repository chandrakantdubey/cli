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
            className={`px-2.5 py-1 text-xs font-bold rounded-md uppercase tracking-wide ${topic.level === "Basic"
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

        <div className="markdown-body [&_:not(pre)>code]:rounded [&_:not(pre)>code]:text-slate-800">
          <ReactMarkdown
            rehypePlugins={[rehypeHighlight]}
            remarkPlugins={[remarkGfm]}
            components={{
              code: ({ className, children }) => {
                if (className?.split(" ").includes("language-mermaid")) {
                  return <Mermaid chart={String(children).replace(/\n$/, "")} />;
                }
                return <code className={className}>{children}</code>;
              },
              pre: ({ children }) => (
                <pre
                  className="rounded-xl shadow-sm border border-slate-200 bg-[#0d1117] overflow-x-auto p-0 my-8 px-4">
                  {children}
                </pre>
              ),
              h1: ({ children }) => (
                <h1
                  className="text-4xl font-extrabold tracking-tight text-slate-900 mb-8 pb-4 border-b border-slate-100">
                  {children}
                </h1>
              ),
              h2: ({ children }) => (
                <h2
                  className="text-2xl font-bold tracking-tight text-slate-900 mt-12 mb-6">
                  {children}
                </h2>
              ),
              h3: ({ children }) => (
                <h3
                  className="text-xl font-bold tracking-tight text-slate-800 mt-8 mb-4">
                  {children}
                </h3>
              ),
              table: ({ children }) => (
                <div className="overflow-x-auto my-8 border border-slate-200 rounded-lg shadow-sm">
                  <table
                    className="min-w-full divide-y divide-slate-200">
                    {children}
                  </table>
                </div>
              ),
              thead: ({ children }) => (
                <thead className="bg-slate-50">{children}</thead>
              ),
              th: ({ children, style }) => (
                <th
                  className="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider" style={style}>
                  {children}
                </th>
              ),
              tbody: ({ children }) => (
                <tbody
                  className="bg-white divide-y divide-slate-100">
                  {children}
                </tbody>
              ),
              td: ({ children, style }) => (
                <td
                  className="px-6 py-4 whitespace-nowrap text-sm text-slate-600 font-medium" style={style}>
                  {children}
                </td>
              ),
              blockquote: ({ children }) => (
                <blockquote
                  className="border-l-4 border-blue-500 bg-blue-50/50 pl-4 py-1 pr-1 my-6 rounded-r italic text-slate-700">
                  {children}
                </blockquote>
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
