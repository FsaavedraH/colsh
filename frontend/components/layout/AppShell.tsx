"use client";

import { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useAuth } from "@/lib/auth";

export interface NavItem {
  label: string;
  href: string;
  icon: keyof typeof ICONOS;
  requiresAuth?: boolean;
}

const ICONOS = {
  home: (
    <path d="M3 10.5 12 3l9 7.5M5 9.5V21h5v-6h4v6h5V9.5" strokeLinecap="round" strokeLinejoin="round" />
  ),
  list: (
    <path d="M8 6h13M8 12h13M8 18h13M3 6h.01M3 12h.01M3 18h.01" strokeLinecap="round" />
  ),
  scan: (
    <path d="M3 7V4a1 1 0 0 1 1-1h3M17 3h3a1 1 0 0 1 1 1v3M21 17v3a1 1 0 0 1-1 1h-3M7 21H4a1 1 0 0 1-1-1v-3M7 12h10" strokeLinecap="round" strokeLinejoin="round" />
  ),
  history: (
    <path d="M3 12a9 9 0 1 0 3-6.7M3 4v5h5M12 7v5l3 3" strokeLinecap="round" strokeLinejoin="round" />
  ),
  box: (
    <path d="M21 8 12 3 3 8l9 5 9-5ZM3 8v9l9 5M21 8v9l-9 5M12 13v9" strokeLinecap="round" strokeLinejoin="round" />
  ),
  users: (
    <path d="M17 21v-2a4 4 0 0 0-4-4H7a4 4 0 0 0-4 4v2M9 11a4 4 0 1 0 0-8 4 4 0 0 0 0 8ZM23 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75" strokeLinecap="round" strokeLinejoin="round" />
  ),
  qrcode: (
    <path d="M4 4h6v6H4V4Zm10 0h6v6h-6V4ZM4 14h6v6H4v-6Zm10 3h3m-3 3h6v-6h-3" strokeLinecap="round" strokeLinejoin="round" />
  ),
  activity: (
    <path d="M22 12h-4l-3 9-6-18-3 9H2" strokeLinecap="round" strokeLinejoin="round" />
  ),
  filetext: (
    <path d="M14 3H6a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9l-6-6ZM14 3v6h6M9 13h6M9 17h6" strokeLinecap="round" strokeLinejoin="round" />
  ),
  truck: (
    <path d="M1 3h13v13H1V3Zm13 5h4l3 3v5h-7V8ZM4.5 21a2 2 0 1 0 0-4 2 2 0 0 0 0 4Zm14 0a2 2 0 1 0 0-4 2 2 0 0 0 0 4Z" strokeLinecap="round" strokeLinejoin="round" />
  ),
};

function Icono({ nombre }: { nombre: keyof typeof ICONOS }) {
  return (
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      {ICONOS[nombre]}
    </svg>
  );
}

const ALTO_HEADER = "4rem";

function BloqueSesion({ accentColor }: { accentColor: string }) {
  const { usuario, logout } = useAuth();

  if (usuario) {
    return (
      <div className="border-t border-gray-100 pt-3 px-1">
        <p className="text-xs text-gray-400 mb-2 truncate">{usuario.nombre}</p>
        <button
          onClick={logout}
          className="text-sm text-left px-2 py-1.5 rounded-lg text-red-600 hover:bg-red-50 w-full"
        >
          Cerrar sesión
        </button>
      </div>
    );
  }

  return (
    <div className="border-t border-gray-100 pt-3 px-1 space-y-1">
      <Link
        href="/login"
        className="block text-sm text-center px-3 py-2 rounded-lg text-white font-medium"
        style={{ backgroundColor: accentColor }}
      >
        Iniciar sesión
      </Link>
      <Link
        href="/registro"
        className="block text-sm text-center px-3 py-2 rounded-lg text-gray-700 hover:bg-gray-100"
      >
        Crear cuenta
      </Link>
    </div>
  );
}

