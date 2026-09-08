-- Migracion: rastrea el progreso de recoleccion/empaque por cada item de un pedido,
-- independiente del estado general del pedido (que solo cambia cuando TODOS los
-- items de la etapa correspondiente ya fueron procesados).

CREATE TABLE progreso_item_pedido (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    id_pedido uuid NOT NULL REFERENCES pedido(id_pedido) ON DELETE CASCADE,
    id_producto uuid NOT NULL REFERENCES producto(id_producto),
    etapa varchar(20) NOT NULL CHECK (etapa IN ('picking', 'empaque')),
    cantidad_procesada integer,
    fecha timestamp NOT NULL DEFAULT now(),
    responsable uuid NOT NULL REFERENCES usuario(id_usuario),
    UNIQUE (id_pedido, id_producto, etapa)
);