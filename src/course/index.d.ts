declare module "virtual:cli-silly-index" {
  interface LessonSummary {
    id: string;
    title: string;
    reviewOf?: string;
  }

  export const workshops: { id: string; title: string; lessons: LessonSummary[] }[];
  export const topics: { id: string; title: string; file: string }[];
}
