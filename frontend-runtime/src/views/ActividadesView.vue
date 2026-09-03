<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import {
  listarActividades,
  crearActividad,
  actualizarActividad,
  eliminarActividad,
  crearCheck,
  listarCategorias,
  obtenerResumen,
  ApiError,
} from '@/services/api'
import type { Actividad, Categoria, ResumenActividad } from '@/services/types'

const actividades = ref<Actividad[]>([])
const categorias = ref<Categoria[]>([])
const resumen = ref<ResumenActividad[]>([])
const cargando = ref(true)
const error = ref<string | null>(null)

const nuevoTitulo = ref('')
const nuevaDescripcion = ref('')
const nuevaCategoriaId = ref<number | null>(null)
const nuevaPrioridad = ref(3)
const creando = ref(false)
const formAbierto = ref(false)
// Si tiene un id, el modal está en modo edición (PUT) en vez de creación (POST).
const editandoId = ref<number | null>(null)

const checkEnCurso = ref<number | null>(null)
const checkeadasHoy = ref<Set<number>>(new Set())
const eliminandoId = ref<number | null>(null)

const overallProgress = computed(() => {
  if (resumen.value.length === 0) return 0
  const total = resumen.value.reduce((sum, r) => sum + (r.porcentaje_cumplimiento ?? 0), 0)
  return Math.round(total / resumen.value.length)
})

// Fecha de hoy, formateada para mostrarse en la pantalla: el check es por día
// (el backend solo permite uno por actividad y fecha), así que conviene que
// se vea a qué día corresponde el estado "hecho" que se está marcando.
const hoyFormateado = computed(() => {
  const texto = new Intl.DateTimeFormat('es-PE', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
  }).format(new Date())
  return texto.charAt(0).toUpperCase() + texto.slice(1)
})

async function cargarActividades() {
  cargando.value = true
  error.value = null
  try {
    actividades.value = await listarActividades()
    // El backend ya nos dice, por actividad, si hoy tiene check registrado
    // (check_hoy) — sincronizamos el set local con ese estado real en vez de
    // depender solo de los clicks hechos en esta sesión.
    checkeadasHoy.value = new Set(actividades.value.filter((a) => a.check_hoy).map((a) => a.id))
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Error desconocido'
  } finally {
    cargando.value = false
  }
}

async function cargarCategorias() {
  try {
    categorias.value = await listarCategorias()
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Error desconocido'
  }
}

async function cargarResumen() {
  try {
    resumen.value = await obtenerResumen()
  } catch {
    // El panel de progreso es informativo: si falla, simplemente no se muestra.
  }
}

function limpiarFormulario() {
  nuevoTitulo.value = ''
  nuevaDescripcion.value = ''
  nuevaCategoriaId.value = null
  nuevaPrioridad.value = 3
}

function abrirFormulario() {
  editandoId.value = null
  limpiarFormulario()
  formAbierto.value = true
}

function abrirEdicion(a: Actividad) {
  editandoId.value = a.id
  nuevoTitulo.value = a.titulo
  nuevaDescripcion.value = a.descripcion ?? ''
  nuevaCategoriaId.value = a.categoria_id
  nuevaPrioridad.value = a.prioridad ?? 3
  formAbierto.value = true
}

function cerrarFormulario() {
  formAbierto.value = false
  editandoId.value = null
}

async function onGuardarActividad() {
  if (!nuevoTitulo.value.trim()) return
  creando.value = true
  error.value = null
  try {
    const datos = {
      titulo: nuevoTitulo.value,
      descripcion: nuevaDescripcion.value || null,
      categoria_id: nuevaCategoriaId.value,
      prioridad: nuevaPrioridad.value,
    }
    if (editandoId.value !== null) {
      await actualizarActividad(editandoId.value, datos)
    } else {
      await crearActividad(datos)
    }
    limpiarFormulario()
    formAbierto.value = false
    editandoId.value = null
    await cargarActividades()
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Error desconocido'
  } finally {
    creando.value = false
  }
}

