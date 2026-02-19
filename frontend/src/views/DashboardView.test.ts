import { createPinia, setActivePinia } from 'pinia'
import { flushPromises, mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { filterStocks, getMetadata, getRecommendations, getStocks, searchStocks } from '@/api/stocks'
import DashboardView from '@/views/DashboardView.vue'

vi.mock('@/api/stocks', () => ({
  getStocks: vi.fn(),
  searchStocks: vi.fn(),
  filterStocks: vi.fn(),
  getMetadata: vi.fn(),
  getRecommendations: vi.fn(),
  syncStocks: vi.fn()
}))

const stock = {
  id: '1',
  ticker: 'AAPL',
  target_from: '100',
  target_to: '110',
  company: 'Apple',
  action: 'upgrade',
  brokerage: 'Broker',
  rating_from: 'Hold',
  rating_to: 'Buy',
  time: '2026-02-01T00:00:00Z',
  created_at: '2026-02-01T00:00:00Z',
  updated_at: '2026-02-01T00:00:00Z'
}

const router = createRouter({
  history: createWebHistory(),
  routes: [{ path: '/', component: DashboardView }]
})

describe('views/DashboardView', () => {
  beforeEach(async () => {
    setActivePinia(createPinia())
    vi.clearAllMocks()

    vi.mocked(getStocks).mockResolvedValue({
      data: [stock],
      total: 1,
      limit: 20,
      offset: 0
    })
    vi.mocked(searchStocks).mockResolvedValue({ query: 'AAPL', data: [stock], total: 1 })
    vi.mocked(filterStocks).mockResolvedValue({
      filters: { action: 'upgrade', rating: 'Buy' },
      data: [stock],
      total: 1,
      limit: 20,
      offset: 0
    })
    vi.mocked(getMetadata).mockResolvedValue({
      actions: ['upgrade'],
      ratings: ['Buy']
    })
    vi.mocked(getRecommendations).mockResolvedValue({
      recommendations: [],
      generated_at: '2026-02-01T00:00:00Z',
      count: 0,
      criteria: {}
    })

    if (!router.hasRoute('dashboard-test')) {
      router.addRoute({ path: '/dashboard-test', name: 'dashboard-test', component: DashboardView })
    }
    await router.push('/dashboard-test')
    await router.isReady()
  })

  it('loads stocks, metadata and recommendations on mount', async () => {
    mount(DashboardView, {
      global: {
        plugins: [createPinia(), router],
        stubs: {
          RecommendationsPanel: true,
          SettingsSyncSection: true
        }
      }
    })

    await flushPromises()

    expect(getStocks).toHaveBeenCalled()
    expect(getMetadata).toHaveBeenCalled()
    expect(getRecommendations).toHaveBeenCalledWith(10)
  })

  it('triggers search endpoint when user searches', async () => {
    const wrapper = mount(DashboardView, {
      global: {
        plugins: [createPinia(), router],
        stubs: {
          RecommendationsPanel: true,
          SettingsSyncSection: true
        }
      }
    })

    await flushPromises()

    const searchInput = wrapper.find('input[placeholder="Ticker or company"]')
    await searchInput.setValue('AAPL')
    await wrapper.find('button').trigger('click')
    await flushPromises()

    expect(searchStocks).toHaveBeenCalledWith({ q: 'AAPL', limit: 20 })
  })
})
