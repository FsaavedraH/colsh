export type ColorBadge = "green" | "red" | "yellow" | "gray";

export function colorPorEstado(estado: string): ColorBadge {
  if (estado === "Entregado") return "green";
  if (estado === "Cancelado") return "red";
  return "yellow";
}