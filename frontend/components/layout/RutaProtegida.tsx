"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/lib/auth";

interface RutaProtegidaProps {
  rolPermitido: string;
  children: React.ReactNode;
}

export default function RutaProtegida({ rolPermitido, children }: RutaProtegidaProps) {
  const { usuario, cargando } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (cargando) return;

    if (!usuario) {
      router.push("/login");
      return;
    }

    if (usuario.rol !== rolPermitido) {
      const rutasPorRol: Record<string, string> = {
        Cliente: "/cliente",
        Picking: "/picking",
        Empaque: "/empaque",
        Transportista: "/transportista",
        Administrador: "/admin",
      };
      router.push(rutasPorRol[usuario.rol] || "/login");
    }
  }, [usuario, cargando, rolPermitido, router]);

  if (cargando) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <p className="text-gray-400 text-sm">Cargando...</p>
      </div>
    );
  }

  if (!usuario || usuario.rol !== rolPermitido) {
    return null;
  }

  return <>{children}</>;
}