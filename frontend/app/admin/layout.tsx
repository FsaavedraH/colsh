import Sidebar from "@/components/layout/Sidebar";
import RutaProtegida from "@/components/layout/RutaProtegida";

const itemsAdmin = [
  { label: "Usuarios", href: "/admin" },
  { label: "Reportes", href: "/admin/reportes" },
];

export default function AdminLayout({ children }: { children: React.ReactNode }) {
  return (
    <RutaProtegida rolPermitido="Administrador">
      <div className="flex min-h-screen">
        <Sidebar items={itemsAdmin} colorRol="#0f172a" />
        <main className="flex-1 p-6 bg-gray-50">{children}</main>
      </div>
    </RutaProtegida>
  );
}