import type {
  Actividad,
  CrearActividadInput,
  CheckDiario,
  CrearCheckInput,
  ResumenActividad,
  Categoria,
  CrearCategoriaInput,
  ActualizarCategoriaInput,
} from './types'

const BASE_URL = import.meta.env.VITE_API_URL

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })

  if (!res.ok) {
    const body = await res.json().catch(() => null)
    throw new Error(body?.error ?? `Error ${res.status} en ${path}`)
  }

  if (res.status === 204) {
    return undefined as T
  }

  return res.json() as Promise<T>
}

export function listarActividades(): Promise<Actividad[]> {
  return request<Actividad[]>('/actividades')
}

export function crearActividad(input: CrearActividadInput): Promise<Actividad> {
  return request<Actividad>('/actividades', {
    method: 'POST',
    body: JSON.stringify(input),
  })
}

export function crearCheck(input: CrearCheckInput): Promise<CheckDiario> {
  return request<CheckDiario>('/checks', {
    method: 'POST',
    body: JSON.stringify(input),
  })
}

export function obtenerResumen(): Promise<ResumenActividad[]> {
  return request<ResumenActividad[]>('/resumen')
}

export function listarCategorias(): Promise<Categoria[]> {
  return request<Categoria[]>('/categorias')
}

export function crearCategoria(input: CrearCategoriaInput): Promise<Categoria> {
  return request<Categoria>('/categorias', {
    method: 'POST',
    body: JSON.stringify(input),
  })
}

export function actualizarCategoria(
  id: number,
  input: ActualizarCategoriaInput,
): Promise<Categoria> {
  return request<Categoria>(`/categorias/${id}`, {
    method: 'PUT',
    body: JSON.stringify(input),
  })
}

export function eliminarCategoria(id: number): Promise<void> {
  return request<void>(`/categorias/${id}`, { method: 'DELETE' })
}
