import Sidebar from "@/components/layout/Sidebar";

const itemsTransportista = [
  { label: "Mis despachos", href: "/transportista" },
  { label: "Historial", href: "/transportista/historial" },
];

export default function TransportistaLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex min-h-screen">
      <Sidebar items={itemsTransportista} colorRol="#b45309" />
      <main className="flex-1 p-6 bg-gray-50">{children}</main>
    </div>
  );
}