export default function AppShell({
  items,
  accentColor = "#111827",
  children,
}: {
  items: NavItem[];
  accentColor?: string;
  children: React.ReactNode;
}) {
  const { usuario } = useAuth();
  const pathname = usePathname();
  const [menuAbierto, setMenuAbierto] = useState(false);

  const itemsVisibles = items.filter((item) => !item.requiresAuth || usuario);

  const Nav = ({ enDrawer = false }: { enDrawer?: boolean }) => (
    <nav className="flex flex-col gap-1 px-2">
      {itemsVisibles.map((item) => {
        const activo = pathname === item.href;
        return (
          <Link
            key={item.href}
            href={item.href}
            onClick={() => enDrawer && setMenuAbierto(false)}
            className={`flex items-center gap-3 text-sm px-3 py-2 rounded-lg transition-colors ${
              activo ? "text-white" : "text-gray-700 hover:bg-gray-100"
            }`}
            style={activo ? { backgroundColor: accentColor } : undefined}
          >
            <Icono nombre={item.icon} />
            {item.label}
          </Link>
        );
      })}
    </nav>
  );

  return (
    <div className="min-h-screen flex flex-col bg-gray-50">
      <header
        className="sticky top-0 z-20 h-16 shrink-0 bg-white border-b border-gray-200 px-4 flex items-center gap-4"
      >
        <button
          className="md:hidden p-2 -ml-2 text-gray-700"
          onClick={() => setMenuAbierto((v) => !v)}
          aria-label="Abrir menu"
        >
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <path d="M4 6h16M4 12h16M4 18h16" strokeLinecap="round" />
          </svg>
        </button>

        <Link href="/" className="font-bold text-lg shrink-0">
          ColSh
        </Link>

        <div className="hidden sm:flex flex-1 max-w-md mx-auto">
          <input
            type="text"
            placeholder="Buscar pedidos, repuestos, clientes..."
            className="w-full border border-gray-300 rounded-lg px-3 py-1.5 text-sm"
          />
        </div>

        <div className="ml-auto flex items-center gap-3">
          {usuario ? (
            <>
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="text-gray-500 hidden sm:block">
                <path d="M15 17h5l-1.4-1.4A2 2 0 0 1 18 14.2V11a6 6 0 1 0-12 0v3.2a2 2 0 0 1-.6 1.4L4 17h5m6 0a3 3 0 1 1-6 0" strokeLinecap="round" strokeLinejoin="round" />
              </svg>
              <div
                className="w-8 h-8 rounded-full text-white flex items-center justify-center text-xs font-semibold shrink-0"
                style={{ backgroundColor: accentColor }}
                title={usuario.nombre}
              >
                {usuario.nombre?.[0]?.toUpperCase() ?? "?"}
              </div>
            </>
          ) : (
            <>
              <Link href="/login" className="text-sm font-medium text-gray-700 hover:text-gray-900">
                Iniciar sesión
              </Link>
              <Link
                href="/registro"
                className="text-sm font-medium text-white px-3 py-1.5 rounded-lg"
                style={{ backgroundColor: accentColor }}
              >
                Crear cuenta
              </Link>
            </>
          )}
        </div>
      </header>

      <div className="flex flex-1 items-start">
        <aside
          className="hidden md:flex w-56 shrink-0 border-r border-gray-200 bg-white flex-col py-4 sticky overflow-y-auto"
          style={{ top: ALTO_HEADER, height: `calc(100vh - ${ALTO_HEADER})` }}
        >
          <Nav />
          <div className="mt-auto px-2">
            <BloqueSesion accentColor={accentColor} />
          </div>
        </aside>

        {menuAbierto && (
          <div className="fixed inset-0 z-30 md:hidden">
            <div className="absolute inset-0 bg-black/30" onClick={() => setMenuAbierto(false)} />
            <aside className="absolute left-0 top-0 h-full w-64 bg-white p-4 flex flex-col">
              <Nav enDrawer />
              <div className="mt-auto">
                <BloqueSesion accentColor={accentColor} />
              </div>
            </aside>
          </div>
        )}

        <main className="flex-1 p-4 md:p-6 min-w-0">{children}</main>
      </div>
    </div>
  );
}