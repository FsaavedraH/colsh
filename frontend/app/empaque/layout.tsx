import Sidebar from "@/components/layout/Sidebar";

const itemsEmpaque = [
  { label: "Recepción", href: "/empaque" },
  { label: "Historial", href: "/empaque/historial" },
];

export default function EmpaqueLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex min-h-screen">
      <Sidebar items={itemsEmpaque} colorRol="#0d9488" />
      <main className="flex-1 p-6 bg-gray-50">{children}</main>
    </div>
  );
}