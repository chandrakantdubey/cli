import { topics, workshops } from "./catalog";

interface Lesson {
  id: string;
  title: string;
  task: string;
  reviewOf?: string;
  tools?: string[];
  demo: { command: string; visual: string; explanation: string }[];
  checks: { id: string; label: string; points: number }[];
}

interface Topic {
  id: string;
  title: string;
  level: "Basic" | "Advanced";
  content?: string;
}

const manifests = import.meta.glob<Lesson>("../../course/workshops/*/*/lesson.json", { import: "default" });
const bodies = import.meta.glob<string>("../../course/workshops/*/*/lesson.md", { import: "default", query: "?raw" });
const references = import.meta.glob<Topic>("../data/modules/*.json", { import: "default" });

export async function loadWorkshopLesson(workshopId?: string, slug?: string) {
  const id = `${workshopId}/${slug}`;
  const workshop = workshops.find((group) => group.lessons.some((lesson) => lesson.id === id));
  if (!workshop) return null;
  const base = `../../course/workshops/${id}`;
  const [lesson, body] = await Promise.all([
    manifests[`${base}/lesson.json`](),
    bodies[`${base}/lesson.md`](),
  ]);
  return { lesson, body, workshopTitle: workshop.title };
}

export async function loadReference(id?: string) {
  const topic = topics.find((entry) => entry.id === id);
  if (!topic) return null;
  const [reference, body] = await Promise.all([
    references[`../data/modules/${topic.file}`](),
    topic.lessonId ? bodies[`../../course/workshops/${topic.lessonId}/lesson.md`]() : undefined,
  ]);
  return {
    ...reference,
    content: body === undefined ? reference.content : `# ${reference.title}\n\n${body}`,
    lessonId: topic.lessonId,
  };
}
