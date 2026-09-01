-- Up Migration
-- Baseline: documenta el esquema que YA EXISTE en la base de datos.
-- Usa IF NOT EXISTS a propósito: si las tablas ya están creadas (como es el caso),
-- esta migración no hace nada; solo queda registrada en la tabla de control
-- "pgmigrations" para que las próximas migraciones partan de aquí.

CREATE TABLE IF NOT EXISTS actividades (
  id              SERIAL PRIMARY KEY,
  titulo          VARCHAR(150)     NOT NULL,
  descripcion     TEXT,
  categoria       VARCHAR(50),
  prioridad       SMALLINT         DEFAULT 3 CHECK (prioridad >= 1 AND prioridad <= 5),
  activa          BOOLEAN          NOT NULL DEFAULT TRUE,
  creado_en       TIMESTAMP        DEFAULT now(),
  actualizado_en  TIMESTAMP        DEFAULT now()
);

CREATE TABLE IF NOT EXISTS checks_diarios (
  id            SERIAL PRIMARY KEY,
  actividad_id  INTEGER NOT NULL REFERENCES actividades(id),
  fecha         DATE    NOT NULL DEFAULT CURRENT_DATE,
  nota          TEXT,
  creado_en     TIMESTAMP DEFAULT now(),
  UNIQUE (actividad_id, fecha)
);

CREATE OR REPLACE VIEW resumen_actividades AS
 SELECT a.id,
    a.titulo,
    count(c.id) AS veces_hecha,
    (((CURRENT_DATE - (a.creado_en)::date) + 1) - count(c.id)) AS veces_no_hecha,
    round((((count(c.id))::numeric / (NULLIF(((CURRENT_DATE - (a.creado_en)::date) + 1), 0))::numeric) * (100)::numeric), 1) AS porcentaje_cumplimiento
   FROM (actividades a
     LEFT JOIN checks_diarios c ON ((c.actividad_id = a.id)))
  WHERE (a.activa = true)
  GROUP BY a.id;

-- Down Migration
-- A propósito no se borra nada: esta es la migración baseline y las tablas
-- ya existían antes de tener migraciones. Revertir esto no debería destruir
-- datos de producción por accidente.
