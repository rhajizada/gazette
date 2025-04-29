import { Link, Navigate } from "react-router-dom"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { cn } from "@/lib/utils"
import { useAuth } from "@/context/AuthContext"
import { jwtDecode } from "jwt-decode"

export function Navbar() {
  const { token } = useAuth()

  if (!token) {
    return <Navigate to="/login" replace />
  }

  // Decode username from token
  let username: string | null = null
  try {
    const decoded = jwtDecode<{ name?: string }>(token)
    username = decoded.name ?? null
  } catch {
    username = null
  }

  const navItems = [
    { to: "/favorites", label: "Favorites" },
    { to: "/feeds", label: "Feeds" },
    { to: "/collections", label: "Collections" },
  ]

  return (
    <nav className="sticky top-0 z-50 w-full border-b bg-background">
      <div className="container mx-auto flex h-16 items-center justify-between px-4">
        {/* Logo */}
        <Link to="/" className="flex items-center">
          <img
            src="/logo512.png"
            alt="Gazette"
            className="h-8 w-8 mr-2"
          />
          <span className="text-2xl font-bold">Gazette</span>
        </Link>

        <div className="flex space-x-4">
          {navItems.map((item) => (
            <Link
              key={item.to}
              to={item.to}
              className={cn(
                "px-3 py-2 rounded-md text-sm font-medium text-muted-foreground",
                "hover:bg-muted hover:text-foreground"
              )}
            >
              {item.label}
            </Link>
          ))}
        </div>

        <Link to="/user" className="flex items-center space-x-2">
          <Avatar>
            <AvatarImage
              src="https://github.com/shadcn.png"
              alt="User Avatar"
            />
            <AvatarFallback>U</AvatarFallback>
          </Avatar>
          {username && (
            <span className="text-sm font-medium text-foreground">
              {username}
            </span>
          )}
        </Link>
      </div>
    </nav>
  )
}

