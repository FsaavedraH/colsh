import Sidebar from "@/components/layout/Sidebar";

const itemsCliente = [
  { label: "Catálogo", href: "/cliente" },
  { label: "Mis pedidos", href: "/cliente/pedidos" },
];

export default function ClienteLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex min-h-screen">
      <Sidebar items={itemsCliente} colorRol="#1e3a5f" />
      <main className="flex-1 p-6 bg-gray-50">{children}</main>
    </div>
  );
}