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
          <input v-model="nombreEditado" class="input" @keyup.enter="guardarEdicion(c.id)" />
          <div class="acciones">
            <button class="btn" :disabled="guardando" @click="guardarEdicion(c.id)">Guardar</button>
            <button class="btn btn-secondary" @click="cancelarEdicion">Cancelar</button>
          </div>
        </template>

        <template v-else>
          <span class="nombre">{{ c.nombre }}</span>
          <div class="acciones">
            <button class="btn btn-secondary" @click="empezarEdicion(c)">Renombrar</button>
            <button
              class="btn btn-secondary btn-danger"
              :disabled="eliminandoId === c.id"
              @click="onEliminar(c.id)"
            >
              {{ eliminandoId === c.id ? '...' : 'Eliminar' }}
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
  gap: var(--space-2);
  flex-shrink: 0;
}

.btn-danger {
  color: var(--color-danger);
  border-color: var(--color-danger-soft);
}

.btn-danger:hover:not(:disabled) {
  background: var(--color-danger-soft);
}
</style>
