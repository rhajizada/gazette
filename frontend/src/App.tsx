import { Routes, Route } from "react-router-dom";
import Feeds from "./pages/Feeds";
import Login from "./pages/Login";
import User from "./pages/User.tsx"
import NotFound from "./pages/NotFound";
import FeedDetail from "./pages/FeedDetail";
import ItemDetail from "./pages/ItemDetail";
import CallbackPage from "./pages/Callback";
import RequireAuth from "./components/RequireAuth";

export default function App() {
  return (
    <Routes>
      <Route
        path="/"
        element={
          <RequireAuth>
            <Feeds />
          </RequireAuth>
        }
      />
      <Route
        path="/feeds"
        element={
          <RequireAuth>
            <Feeds />
          </RequireAuth>
        }
      />
      <Route path="/feeds/:feedID" element={
        <RequireAuth>
          <FeedDetail />
        </RequireAuth>
      } />
      <Route path="/items/:itemID" element={
        <RequireAuth>
          <ItemDetail />
        </RequireAuth>
      } />
      <Route
        path="/user"
        element={
          <RequireAuth>
            <User />
          </RequireAuth>
        }
      />
      <Route path="/login" element={<Login />} />
      <Route path="/callback" element={<CallbackPage />} />
      <Route path="*" element={<NotFound />} />
    </Routes>
  );
}

