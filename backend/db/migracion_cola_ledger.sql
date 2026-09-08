-- Migracion: cola de eventos que no se pudieron registrar en el ledger al momento
-- del evento (ej. Hyperledger Fabric caido), para reintentar automaticamente mas
-- adelante SIN perder la hora original en que ocurrio el evento real.

CREATE TABLE eventos_ledger_pendientes (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    id_evento_original varchar(100) NOT NULL,
    id_pedido uuid NOT NULL,
    estado varchar(50) NOT NULL,
    fecha_evento_real timestamp NOT NULL,
    responsable varchar(100) NOT NULL,
    intentos integer NOT NULL DEFAULT 0,
    ultimo_intento timestamp,
    creado_en timestamp NOT NULL DEFAULT now()
);