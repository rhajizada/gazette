import { Footer } from "@/components/Footer";
import { LikedItems } from "@/components/LikedItems";
import { Navbar } from "@/components/Navbar";
import { SubscribedFeeds } from "@/components/SubscribedFeeds";
import { Spinner } from "@/components/ui/spinner";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { UserProfile } from "@/components/UserProfile";
import {
  Heart as HeartIcon,
  Rss as RssIcon,
  User as UserIcon,
} from "lucide-react";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { toast } from "sonner";
import type { GithubComRhajizadaGazetteInternalServiceUser as UserModel } from "../api/data-contracts";
import { useAuth } from "../context/AuthContext";

const tabs = [
  { name: "Profile", value: "profile", icon: UserIcon },
  { name: "Subscribed feeds", value: "subscribed", icon: RssIcon },
  { name: "Liked items", value: "liked", icon: HeartIcon },
];

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
        if (err.error === "Unauthorized") {
          logout();
          navigate("/login", { replace: true });
        } else {
          toast.error(err.text() || "Failed to fetch user");
          setError("Failed to load user data.");
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

      <main className="flex-1 bg-gray-50 py-6">
        <div className="container mx-auto px-4">
          <Tabs defaultValue="profile" className="w-full">
            {/* Responsive tab list: icons only on small, text+icon on sm+ */}
            <TabsList className="flex w-full bg-white rounded-t-lg shadow-sm overflow-hidden">
              {tabs.map((tab) => (
                <TabsTrigger
                  key={tab.value}
                  value={tab.value}
                  className="flex-1 flex items-center justify-center gap-1 py-3 px-3 sm:py-4 sm:px-6
                             text-sm sm:text-lg font-medium text-gray-600 hover:text-gray-800 hover:bg-gray-50
                             data-[state=active]:text-gray-900 data-[state=active]:border-b-4
                             data-[state=active]:border-primary bg-transparent"
                >
                  <tab.icon className="h-6 w-6" />
                  <span className="hidden sm:inline">{tab.name}</span>
                </TabsTrigger>
              ))}
            </TabsList>

            <div className="bg-white rounded-b-lg shadow-sm p-6">
              <TabsContent value="profile">
                <UserProfile user={user!} onLogout={logout} />
              </TabsContent>

              <TabsContent value="subscribed">
                <SubscribedFeeds />
              </TabsContent>

              <TabsContent value="liked">
                <LikedItems />
              </TabsContent>
            </div>
          </Tabs>
        </div>
      </main>

      <Footer />
    </div>
  );
}
