import ReactMarkdown from "react-markdown";
import guide from "../../docs/install.md?raw";

export default function Install() {
  return (
    <article className="prose prose-slate mx-auto max-w-3xl prose-pre:overflow-x-auto">
      <ReactMarkdown>{guide}</ReactMarkdown>
    </article>
  );
}
