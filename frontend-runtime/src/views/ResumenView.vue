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
  <main class="page">
    <h1>Resumen</h1>

    <p v-if="error" class="banner banner-error">⚠️ {{ error }}</p>
    <p v-if="cargando" class="empty-state">Cargando resumen…</p>

    <div v-else-if="resumen.length === 0" class="empty-state">
      <p>Todavía no hay actividades para resumir.</p>
    </div>

    <ul v-else class="lista">
      <li v-for="r in resumen" :key="r.id" class="card item">
        <div class="item-head">
          <strong>{{ r.titulo }}</strong>
          <span class="pct">{{ r.porcentaje_cumplimiento ?? 0 }}%</span>
        </div>

        <div class="progress-track">
          <div class="progress-fill" :style="{ width: `${r.porcentaje_cumplimiento ?? 0}%` }" />
        </div>

        <div class="item-stats">
          <span>✅ {{ r.veces_hecha }} hechas</span>
          <span>⬜ {{ r.veces_no_hecha }} pendientes</span>
        </div>
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

.lista {
  list-style: none;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.item {
  padding: var(--space-4) var(--space-5);
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.item-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
}

.pct {
  font-weight: 800;
  font-size: 1.1rem;
  background: var(--gradient-accent);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
}

.progress-track {
  height: 8px;
  border-radius: var(--radius-full);
  background: var(--color-border);
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: var(--gradient-accent);
  border-radius: var(--radius-full);
  transition: width 0.3s ease;
}

.item-stats {
  display: flex;
  gap: var(--space-4);
  font-size: 0.85rem;
  color: var(--color-text-muted);
}
</style>
