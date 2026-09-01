// Tipos que reflejan las structs del backend (backend-runtime/models.go)

export interface Actividad {
  id: number
  titulo: string
  descripcion: string | null
  categoria_id: number | null
  categoria_nombre: string | null
  prioridad: number | null
  activa: boolean
  creado_en: string
  actualizado_en: string
}

export interface CrearActividadInput {
  titulo: string
  descripcion?: string | null
  categoria_id?: number | null
  prioridad?: number | null
}

export interface CheckDiario {
  id: number
  actividad_id: number
  fecha: string
  nota: string | null
  creado_en: string
}

export interface CrearCheckInput {
  actividad_id: number
  nota?: string | null
}

export interface ResumenActividad {
  id: number
  titulo: string
  veces_hecha: number
  veces_no_hecha: number
  porcentaje_cumplimiento: number | null
}

export interface Categoria {
  id: number
  nombre: string
}

export interface CrearCategoriaInput {
  nombre: string
}

export interface ActualizarCategoriaInput {
  nombre: string
}
