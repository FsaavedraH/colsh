"use client";

import { createContext, useContext, useEffect, useRef, useState } from "react";
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

const MINUTOS_INACTIVIDAD = 15;
const CLAVE_ULTIMA_ACTIVIDAD = "colsh_ultima_actividad";

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [usuario, setUsuario] = useState<Usuario | null>(null);
  const [cargando, setCargando] = useState(true);
  const router = useRouter();
  const usuarioRef = useRef<Usuario | null>(null);

  useEffect(() => {
    usuarioRef.current = usuario;
  }, [usuario]);

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

    // Registra actividad real del usuario (clics, teclado, scroll, movimiento
    // del mouse) para reiniciar el conteo de inactividad. No cierra sesion por
    // si sola con el tiempo fijo: solo si pasan 15 min SIN ninguna de estas
    // interacciones.
    function registrarActividad() {
      if (usuarioRef.current) {
        localStorage.setItem(CLAVE_ULTIMA_ACTIVIDAD, Date.now().toString());
      }
    }
    const eventosActividad = ["click", "keydown", "scroll", "mousemove", "touchstart"];
    eventosActividad.forEach((evento) => window.addEventListener(evento, registrarActividad));
    registrarActividad();

    // Revisa cada 30 segundos si ya pasaron 15 minutos desde la ultima actividad.
    const verificador = setInterval(() => {
      if (!usuarioRef.current) return;
      const ultima = Number(localStorage.getItem(CLAVE_ULTIMA_ACTIVIDAD) || 0);
      const minutosInactivo = (Date.now() - ultima) / 1000 / 60;
      if (minutosInactivo >= MINUTOS_INACTIVIDAD) {
        cerrarSesionPorInactividad();
      }
    }, 30000);

    return () => {
      window.removeEventListener("storage", alCambiarStorage);
      eventosActividad.forEach((evento) => window.removeEventListener(evento, registrarActividad));
      clearInterval(verificador);
    };
  }, []);

  function cargarDesdeStorage() {
    const guardado = localStorage.getItem("colsh_usuario");
    setUsuario(guardado ? JSON.parse(guardado) : null);
  }

  function cerrarSesionPorInactividad() {
    setUsuario(null);
    localStorage.removeItem("colsh_usuario");
    localStorage.removeItem(CLAVE_ULTIMA_ACTIVIDAD);
    router.push("/login?motivo=inactividad");
  }

  async function login(email: string, password: string) {
    const data = await apiFetch<Usuario>("/api/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });
    setUsuario(data);
    localStorage.setItem("colsh_usuario", JSON.stringify(data));
    localStorage.setItem(CLAVE_ULTIMA_ACTIVIDAD, Date.now().toString());

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
    localStorage.removeItem(CLAVE_ULTIMA_ACTIVIDAD);
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