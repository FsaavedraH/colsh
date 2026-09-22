import AppShell from "@/components/layout/AppShell";
import RutaProtegida from "@/components/layout/RutaProtegida";

const itemsAdmin = [
  { label: "Usuarios", href: "/admin", icon: "users" as const },
  { label: "Inventario", href: "/admin/inventario", icon: "box" as const },
  { label: "Códigos QR", href: "/admin/codigos-qr", icon: "qrcode" as const },
  { label: "Ledger", href: "/admin/ledger", icon: "activity" as const },
  { label: "Reportes", href: "/admin/reportes", icon: "filetext" as const },
];

export default function AdminLayout({ children }: { children: React.ReactNode }) {
  return (
    <RutaProtegida rolPermitido="Administrador">
      <AppShell items={itemsAdmin} accentColor="#0f172a">
        {children}
      </AppShell>
    </RutaProtegida>
  );
}