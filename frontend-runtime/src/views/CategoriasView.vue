<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { listarCategorias, crearCategoria, actualizarCategoria, eliminarCategoria } from '@/services/api'
import type { Categoria } from '@/services/types'

const categorias = ref<Categoria[]>([])
const cargando = ref(true)
const error = ref<string | null>(null)

const nuevoNombre = ref('')
const creando = ref(false)

const editandoId = ref<number | null>(null)
const nombreEditado = ref('')
const guardando = ref(false)
const eliminandoId = ref<number | null>(null)

async function cargar() {
  cargando.value = true
  error.value = null
  try {
    categorias.value = await listarCategorias()
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Error desconocido'
  } finally {
    cargando.value = false
  }
}

async function onCrear() {
  if (!nuevoNombre.value.trim()) return
  creando.value = true
  error.value = null
  try {
    await crearCategoria({ nombre: nuevoNombre.value.trim() })
    nuevoNombre.value = ''
    await cargar()
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Error desconocido'
  } finally {
    creando.value = false
  }
}

function empezarEdicion(c: Categoria) {
  editandoId.value = c.id
  nombreEditado.value = c.nombre
}

function cancelarEdicion() {
  editandoId.value = null
}

async function guardarEdicion(id: number) {
  if (!nombreEditado.value.trim()) return
  guardando.value = true
  error.value = null
  try {
    await actualizarCategoria(id, { nombre: nombreEditado.value.trim() })
    editandoId.value = null
    await cargar()
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Error desconocido'
  } finally {
    guardando.value = false
  }
}

async function onEliminar(id: number) {
  eliminandoId.value = id
  error.value = null
  try {
    await eliminarCategoria(id)
    await cargar()
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Error desconocido'
  } finally {
    eliminandoId.value = null
  }
}

onMounted(cargar)
</script>

<template>
  <main class="page">
    <h1>Categorías</h1>

    <form class="card form-inline" @submit.prevent="onCrear">
      <input v-model="nuevoNombre" class="input" placeholder="Nueva categoría" required />
      <button type="submit" class="btn" :disabled="creando">
        {{ creando ? 'Creando…' : 'Agregar' }}
      </button>
    </form>

    <p v-if="error" class="banner banner-error">⚠️ {{ error }}</p>
    <p v-if="cargando" class="empty-state">Cargando categorías…</p>

    <div v-else-if="categorias.length === 0" class="empty-state">
      <p>No hay categorías todavía.</p>
    </div>

    <ul v-else class="lista">
      <li v-for="c in categorias" :key="c.id" class="card item">
        <template v-if="editandoId === c.id">
          <input
            v-model="nombreEditado"
            class="input"
            autofocus
            @keyup.enter="guardarEdicion(c.id)"
            @keyup.esc="cancelarEdicion"
          />
          <div class="acciones">
            <button class="btn" :disabled="guardando" @click="guardarEdicion(c.id)">Guardar</button>
            <button class="btn btn-secondary" @click="cancelarEdicion">Cancelar</button>
          </div>
        </template>

        <template v-else>
          <span class="nombre">{{ c.nombre }}</span>
          <div class="acciones">
            <button class="icon-btn" aria-label="Renombrar" title="Renombrar" @click="empezarEdicion(c)">
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
              class="icon-btn icon-btn-danger"
              aria-label="Eliminar"
              title="Eliminar"
              :disabled="eliminandoId === c.id"
              @click="onEliminar(c.id)"
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
        </template>
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

.page h1 {
  font-size: 1.5rem;
  font-weight: 800;
}

.form-inline {
  display: flex;
  gap: var(--space-3);
  padding: var(--space-4);
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
  justify-content: space-between;
  gap: var(--space-3);
  padding: var(--space-3) var(--space-5);
}

.nombre {
  font-weight: 600;
}

.acciones {
  display: flex;
  gap: var(--space-1);
  flex-shrink: 0;
}
</style>
