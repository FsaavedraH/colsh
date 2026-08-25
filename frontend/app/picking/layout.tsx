import Sidebar from "@/components/layout/Sidebar";

const itemsPicking = [
  { label: "Órdenes de Picking", href: "/picking" },
  { label: "Historial", href: "/picking/historial" },
];

export default function PickingLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex min-h-screen">
      <Sidebar items={itemsPicking} colorRol="#d97706" />
      <main className="flex-1 p-6 bg-gray-50">{children}</main>
    </div>
  );
}