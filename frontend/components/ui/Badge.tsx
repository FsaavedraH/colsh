interface BadgeProps {
  children: React.ReactNode;
  color?: "green" | "red" | "yellow" | "gray";
}

export default function Badge({ children, color = "gray" }: BadgeProps) {
  const estilos = {
    green: "bg-green-100 text-green-800",
    red: "bg-red-100 text-red-800",
    yellow: "bg-yellow-100 text-yellow-800",
    gray: "bg-gray-100 text-gray-800",
  };

  return (
    <span className={`px-2 py-1 rounded-full text-xs font-semibold ${estilos[color]}`}>
      {children}
    </span>
  );
}