package main

import "time"

// Actividad refleja la tabla actividades (ver migrations/*_baseline-esquema-inicial.sql)
type Actividad struct {
	ID             int32     `json:"id" db:"id"`
	Titulo         string    `json:"titulo" db:"titulo"`
	Descripcion    *string   `json:"descripcion" db:"descripcion"`
	Categoria      *string   `json:"categoria" db:"categoria"`
	Prioridad      *int16    `json:"prioridad" db:"prioridad"`
	Activa         bool      `json:"activa" db:"activa"`
	CreadoEn       time.Time `json:"creado_en" db:"creado_en"`
	ActualizadoEn  time.Time `json:"actualizado_en" db:"actualizado_en"`
}

// CrearActividadInput es el body esperado en POST /actividades
type CrearActividadInput struct {
	Titulo      string  `json:"titulo"`
	Descripcion *string `json:"descripcion"`
	Categoria   *string `json:"categoria"`
	Prioridad   *int16  `json:"prioridad"`
}

// CheckDiario refleja la tabla checks_diarios
type CheckDiario struct {
	ID          int32     `json:"id" db:"id"`
	ActividadID int32     `json:"actividad_id" db:"actividad_id"`
	Fecha       time.Time `json:"fecha" db:"fecha"`
	Nota        *string   `json:"nota" db:"nota"`
	CreadoEn    time.Time `json:"creado_en" db:"creado_en"`
}

// CrearCheckInput es el body esperado en POST /checks
type CrearCheckInput struct {
	ActividadID int32   `json:"actividad_id"`
	Nota        *string `json:"nota"`
}

// ResumenActividad refleja la vista resumen_actividades
type ResumenActividad struct {
	ID                     int32   `json:"id" db:"id"`
	Titulo                 string  `json:"titulo" db:"titulo"`
	VecesHecha             int64   `json:"veces_hecha" db:"veces_hecha"`
	VecesNoHecha           int64   `json:"veces_no_hecha" db:"veces_no_hecha"`
	PorcentajeCumplimiento *float64 `json:"porcentaje_cumplimiento" db:"porcentaje_cumplimiento"`
}
