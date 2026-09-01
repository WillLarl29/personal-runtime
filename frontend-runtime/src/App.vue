<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <header>
    <nav>
      <RouterLink to="/">Actividades</RouterLink>
      <RouterLink to="/resumen">Resumen</RouterLink>
    </nav>
  </header>

  <RouterView />
</template>

<style scoped>
header {
  padding: 1rem;
  border-bottom: 1px solid var(--color-border);
}

nav {
  display: flex;
  gap: 1rem;
  justify-content: center;
}

nav a.router-link-exact-active {
  font-weight: bold;
}
</style>