async function onEliminar(a: Actividad) {
  if (!confirm(`¿Eliminar "${a.titulo}"? Esta acción no se puede deshacer desde la pantalla.`)) {
    return
  }
  eliminandoId.value = a.id
  error.value = null
  try {
    await eliminarActividad(a.id)
    await cargarActividades()
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Error desconocido'
  } finally {
    eliminandoId.value = null
  }
}

async function onCheck(actividadId: number) {
  if (checkeadasHoy.value.has(actividadId)) return
  checkEnCurso.value = actividadId
  error.value = null
  try {
    await crearCheck({ actividad_id: actividadId })
    checkeadasHoy.value = new Set(checkeadasHoy.value).add(actividadId)
  } catch (e) {
    if (e instanceof ApiError && e.status === 409) {
      // Ya existía un check para hoy (p. ej. otra pestaña se adelantó): no es
      // un error real, solo sincronizamos el estado local con el del backend.
      checkeadasHoy.value = new Set(checkeadasHoy.value).add(actividadId)
    } else {
      error.value = e instanceof Error ? e.message : 'Error desconocido'
    }
  } finally {
    checkEnCurso.value = null
  }
}

function prioridadInfo(p: number | null) {
  const nivel = p ?? 0
  if (nivel >= 4) return { label: 'Alta', bars: 3 }
  if (nivel === 3) return { label: 'Media', bars: 2 }
  if (nivel >= 1) return { label: 'Baja', bars: 1 }
  return { label: 'Sin definir', bars: 0 }
}

function prioridadLabel(p: number | null) {
  if (!p) return null
  return ['', 'Muy baja', 'Baja', 'Media', 'Alta', 'Urgente'][p] ?? `P${p}`
}

onMounted(() => {
  cargarActividades()
  cargarCategorias()
  cargarResumen()
})
</script>

