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

const checkEnCurso = ref<number | null>(null)
const mensaje = ref<string | null>(null)

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
    await cargarActividades()
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Error desconocido'
  } finally {
    creando.value = false
  }
}

async function onCheck(actividadId: number) {
  checkEnCurso.value = actividadId
  mensaje.value = null
  error.value = null
  try {
    await crearCheck({ actividad_id: actividadId })
    mensaje.value = '¡Check registrado!'
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Error desconocido'
  } finally {
    checkEnCurso.value = null
  }
}

onMounted(cargarActividades)
</script>

<template>
  <main>
    <h1>Actividades</h1>

    <form class="form-actividad" @submit.prevent="onCrearActividad">
      <h2>Nueva actividad</h2>
      <input v-model="nuevoTitulo" placeholder="Título" required />
      <input v-model="nuevaDescripcion" placeholder="Descripción (opcional)" />
      <input v-model="nuevaCategoria" placeholder="Categoría (opcional)" />
      <label>
        Prioridad (1-5)
        <input v-model.number="nuevaPrioridad" type="number" min="1" max="5" />
      </label>
      <button type="submit" :disabled="creando">
        {{ creando ? 'Creando...' : 'Crear actividad' }}
      </button>
    </form>

    <p v-if="error" class="error">{{ error }}</p>
    <p v-if="mensaje" class="ok">{{ mensaje }}</p>

    <p v-if="cargando">Cargando actividades...</p>
    <p v-else-if="actividades.length === 0">No hay actividades todavía.</p>

    <ul v-else class="lista">
      <li v-for="a in actividades" :key="a.id">
        <div>
          <strong>{{ a.titulo }}</strong>
          <span v-if="a.categoria"> · {{ a.categoria }}</span>
          <span v-if="a.prioridad"> · prioridad {{ a.prioridad }}</span>
          <p v-if="a.descripcion">{{ a.descripcion }}</p>
        </div>
        <button :disabled="checkEnCurso === a.id" @click="onCheck(a.id)">
          {{ checkEnCurso === a.id ? '...' : 'Check hoy' }}
        </button>
      </li>
    </ul>
  </main>
</template>

<style scoped>
main {
  max-width: 640px;
  margin: 0 auto;
  padding: 1rem;
}

.form-actividad {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-bottom: 2rem;
}

.form-actividad input {
  padding: 0.4rem;
}

.lista {
  list-style: none;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.lista li {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
  padding: 0.75rem;
  border: 1px solid var(--color-border);
  border-radius: 6px;
}

.error {
  color: #d33;
}

.ok {
  color: #2a2;
}
</style>
