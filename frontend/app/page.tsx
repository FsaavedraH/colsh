"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/lib/auth";

const RUTAS_POR_ROL: Record<string, string> = {
  Cliente: "/cliente",
  Picking: "/picking",
  Empaque: "/empaque",
  Transportista: "/transportista",
  Administrador: "/admin",
};

export default function Home() {
  const { usuario, cargando } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (cargando) return;
       router.replace(usuario ? RUTAS_POR_ROL[usuario.rol] ?? "/cliente" : "/cliente");
  }, [usuario, cargando, router]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
           <p className="text-sm text-gray-400">Cargando ColSh...</p>
    </div>
  );
}