<template>
  <main class="page">
    <div class="layout">
      <section class="main-col">
        <div class="today-row">
          <span class="today-dot" aria-hidden="true" />
          <span>
            Hoy · <strong>{{ hoyFormateado }}</strong>
          </span>
          <span class="today-hint">— cada actividad admite un check por día</span>
        </div>

        <p v-if="error" class="banner banner-error">⚠️ {{ error }}</p>

        <p v-if="cargando" class="empty-state">Cargando actividades…</p>

        <div v-else-if="actividades.length === 0" class="empty-state">
          <p>Aún no tienes actividades.</p>
          <p>Crea la primera con el botón "+ Nueva actividad".</p>
        </div>

        <ul v-else class="lista">
          <li v-for="a in actividades" :key="a.id" class="card item">
            <button
              type="button"
              class="check-circle"
              :class="{ 'is-checked': checkeadasHoy.has(a.id) }"
              :disabled="checkEnCurso === a.id || checkeadasHoy.has(a.id)"
              :aria-label="checkeadasHoy.has(a.id) ? `Hecho hoy (${hoyFormateado})` : `Marcar como hecho hoy (${hoyFormateado})`"
              :title="checkeadasHoy.has(a.id) ? `Hecho hoy · ${hoyFormateado}` : `Marcar como hecho hoy · ${hoyFormateado}`"
              @click="onCheck(a.id)"
            >
              <svg v-if="checkeadasHoy.has(a.id)" width="14" height="14" viewBox="0 0 14 14" fill="none">
                <path d="M2.5 7.2L5.4 10L11.5 3.8" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
              </svg>
              <span v-else-if="checkEnCurso === a.id" class="spinner" />
            </button>

            <div class="item-body">
              <strong class="item-title">{{ a.titulo }}</strong>
              <p v-if="a.descripcion" class="item-desc">{{ a.descripcion }}</p>
              <span v-if="a.categoria_nombre" class="pill">{{ a.categoria_nombre }}</span>
            </div>

            <div class="priority-indicator" :title="`Prioridad: ${prioridadInfo(a.prioridad).label}`">
              <span
                v-for="bar in 3"
                :key="bar"
                class="priority-bar"
                :style="{ height: `${bar * 4 + 4}px` }"
                :class="{ 'is-filled': bar <= prioridadInfo(a.prioridad).bars }"
              />
            </div>

            <div class="acciones">
              <button
                type="button"
                class="icon-btn"
                aria-label="Editar actividad"
                title="Editar actividad"
                @click="abrirEdicion(a)"
              >
                <svg width="17" height="17" viewBox="0 0 17 17" fill="none">
                  <path
                    d="M11.5 2.5L14.5 5.5L5.83 14.17L2.5 15L3.33 11.67L11.5 2.5Z"
                    stroke="currentColor"
                    stroke-width="1.4"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  />
                </svg>
              </button>
              <button
                type="button"
                class="icon-btn icon-btn-danger"
                aria-label="Eliminar actividad"
                title="Eliminar actividad"
                :disabled="eliminandoId === a.id"
                @click="onEliminar(a)"
              >
                <svg width="17" height="17" viewBox="0 0 17 17" fill="none">
                  <path
                    d="M3.5 4.5H13.5M6.5 4.5V2.83C6.5 2.37 6.87 2 7.33 2H9.67C10.13 2 10.5 2.37 10.5 2.83V4.5M7 7.5V12M10 7.5V12M4.5 4.5L5 13.33C5.04 13.98 5.58 14.5 6.23 14.5H10.77C11.42 14.5 11.96 13.98 12 13.33L12.5 4.5"
                    stroke="currentColor"
                    stroke-width="1.4"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  />
                </svg>
              </button>
            </div>
          </li>
        </ul>
      </section>

      <aside class="side-col">
        <div class="card summary-inset">
          <span class="summary-label">Overall Progress</span>
          <div class="summary-pct">{{ overallProgress }}%</div>
          <div class="progress-track">
            <div class="progress-fill" :style="{ width: `${overallProgress}%` }" />
          </div>
          <RouterLink to="/resumen" class="summary-link">Ver resumen completo →</RouterLink>
        </div>
      </aside>
    </div>

    <button type="button" class="fab" @click="abrirFormulario">
      <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
        <path d="M8 2.5V13.5M2.5 8H13.5" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
      </svg>
      Nueva actividad
    </button>

    <div v-if="formAbierto" class="modal-overlay" @click.self="cerrarFormulario">
      <form class="modal-card" @submit.prevent="onGuardarActividad">
        <div class="modal-header">
          <h2>{{ editandoId !== null ? 'Editar actividad' : 'Nueva actividad' }}</h2>
          <button type="button" class="icon-btn" aria-label="Cerrar" @click="cerrarFormulario">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
              <path d="M3 3L13 13M13 3L3 13" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" />
            </svg>
          </button>
        </div>

        <label class="field">
          <span>Título</span>
          <input v-model="nuevoTitulo" class="input" placeholder="p. ej. Practicar básquet" required autofocus />
        </label>

        <label class="field">
          <span>Descripción</span>
          <textarea v-model="nuevaDescripcion" class="input" placeholder="Descripción (opcional)" rows="3" />
        </label>

        <div class="form-row">
          <label class="field">
            <span>Categoría</span>
            <select v-model="nuevaCategoriaId" class="input">
              <option :value="null">Sin categoría</option>
              <option v-for="c in categorias" :key="c.id" :value="c.id">{{ c.nombre }}</option>
            </select>
          </label>
          <label class="field">
            <span>Prioridad</span>
            <select v-model.number="nuevaPrioridad" class="input">
              <option v-for="p in 5" :key="p" :value="p">{{ p }} · {{ prioridadLabel(p) }}</option>
            </select>
          </label>
        </div>

        <p v-if="categorias.length === 0" class="hint">
          No tienes categorías todavía — puedes crearlas en la sección
          <RouterLink to="/categorias">Categorías</RouterLink>.
        </p>

        <div class="modal-actions">
          <button type="button" class="btn btn-secondary" @click="cerrarFormulario">Cancelar</button>
          <button type="submit" class="btn" :disabled="creando">
            {{ creando ? 'Guardando…' : editandoId !== null ? 'Guardar cambios' : 'Guardar actividad' }}
          </button>
        </div>
      </form>
    </div>
  </main>
