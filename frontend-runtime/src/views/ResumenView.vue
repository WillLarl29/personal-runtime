<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { obtenerResumen } from '@/services/api'
import type { ResumenActividad } from '@/services/types'

const resumen = ref<ResumenActividad[]>([])
const cargando = ref(true)
const error = ref<string | null>(null)

async function cargar() {
  cargando.value = true
  error.value = null
  try {
    resumen.value = await obtenerResumen()
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Error desconocido'
  } finally {
    cargando.value = false
  }
}

onMounted(cargar)
</script>

<template>
  <main>
    <h1>Resumen</h1>

    <p v-if="error" class="error">{{ error }}</p>
    <p v-if="cargando">Cargando resumen...</p>
    <p v-else-if="resumen.length === 0">No hay datos todavía.</p>

    <table v-else>
      <thead>
        <tr>
          <th>Actividad</th>
          <th>Veces hecha</th>
          <th>Veces no hecha</th>
          <th>% Cumplimiento</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="r in resumen" :key="r.id">
          <td>{{ r.titulo }}</td>
          <td>{{ r.veces_hecha }}</td>
          <td>{{ r.veces_no_hecha }}</td>
          <td>{{ r.porcentaje_cumplimiento ?? '-' }}%</td>
        </tr>
      </tbody>
    </table>
  </main>
</template>

<style scoped>
main {
  max-width: 640px;
  margin: 0 auto;
  padding: 1rem;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th,
td {
  text-align: left;
  padding: 0.5rem;
  border-bottom: 1px solid var(--color-border);
}

.error {
  color: #d33;
}
</style>
