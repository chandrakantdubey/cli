import { createBrowserRouter, redirect, RouterProvider } from "react-router";
import { loadReference, loadWorkshopLesson } from "./course/content";
import CourseLayout from "./layouts/CourseLayout";
import Loading from "./pages/Loading";
import Home from "./pages/Home";

const router = createBrowserRouter([
  {
    path: "/",
    Component: CourseLayout,
    HydrateFallback: Loading,
    children: [
      { index: true, Component: Home },
      { path: "install", lazy: async () => ({ Component: (await import("./pages/Install")).default }) },
      {
        path: "workshop/:workshopId/:lessonId",
        loader: ({ params }) => loadWorkshopLesson(params.workshopId, params.lessonId),
        lazy: async () => ({ Component: (await import("./pages/WorkshopLesson")).default }),
      },
      { path: "intro", loader: () => redirect("/1-intro") },
      {
        path: ":topicId",
        loader: async ({ params }) => (await loadReference(params.topicId)) ?? redirect("/intro"),
        lazy: async () => ({ Component: (await import("./pages/TopicViewer")).default }),
      },
    ],
  },
]);

export default function App() {
  return <RouterProvider router={router} />;
}
