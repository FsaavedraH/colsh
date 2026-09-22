import AppShell from "@/components/layout/AppShell";
import RutaProtegida from "@/components/layout/RutaProtegida";

const itemsPicking = [
  { label: "Órdenes de Picking", href: "/picking", icon: "scan" as const },
  { label: "Historial", href: "/picking/historial", icon: "history" as const },
];

export default function PickingLayout({ children }: { children: React.ReactNode }) {
  return (
    <RutaProtegida rolPermitido="Picking">
      <AppShell items={itemsPicking} accentColor="#d97706">
        {children}
      </AppShell>
    </RutaProtegida>
  );
}