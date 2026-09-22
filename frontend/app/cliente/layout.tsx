import AppShell from "@/components/layout/AppShell";

const itemsCliente = [
  { label: "Catálogo", href: "/cliente", icon: "list" as const },
  { label: "Mis pedidos", href: "/cliente/pedidos", icon: "box" as const, requiresAuth: true },
];

export default function ClienteLayout({ children }: { children: React.ReactNode }) {
  return (
    <AppShell items={itemsCliente} accentColor="#1e3a5f">
      {children}
    </AppShell>
  );
}