</template>

<style scoped>
.page {
  max-width: 960px;
  margin: 0 auto;
  padding: var(--space-8) var(--space-6) 7rem;
}

.layout {
  display: grid;
  grid-template-columns: 1fr 17rem;
  gap: var(--space-6);
  align-items: start;
}

.main-col {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
  min-width: 0;
}

.today-row {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-size: 0.85rem;
  color: var(--color-text-muted);
}

.today-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--gradient-accent);
  flex-shrink: 0;
}

.today-row strong {
  color: var(--color-text);
}

.today-hint {
  color: var(--color-text-subtle);
}

@media (max-width: 500px) {
  .today-hint {
    display: none;
  }
}

.side-col {
  position: sticky;
  top: calc(var(--space-8) + 3.5rem);
}

.summary-inset {
  padding: var(--space-5);
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.summary-label {
  font-size: 0.8rem;
  font-weight: 700;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.summary-pct {
  font-size: 2rem;
  font-weight: 800;
}

.progress-track {
  height: 8px;
  border-radius: var(--radius-full);
  background: var(--color-border);
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  border-radius: var(--radius-full);
  background: var(--gradient-accent);
  transition: width 0.3s ease;
}

.summary-link {
  margin-top: var(--space-1);
  font-size: 0.82rem;
  font-weight: 600;
}

.lista {
  list-style: none;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.item {
  display: flex;
  align-items: center;
  gap: var(--space-4);
  padding: var(--space-4) var(--space-5);
  transition: box-shadow 0.15s;
}

.item:hover {
  box-shadow: var(--shadow-md);
}

.check-circle {
  flex-shrink: 0;
  width: 1.6rem;
  height: 1.6rem;
  border-radius: 50%;
  border: 1.6px solid var(--color-border);
  background: var(--color-surface);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition:
    background-color 0.15s,
    border-color 0.15s,
    transform 0.1s;
}

.check-circle:hover:not(:disabled) {
  border-color: var(--color-accent);
}

.check-circle:active:not(:disabled) {
  transform: scale(0.92);
}

.check-circle.is-checked {
  background: var(--color-accent);
  border-color: var(--color-accent);
  cursor: default;
}

.spinner {
  width: 0.7rem;
  height: 0.7rem;
  border: 2px solid var(--color-border);
  border-top-color: var(--color-accent);
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.item-body {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
}

.item-title {
  font-weight: 700;
}

.item-desc {
  color: var(--color-text-muted);
  font-size: 0.9rem;
}

.pill {
  align-self: flex-start;
}

.priority-indicator {
  flex-shrink: 0;
  display: flex;
  align-items: flex-end;
  gap: 3px;
  height: 16px;
}

.priority-bar {
  width: 4px;
  border-radius: 2px;
  background: var(--color-neutral-soft-hover);
}

.priority-bar.is-filled {
  background: var(--color-text);
}

.acciones {
  display: flex;
  gap: var(--space-1);
  flex-shrink: 0;
}

.hint {
  font-size: 0.85rem;
  color: var(--color-text-muted);
}

.field {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--color-text-muted);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-3);
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-3);
  margin-top: var(--space-2);
}

@media (max-width: 760px) {
  .layout {
    grid-template-columns: 1fr;
  }

  .side-col {
    position: static;
  }

  .form-row {
    grid-template-columns: 1fr;
  }

  .fab {
    right: var(--space-4);
    bottom: var(--space-4);
  }
}
</style>
