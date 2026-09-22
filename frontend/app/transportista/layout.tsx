import AppShell from "@/components/layout/AppShell";
import RutaProtegida from "@/components/layout/RutaProtegida";

const itemsTransportista = [
  { label: "Mis despachos", href: "/transportista", icon: "truck" as const },
  { label: "Historial", href: "/transportista/historial", icon: "history" as const },
];

export default function TransportistaLayout({ children }: { children: React.ReactNode }) {
  return (
    <RutaProtegida rolPermitido="Transportista">
      <AppShell items={itemsTransportista} accentColor="#b45309">
        {children}
      </AppShell>
    </RutaProtegida>
  );
}