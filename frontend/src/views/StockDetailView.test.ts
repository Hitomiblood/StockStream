import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { getStockById, getStocksByTicker } from '@/api/stocks'
import StockDetailView from '@/views/StockDetailView.vue'

vi.mock('@/api/stocks', () => ({
  getStockById: vi.fn(),
  getStocksByTicker: vi.fn()
}))

vi.mock('vue-router', () => ({
  useRoute: () => ({
    params: {
      id: '1'
    }
  }),
  RouterLink: {
    template: '<a><slot /></a>'
  }
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

describe('views/StockDetailView', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('loads and renders stock detail and history', async () => {
    vi.mocked(getStockById).mockResolvedValue(stock)
    vi.mocked(getStocksByTicker).mockResolvedValue({
      ticker: 'AAPL',
      company: 'Apple',
      history: [stock],
      total: 1
    })

    const wrapper = mount(StockDetailView)
    await flushPromises()

    expect(wrapper.text()).toContain('AAPL · Apple')
    expect(wrapper.text()).toContain('Ticker history')
    expect(getStockById).toHaveBeenCalledWith('1')
    expect(getStocksByTicker).toHaveBeenCalledWith('AAPL')
  })

  it('renders top-level error when stock fetch fails', async () => {
    vi.mocked(getStockById).mockRejectedValue(new Error('detail failed'))

    const wrapper = mount(StockDetailView)
    await flushPromises()

    expect(wrapper.text()).toContain('detail failed')
  })

  it('renders history error while keeping stock detail', async () => {
    vi.mocked(getStockById).mockResolvedValue(stock)
    vi.mocked(getStocksByTicker).mockRejectedValue(new Error('history failed'))

    const wrapper = mount(StockDetailView)
    await flushPromises()

    expect(wrapper.text()).toContain('AAPL · Apple')
    expect(wrapper.text()).toContain('history failed')
  })
})
