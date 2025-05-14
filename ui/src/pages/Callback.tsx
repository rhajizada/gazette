import { useEffect } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

export default function CallbackPage() {
  const { login } = useAuth();
  const loc = useLocation();
  const nav = useNavigate();

  useEffect(() => {
    const params = new URLSearchParams(loc.search);
    const token = params.get("token");
    if (token) {
      login(token);
      nav("/", { replace: true });
    } else {
      nav("/login", { replace: true });
    }
  }, [loc.search, login, nav]);

  return <p>Signing you in…</p>;
}
