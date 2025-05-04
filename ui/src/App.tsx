import { Route, Routes } from "react-router-dom";
import RequireAuth from "./components/RequireAuth";
import CallbackPage from "./pages/Callback";
import CollectionDetails from "./pages/CollectionDetails.tsx";
import Collections from "./pages/Collections.tsx";
import FeedDetails from "./pages/FeedDetails.tsx";
import Feeds from "./pages/Feeds";
import ItemDetails from "./pages/ItemDetails.tsx";
import Login from "./pages/Login";
import NotFound from "./pages/NotFound";
import User from "./pages/User.tsx";

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
      <Route
        path="/collections"
        element={
          <RequireAuth>
            <Collections />
          </RequireAuth>
        }
      />
      <Route
        path="/feeds/:feedID"
        element={
          <RequireAuth>
            <FeedDetails />
          </RequireAuth>
        }
      />
      <Route
        path="/items/:itemID"
        element={
          <RequireAuth>
            <ItemDetails />
          </RequireAuth>
        }
      />
      <Route
        path="/collections/:collectionID"
        element={
          <RequireAuth>
            <CollectionDetails />
          </RequireAuth>
        }
      />
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
