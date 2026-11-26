import { lazy, Suspense } from "react";
import { BrowserRouter, Routes, Route, Navigate } from "react-router";
import Loading from "./pages/Loading";
import Home from "./pages/Home";

const Install = lazy(() => import("./pages/Install"));
const WorkshopLesson = lazy(() => import("./pages/WorkshopLesson"));
const CourseLayout = lazy(() => import("./layouts/CourseLayout"));
const TopicViewer = lazy(() => import("./pages/TopicViewer"));

export default function App() {
  return (
    <BrowserRouter>
      <Suspense fallback={<Loading />}>
        <Routes>
          <Route path="/" element={<CourseLayout />}>
            <Route index element={<Home />} />
            <Route path="install" element={<Install />} />
            <Route path="workshop/:workshopId/:lessonId" element={<WorkshopLesson />} />
            <Route path="intro" element={<Navigate to="/1-intro" replace />} />
            <Route path=":topicId" element={<TopicViewer />} />
          </Route>
        </Routes>
      </Suspense>
    </BrowserRouter>
  );
}
