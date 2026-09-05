-- Limpieza de datos de prueba y carga del catalogo real (48 repuestos)
BEGIN;
TRUNCATE TABLE compra_inventario, intento_escaneo, evento_trazabilidad, detalle_pedido, pedido, inventario, producto CASCADE;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Llanta delantera', 150.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 50, 'B-01-01-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Llanta trasera', 180.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 100, 'B-01-02-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Pastillas de freno delanteras', 25.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 100, 'B-03-01-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Pastillas de freno traseras', 20.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 100, 'B-03-02-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Disco de freno delantero', 80.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 50, 'B-02-01-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Disco de freno trasero', 70.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 50, 'B-02-02-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Pinza de freno delantera', 120.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 20, 'B-02-03-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Bomba de freno', 60.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 20, 'B-02-04-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Pastillas de embrague', 30.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 100, 'B-03-03-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Disco de embrague', 50.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 50, 'B-02-05-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Kit de arrastre', 120.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 100, 'B-02-06-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Cadena de transmision', 60.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 100, 'B-02-07-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Pinon delantero', 25.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 100, 'B-02-08-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Corona trasera', 35.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 100, 'B-02-09-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Amortiguador trasero', 180.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 20, 'B-01-03-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Horquilla delantera', 250.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 20, 'B-01-04-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Botellas de horquilla', 100.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 20, 'B-01-05-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Muelles suspension trasera', 40.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 20, 'B-03-04-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Bujia', 15.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 100, 'B-03-05-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('CDI / unidad encendido', 80.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 20, 'B-03-06-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Bobina de alta', 40.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 20, 'B-03-07-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Alternador / estator', 90.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 20, 'B-03-08-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Regulador / rectificador', 35.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 50, 'B-03-09-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Motor de arranque', 200.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 20, 'B-03-10-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Bomba de aceite', 70.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 20, 'B-02-10-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Juntas de motor', 60.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 50, 'B-03-11-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Piston', 90.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 50, 'B-03-12-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Segmentos (anillos)', 25.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 100, 'B-03-13-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Bielas', 120.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 20, 'B-03-14-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Ciguenal', 180.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 20, 'B-02-11-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Carter', 80.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 20, 'B-01-06-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Manillar', 40.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 50, 'B-02-12-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Espejos', 15.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 100, 'B-03-15-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Guardabarros', 25.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 50, 'B-01-07-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Carenaje / tapas plasticas', 60.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 50, 'B-01-08-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Farol delantero', 50.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 50, 'B-02-13-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Intermitentes', 10.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 100, 'B-03-16-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Manetas freno/embrague', 20.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 100, 'B-03-17-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Pedal freno/cambio', 15.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 50, 'B-03-18-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Filtro de aire', 8.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 100, 'B-03-19-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Filtro de aceite', 12.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 100, 'B-03-20-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Vulcanizaciones / neumatico interno', 20.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 100, 'B-03-21-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Cable de embrague', 10.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 100, 'B-03-22-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Cable de acelerador', 10.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 100, 'B-03-23-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Cable de freno (mecanico)', 10.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 100, 'B-03-24-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Sensor de velocidad', 25.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 50, 'B-03-25-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Juego de tornilleria', 15.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 100, 'B-03-26-01' FROM nuevo;

WITH nuevo AS (
  INSERT INTO producto (nombre, costo_unitario) VALUES ('Protector / slider de motor', 30.00)
  RETURNING id_producto
)
INSERT INTO inventario (id_producto, stock, ubicacion)
SELECT id_producto, 50, 'B-02-14-01' FROM nuevo;

COMMIT;
