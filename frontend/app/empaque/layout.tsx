import AppShell from "@/components/layout/AppShell";
import RutaProtegida from "@/components/layout/RutaProtegida";

const itemsEmpaque = [
  { label: "Recepción", href: "/empaque", icon: "box" as const },
  { label: "Historial", href: "/empaque/historial", icon: "history" as const },
];

export default function EmpaqueLayout({ children }: { children: React.ReactNode }) {
  return (
    <RutaProtegida rolPermitido="Empaque">
      <AppShell items={itemsEmpaque} accentColor="#0d9488">
        {children}
      </AppShell>
    </RutaProtegida>
  );
}