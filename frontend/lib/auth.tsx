"use client";

import { createContext, useContext, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { apiFetch } from "./api";

interface Usuario {
  id_usuario: string;
  nombre: string;
  email: string;
  rol: string;
}

interface AuthContextType {
  usuario: Usuario | null;
  cargando: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [usuario, setUsuario] = useState<Usuario | null>(null);
  const [cargando, setCargando] = useState(true);
  const router = useRouter();

  useEffect(() => {
    cargarDesdeStorage();
    setCargando(false);

    // Si el usuario inicia o cierra sesion en OTRA pestana del mismo navegador,
    // esta pestana se entera de inmediato y actualiza el nombre/rol mostrado,
    // en vez de quedarse con la sesion vieja que tenia al momento de abrirse.
    function alCambiarStorage(evento: StorageEvent) {
      if (evento.key === "colsh_usuario") {
        cargarDesdeStorage();
      }
    }

    window.addEventListener("storage", alCambiarStorage);
    return () => window.removeEventListener("storage", alCambiarStorage);
  }, []);

  function cargarDesdeStorage() {
    const guardado = localStorage.getItem("colsh_usuario");
    setUsuario(guardado ? JSON.parse(guardado) : null);
  };

  async function login(email: string, password: string) {
    const data = await apiFetch<Usuario>("/api/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });
    setUsuario(data);
    localStorage.setItem("colsh_usuario", JSON.stringify(data));

    const rutasPorRol: Record<string, string> = {
      Cliente: "/cliente",
      Picking: "/picking",
      Empaque: "/empaque",
      Transportista: "/transportista",
      Administrador: "/admin",
    };
    router.push(rutasPorRol[data.rol] || "/");
  }

  function logout() {
    setUsuario(null);
    localStorage.removeItem("colsh_usuario");
    router.push("/login");
  }

  return (
    <AuthContext.Provider value={{ usuario, cargando, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) throw new Error("useAuth debe usarse dentro de AuthProvider");
  return context;
}