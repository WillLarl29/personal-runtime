<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <div class="shell">
    <header class="topbar">
      <div class="topbar-inner">
        <span class="brand">✅ Actividades</span>
        <nav class="nav">
          <RouterLink to="/" class="nav-link">Actividades</RouterLink>
          <RouterLink to="/resumen" class="nav-link">Resumen</RouterLink>
        </nav>
      </div>
    </header>

    <RouterView />
  </div>
</template>

<style scoped>
.shell {
  min-height: 100vh;
}

.topbar {
  position: sticky;
  top: 0;
  z-index: 10;
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border);
}

.topbar-inner {
  max-width: 720px;
  margin: 0 auto;
  padding: var(--space-4) var(--space-6);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.brand {
  font-weight: 700;
  font-size: 1.05rem;
}

.nav {
  display: flex;
  gap: var(--space-2);
}

.nav-link {
  padding: 0.4rem 0.85rem;
  border-radius: var(--radius-sm);
  color: var(--color-text-muted);
  font-weight: 500;
  font-size: 0.9rem;
}

.nav-link:hover {
  background: var(--color-surface-hover);
  color: var(--color-text);
}

.nav-link.router-link-exact-active {
  background: var(--color-accent-soft);
  color: var(--color-accent);
}
</style>
