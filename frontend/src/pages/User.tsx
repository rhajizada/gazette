import { Footer } from "@/components/Footer";
import { Navbar } from "@/components/Navbar";
import { Button } from "@/components/ui/button";
import { Spinner } from "@/components/ui/spinner";
import { LogOut as LogOutIcon, User as UserIcon } from "lucide-react";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import type { GithubComRhajizadaGazetteInternalServiceUser as UserModel } from "../api/data-contracts";
import { useAuth } from "../context/AuthContext";

export default function UserPage() {
  const { api, logout } = useAuth();
  const navigate = useNavigate();
  const [user, setUser] = useState<UserModel | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setLoading(true);
    api
      .userList({ format: "json" })
      .then((res) => setUser(res.data))
      .catch((err) => {
        console.error(err);
        setError("Failed to load user data.");
        if (err.error === "Unauthorized") {
          logout();
          navigate("/login", { replace: true });
        }
      })
      .finally(() => setLoading(false));
  }, [api, logout, navigate]);

  if (loading) {
    return (
      <div className="min-h-screen flex flex-col">
        <Navbar />
        <div className="flex-1 flex items-center justify-center">
          <Spinner />
        </div>
        <Footer />
      </div>
    );
  }

  if (error || !user) {
    return (
      <div className="min-h-screen flex flex-col">
        <Navbar />
        <div className="flex-1 p-6 text-red-600">
          {error || "User not found."}
        </div>
        <Footer />
      </div>
    );
  }

  return (
    <div className="min-h-screen flex flex-col">
      <Navbar />
      <main className="flex-1 mx-auto max-w-md px-6 py-8">
        <div className="flex flex-col items-center space-y-4">
          <div className="h-24 w-24 rounded-full bg-gray-200 flex items-center justify-center">
            {user.name ? (
              <span className="text-3xl font-bold text-gray-600">
                {user.name.charAt(0).toUpperCase()}
              </span>
            ) : (
              <UserIcon className="w-12 h-12 text-gray-500" />
            )}
          </div>
          <h1 className="text-2xl font-extrabold text-gray-900">
            {user.name || "Unnamed User"}
          </h1>
          <p className="text-gray-600">{user.email}</p>
          <div className="w-full border-t" />
          <div className="w-full space-y-2">
            <p>
              <span className="font-medium">User ID:</span> {user.id}
            </p>
            <p>
              <span className="font-medium">Joined:</span>{" "}
              {user.createdAt && new Date(user.createdAt).toLocaleDateString()}
            </p>
            <p>
              <span className="font-medium">Last Updated:</span>{" "}
              {user.lastUpdatedAt &&
                new Date(user.lastUpdatedAt).toLocaleDateString()}
            </p>
          </div>
          <Button
            variant="outline"
            size="sm"
            onClick={logout}
            className="flex items-center mt-4"
          >
            <LogOutIcon className="w-4 h-4 mr-2" /> Logout
          </Button>
        </div>
      </main>
      <Footer />
    </div>
  );
}
