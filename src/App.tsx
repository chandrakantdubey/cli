import { lazy, Suspense } from "react";
import { BrowserRouter, Routes, Route } from "react-router";
import Loading from "./pages/Loading";
import Home from "./pages/Home";

const CourseLayout = lazy(() => import("./layouts/CourseLayout"));
const TopicViewer = lazy(() => import("./pages/TopicViewer"));

export default function App() {
  return (
    <BrowserRouter>
      <Suspense fallback={<Loading />}>
        <Routes>
          <Route path="/" element={<CourseLayout />}>
            <Route index element={<Home />} />
            <Route path=":topicId" element={<TopicViewer />} />
          </Route>
        </Routes>
      </Suspense>
    </BrowserRouter>
  );
}
