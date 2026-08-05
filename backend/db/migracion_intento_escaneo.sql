CREATE TABLE intento_escaneo (
    id_intento  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_pedido   UUID NOT NULL REFERENCES pedido(id_pedido) ON DELETE CASCADE,
    tipo        VARCHAR(20) NOT NULL CHECK (tipo IN ('ubicacion', 'producto')),
    resultado   VARCHAR(20) NOT NULL CHECK (resultado IN ('correcto', 'incorrecto')),
    etapa       VARCHAR(20) NOT NULL CHECK (etapa IN ('picking', 'empaque')),
    fecha       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_intento_escaneo_pedido ON intento_escaneo(id_pedido);
CREATE INDEX idx_intento_escaneo_fecha ON intento_escaneo(fecha);