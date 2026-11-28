import ReactMarkdown from "react-markdown";
import guide from "../../docs/install.md?raw";

export default function Install() {
  return (
    <article className="prose prose-slate mx-auto max-w-3xl prose-pre:overflow-x-auto">
      <ReactMarkdown components={{ a: ({ href, children }) => <a href={href === "recovery.md" ? "https://github.com/chandrakantdubey/cli/blob/main/docs/recovery.md" : href}>{children}</a> }}>{guide}</ReactMarkdown>
    </article>
  );
}
