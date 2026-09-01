<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { listarActividades, crearActividad, crearCheck } from '@/services/api'
import type { Actividad } from '@/services/types'

const actividades = ref<Actividad[]>([])
const cargando = ref(true)
const error = ref<string | null>(null)

const nuevoTitulo = ref('')
const nuevaDescripcion = ref('')
const nuevaCategoria = ref('')
const nuevaPrioridad = ref(3)
const creando = ref(false)
const formAbierto = ref(false)

const checkEnCurso = ref<number | null>(null)
const checkeadasHoy = ref<Set<number>>(new Set())

async function cargarActividades() {
  cargando.value = true
  error.value = null
  try {
    actividades.value = await listarActividades()
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Error desconocido'
  } finally {
    cargando.value = false
  }
}

async function onCrearActividad() {
  if (!nuevoTitulo.value.trim()) return
  creando.value = true
  error.value = null
  try {
    await crearActividad({
      titulo: nuevoTitulo.value,
      descripcion: nuevaDescripcion.value || null,
      categoria: nuevaCategoria.value || null,
      prioridad: nuevaPrioridad.value,
    })
    nuevoTitulo.value = ''
    nuevaDescripcion.value = ''
    nuevaCategoria.value = ''
    nuevaPrioridad.value = 3
    formAbierto.value = false
    await cargarActividades()
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Error desconocido'
  } finally {
    creando.value = false
  }
}

async function onCheck(actividadId: number) {
  checkEnCurso.value = actividadId
  error.value = null
  try {
    await crearCheck({ actividad_id: actividadId })
    checkeadasHoy.value = new Set(checkeadasHoy.value).add(actividadId)
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Error desconocido'
  } finally {
    checkEnCurso.value = null
  }
}

function prioridadLabel(p: number | null) {
  if (!p) return null
  return ['', 'Muy baja', 'Baja', 'Media', 'Alta', 'Urgente'][p] ?? `P${p}`
}

onMounted(cargarActividades)
</script>

<template>
  <main class="page">
    <div class="page-header">
      <h1>Actividades</h1>
      <button class="btn" @click="formAbierto = !formAbierto">
        {{ formAbierto ? 'Cancelar' : '+ Nueva actividad' }}
      </button>
    </div>

    <form v-if="formAbierto" class="card form-actividad" @submit.prevent="onCrearActividad">
      <input v-model="nuevoTitulo" class="input" placeholder="Título" required autofocus />
      <input v-model="nuevaDescripcion" class="input" placeholder="Descripción (opcional)" />
      <div class="form-row">
        <input v-model="nuevaCategoria" class="input" placeholder="Categoría (opcional)" />
        <label class="prioridad-field">
          <span>Prioridad</span>
          <select v-model.number="nuevaPrioridad" class="input">
            <option v-for="p in 5" :key="p" :value="p">{{ p }} · {{ prioridadLabel(p) }}</option>
          </select>
        </label>
      </div>
      <button type="submit" class="btn" :disabled="creando">
        {{ creando ? 'Creando…' : 'Guardar actividad' }}
      </button>
    </form>

    <p v-if="error" class="banner banner-error">⚠️ {{ error }}</p>

    <p v-if="cargando" class="empty-state">Cargando actividades…</p>

    <div v-else-if="actividades.length === 0" class="empty-state">
      <p>Aún no tienes actividades.</p>
      <p>Crea la primera con el botón de arriba.</p>
    </div>

    <ul v-else class="lista">
      <li v-for="a in actividades" :key="a.id" class="card item">
        <span class="prioridad-dot" :data-nivel="a.prioridad ?? 0" />

        <div class="item-body">
          <div class="item-title-row">
            <strong>{{ a.titulo }}</strong>
            <span v-if="a.categoria" class="chip">{{ a.categoria }}</span>
          </div>
          <p v-if="a.descripcion" class="item-desc">{{ a.descripcion }}</p>
        </div>

        <button
          class="btn check-btn"
          :class="{ 'btn-secondary': checkeadasHoy.has(a.id) }"
          :disabled="checkEnCurso === a.id || checkeadasHoy.has(a.id)"
          @click="onCheck(a.id)"
        >
          <template v-if="checkeadasHoy.has(a.id)">✓ Hecho hoy</template>
          <template v-else-if="checkEnCurso === a.id">…</template>
          <template v-else>Check hoy</template>
        </button>
      </li>
    </ul>
  </main>
</template>

<style scoped>
.page {
  max-width: 720px;
  margin: 0 auto;
  padding: var(--space-8) var(--space-6);
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.page-header h1 {
  font-size: 1.5rem;
}

.form-actividad {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  padding: var(--space-6);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: var(--space-3);
}

.prioridad-field {
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
  font-size: 0.8rem;
  color: var(--color-text-muted);
}

.prioridad-field select {
  min-width: 11rem;
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

.prioridad-dot {
  flex-shrink: 0;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--color-text-subtle);
}

.prioridad-dot[data-nivel='4'],
.prioridad-dot[data-nivel='5'] {
  background: var(--color-danger);
}

.prioridad-dot[data-nivel='3'] {
  background: var(--color-accent);
}

.item-body {
  flex: 1;
  min-width: 0;
}

.item-title-row {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  flex-wrap: wrap;
}

.chip {
  font-size: 0.75rem;
  padding: 0.15rem 0.55rem;
  border-radius: 999px;
  background: var(--color-accent-soft);
  color: var(--color-accent);
  font-weight: 600;
}

.item-desc {
  margin-top: 0.2rem;
  color: var(--color-text-muted);
  font-size: 0.9rem;
}

.check-btn {
  flex-shrink: 0;
  white-space: nowrap;
}
</style>
