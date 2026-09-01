<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <div class="shell">
    <header class="topbar">
      <div class="topbar-inner">
        <div class="brand">
          <svg class="brand-mark" width="30" height="30" viewBox="0 0 30 30" fill="none" aria-hidden="true">
            <rect x="1" y="1" width="28" height="28" rx="9" fill="#000004" />
            <path d="M9 15.5L13 19.5L21 10.5" stroke="url(#brandGradient)" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" />
            <defs>
              <linearGradient id="brandGradient" x1="9" y1="10.5" x2="21" y2="19.5" gradientUnits="userSpaceOnUse">
                <stop stop-color="#04b5f4" />
                <stop offset="1" stop-color="#9366f4" />
              </linearGradient>
            </defs>
          </svg>
          <span class="brand-text">ACTIVIDADES</span>
        </div>

        <nav class="nav">
          <RouterLink to="/" class="nav-link">Actividades</RouterLink>
          <RouterLink to="/resumen" class="nav-link">Resumen</RouterLink>
          <RouterLink to="/categorias" class="nav-link">Categorías</RouterLink>
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
  max-width: 960px;
  margin: 0 auto;
  padding: var(--space-4) var(--space-6);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.brand {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.brand-mark {
  flex-shrink: 0;
}

.brand-text {
  font-weight: 800;
  font-size: 1rem;
  letter-spacing: 0.06em;
  color: var(--color-text);
}

.nav {
  display: flex;
  gap: var(--space-6);
}

.nav-link {
  position: relative;
  padding: 0.4rem 0.1rem;
  color: var(--color-text-muted);
  font-weight: 600;
  font-size: 0.9rem;
  transition: color 0.15s;
}

.nav-link::after {
  content: '';
  position: absolute;
  left: 0;
  right: 0;
  bottom: -1px;
  height: 2px;
  border-radius: var(--radius-full);
  background: var(--color-accent);
  transform: scaleX(0);
  transition: transform 0.18s ease;
}

.nav-link:hover {
  color: var(--color-text);
}

.nav-link.router-link-exact-active {
  color: var(--color-text);
}

.nav-link.router-link-exact-active::after {
  transform: scaleX(1);
}

@media (max-width: 600px) {
  .brand-text {
    display: none;
  }

  .nav {
    gap: var(--space-4);
  }
}
</style>
