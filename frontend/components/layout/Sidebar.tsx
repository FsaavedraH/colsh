"use client";

import { useAuth } from "@/lib/auth";

interface SidebarProps {
  items: { label: string; href: string }[];
  colorRol?: string;
}

export default function Sidebar({ items, colorRol = "#1e3a5f" }: SidebarProps) {
  const { usuario, logout } = useAuth();

  return (
    <nav
      style={{ backgroundColor: colorRol }}
      className="w-56 min-h-screen text-white p-4 flex flex-col"
    >
      <div className="font-bold text-lg mb-6">ColSh</div>
      <ul className="space-y-2 flex-1">
        {items.map((item) => (
          <li key={item.href}>
            <a href={item.href} className="block px-3 py-2 rounded hover:bg-white/10 transition-colors">
              {item.label}
            </a>
          </li>
        ))}
      </ul>
      {usuario && (
        <div className="border-t border-white/20 pt-4 mt-4">
          <p className="text-xs text-white/70 mb-2 truncate">{usuario.nombre}</p>
          <button
            onClick={logout}
            className="text-sm text-white/80 hover:text-white underline"
          >
            Cerrar sesión
          </button>
        </div>
      )}
    </nav>
  );
}