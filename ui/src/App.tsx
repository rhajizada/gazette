import { Route, Routes } from "react-router-dom";
import RequireAuth from "./components/RequireAuth";
import CallbackPage from "./pages/Callback";
import SubscribedItems from "./pages/SubscribedItems.tsx";
import SuggestedItems from "./pages/SuggestedItems.tsx";
import CategoriesPage from "./pages/Categories.tsx";
import CollectionDetails from "./pages/CollectionItems.tsx";
import Collections from "./pages/Collections.tsx";
import FeedItems from "./pages/FeedItems.tsx";
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
            <SubscribedItems />
          </RequireAuth>
        }
      />
      <Route path="/categories" element={<CategoriesPage />} />
      <Route
        path="/collections"
        element={
          <RequireAuth>
            <Collections />
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
        path="/feeds"
        element={
          <RequireAuth>
            <Feeds />
          </RequireAuth>
        }
      />
      <Route
        path="/feeds/:feedID"
        element={
          <RequireAuth>
            <FeedItems />
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
        path="/suggested"
        element={
          <RequireAuth>
            <SuggestedItems />
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
