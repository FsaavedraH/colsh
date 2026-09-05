-- Migracion: agrega costo unitario al producto, y crea tabla de historial de compras

ALTER TABLE producto ADD COLUMN costo_unitario numeric(12,2);

CREATE TABLE compra_inventario (
    id_compra uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    id_producto uuid NOT NULL REFERENCES producto(id_producto) ON DELETE RESTRICT,
    cantidad integer NOT NULL CHECK (cantidad > 0),
    costo_unitario_momento numeric(12,2),
    costo_total numeric(12,2),
    fecha timestamp NOT NULL DEFAULT now(),
    responsable uuid NOT NULL REFERENCES usuario(id_usuario)
);