import { Button } from "@/components/ui/button";
import { LogOut as LogOutIcon, User as UserIcon } from "lucide-react";
import type { GithubComRhajizadaGazetteInternalServiceUser as UserModel } from "../api/data-contracts";

interface UserProfileProps {
  user: UserModel;
  onLogout: () => void;
}

export function UserProfile({ user, onLogout }: UserProfileProps) {
  return (
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
      <div className="w-full space-y-2 max-w-md">
        <p>
          <span className="font-medium">User ID:</span> {user.id}
        </p>
        <p>
          <span className="font-medium">Joined:</span>{" "}
          {user.createdAt ? new Date(user.createdAt).toLocaleDateString() : "—"}
        </p>
        <p>
          <span className="font-medium">Last Updated:</span>{" "}
          {user.lastUpdatedAt
            ? new Date(user.lastUpdatedAt).toLocaleDateString()
            : "—"}
        </p>
      </div>

      <Button
        variant="outline"
        size="sm"
        onClick={onLogout}
        className="flex items-center mt-2"
      >
        <LogOutIcon className="w-4 h-4 mr-2" /> Logout
      </Button>
    </div>
  );
}
