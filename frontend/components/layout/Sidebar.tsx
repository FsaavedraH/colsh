interface SidebarProps {
  items: { label: string; href: string }[];
  colorRol?: string;
}

export default function Sidebar({ items, colorRol = "#1e3a5f" }: SidebarProps) {
  return (
    <nav
      style={{ backgroundColor: colorRol }}
      className="w-56 min-h-screen text-white p-4"
    >
      <div className="font-bold text-lg mb-6">ColSh</div>
      <ul className="space-y-2">
        {items.map((item) => (
          <li key={item.href}>
            <a href={item.href} className="block px-3 py-2 rounded hover:bg-white/10 transition-colors">
              {item.label}
            </a>
          </li>
        ))}
      </ul>
    </nav>
  );
}