export default function Brand() {
  return (
    <span className="inline-flex items-center gap-2.5">
      <span className="brand-mark" aria-hidden="true">
        <span className="brand-chevron" />
        <span className="brand-cursor" />
      </span>
      <span className="font-mono text-lg font-bold tracking-tight text-slate-950">cli-silly</span>
    </span>
  );
}
