import { createRouter, createWebHistory } from 'vue-router'

import DashboardView from '@/views/DashboardView.vue'
import StockDetailView from '@/views/StockDetailView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'dashboard',
      component: DashboardView
    },
    {
      path: '/stocks/:id',
      name: 'stock-detail',
      component: StockDetailView,
      props: true
    }
  ]
})

export default router
