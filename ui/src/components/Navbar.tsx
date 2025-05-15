import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { useAuth } from "@/context/AuthContext";
import { cn } from "@/lib/utils";
import { jwtDecode } from "jwt-decode";
import { Menu, X } from "lucide-react";
import { useState } from "react";
import { Link, Navigate } from "react-router-dom";

export function Navbar() {
  const { token } = useAuth();
  const [mobileOpen, setMobileOpen] = useState(false);

  if (!token) {
    return <Navigate to="/login" replace />;
  }

  let username: string | null = null;
  try {
    const decoded = jwtDecode<{ name?: string }>(token);
    username = decoded.name ?? null;
  } catch {
    username = null;
  }

  const navItems = [
    { to: "/feeds", label: "Feeds" },
    { to: "/suggested", label: "Suggested" },
    { to: "/collections", label: "Collections" },
  ];

  return (
    <nav className="sticky top-0 z-50 w-full bg-background border-b">
      <div className="container mx-auto flex h-16 items-center justify-between px-4">
        <Link to="/" className="flex items-center">
          <img src="/logo512.png" alt="Gazette" className="h-8 w-8 mr-2" />
          <span className="text-2xl font-bold">Gazette</span>
        </Link>

        <div className="hidden md:flex space-x-2">
          {navItems.map((item) => (
            <Link
              key={item.to}
              to={item.to}
              className={cn(
                "px-3 py-2 rounded-md text-sm font-medium text-muted-foreground",
                "hover:bg-muted hover:text-foreground",
              )}
            >
              {item.label}
            </Link>
          ))}
        </div>

        <div className="hidden md:flex items-center space-x-2">
          <Link to="/user" className="flex items-center space-x-2">
            <Avatar>
              {username ? (
                <AvatarFallback>
                  {username.charAt(0).toUpperCase()}
                </AvatarFallback>
              ) : (
                <AvatarImage
                  src="https://github.com/shadcn.png"
                  alt="User Avatar"
                />
              )}
            </Avatar>
          </Link>
        </div>

        <button
          className="md:hidden p-2"
          onClick={() => setMobileOpen(!mobileOpen)}
          aria-label="Toggle menu"
        >
          {mobileOpen ? (
            <X className="w-6 h-6" />
          ) : (
            <Menu className="w-6 h-6" />
          )}
        </button>
      </div>

      {mobileOpen && (
        <div className="md:hidden bg-background border-t">
          <div className="flex flex-col px-4 py-2 space-y-1">
            {navItems.map((item) => (
              <Link
                key={item.to}
                to={item.to}
                onClick={() => setMobileOpen(false)}
                className={cn(
                  "block px-3 py-2 rounded-md text-base font-medium text-muted-foreground",
                  "hover:bg-muted hover:text-foreground",
                )}
              >
                {item.label}
              </Link>
            ))}
            <Link
              to="/user"
              onClick={() => setMobileOpen(false)}
              className={cn(
                "flex items-center px-3 py-2 rounded-md text-base font-medium text-muted-foreground",
                "hover:bg-muted hover:text-foreground",
              )}
            >
              User
            </Link>
          </div>
        </div>
      )}
    </nav>
  );
}
