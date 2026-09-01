import { createRouter, createWebHistory } from 'vue-router'
import ActividadesView from '../views/ActividadesView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'actividades',
      component: ActividadesView,
    },
    {
      path: '/resumen',
      name: 'resumen',
      component: () => import('../views/ResumenView.vue'),
    },
    {
      path: '/categorias',
      name: 'categorias',
      component: () => import('../views/CategoriasView.vue'),
    },
  ],
})

export default router
