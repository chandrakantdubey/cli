import catalog from "../../course/catalog.json";

interface Lesson {
  id: string;
  title: string;
  task: string;
  tools?: string[];
  demo: { command: string; visual: string; explanation: string }[];
  checks: { id: string; label: string; points: number }[];
}

const manifests = import.meta.glob<Lesson>("../../course/workshops/*/*/lesson.json", {
  eager: true,
  import: "default",
});

export const workshops = catalog.workshops.map((workshop) => ({
  ...workshop,
  lessons: workshop.lessons.map((id) => {
    const lesson = manifests[`../../course/workshops/${id}/lesson.json`];
    if (!lesson || lesson.id !== id) throw new Error(`Missing lesson manifest for ${id}`);
    return lesson;
  }),
}));
