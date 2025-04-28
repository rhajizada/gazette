import { ReactNode } from "react";
import { Navigate, useLocation } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

export default function RequireAuth({ children }: { children: ReactNode }) {
  const { token } = useAuth();
  const loc = useLocation();
  if (!token) {
    // send to login, preserve where they were going
    return <Navigate to="/login" state={{ from: loc }} replace />;
  }
  return <>{children}</>;
}
