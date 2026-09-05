-- Migracion: permite el estado "Cancelado" en pedido
ALTER TABLE pedido DROP CONSTRAINT pedido_estado_check;

ALTER TABLE pedido ADD CONSTRAINT pedido_estado_check
CHECK (estado IN ('Pendiente', 'En espera por inventario', 'En recoleccion', 'En empaque', 'En despacho', 'Entregado', 'Cancelado'));