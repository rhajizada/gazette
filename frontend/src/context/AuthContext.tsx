import { jwtDecode } from "jwt-decode";
import {
  createContext,
  ReactNode,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { Api } from "../api/Api";

interface AuthContextType {
  token: string | null;
  api: Api<unknown>;
  login: (token: string) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  // initialize token from localStorage (if present)
  const [token, setToken] = useState<string | null>(
    () => window.localStorage.getItem("token")
  );

  // one shared Api client
  const api = useMemo(() => {
    return new Api<unknown>({
      baseUrl: import.meta.env.VITE_API_BASE_URL!,
      baseApiParams: { secure: true, format: "json" },
      securityWorker: () => {
        const t = window.localStorage.getItem("token");
        return t ? { headers: { Authorization: `Bearer ${t}` } } : {};
      },
    });
  }, []);

  // whenever token changes, inform the client
  useEffect(() => {
    api.setSecurityData(token);
    if (token) {
      window.localStorage.setItem("token", token);
    } else {
      window.localStorage.removeItem("token");
    }
  }, [token, api]);

  // auto‐logout when JWT expires
  useEffect(() => {
    if (!token) return;
    // assume JWT has exp claim
    const { exp } = jwtDecode<{ exp: number }>(token);
    const msUntilExpiry = exp * 1000 - Date.now();
    if (msUntilExpiry <= 0) {
      logout();
    } else {
      const timer = setTimeout(logout, msUntilExpiry);
      return () => clearTimeout(timer);
    }
  }, [token]);

  function login(newToken: string) {
    setToken(newToken);
  }

  function logout() {
    setToken(null);
    // force nav to login
    window.location.href = "/login";
  }

  return (
    <AuthContext.Provider value={{ token, api, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be inside AuthProvider");
  return ctx;
}
