-- Up Migration
-- Convierte "categoria" (texto libre en actividades) en su propia tabla,
-- con CRUD independiente y relación por FK.

CREATE TABLE IF NOT EXISTS categorias (
  id     SERIAL PRIMARY KEY,
  nombre VARCHAR(50) UNIQUE NOT NULL
);

-- Migra los valores de texto que ya existan a la tabla nueva
INSERT INTO categorias (nombre)
SELECT DISTINCT categoria FROM actividades
WHERE categoria IS NOT NULL
ON CONFLICT (nombre) DO NOTHING;

-- Nueva columna de relación
ALTER TABLE actividades ADD COLUMN IF NOT EXISTS categoria_id INTEGER REFERENCES categorias(id);

-- Backfill: enlaza cada actividad con su categoría ya migrada
UPDATE actividades a
SET categoria_id = c.id
FROM categorias c
WHERE a.categoria = c.nombre;

-- La columna de texto libre ya no se necesita
ALTER TABLE actividades DROP COLUMN IF EXISTS categoria;

-- Down Migration
ALTER TABLE actividades ADD COLUMN categoria VARCHAR(50);

UPDATE actividades a
SET categoria = c.nombre
FROM categorias c
WHERE a.categoria_id = c.id;

ALTER TABLE actividades DROP COLUMN categoria_id;

DROP TABLE categorias;
