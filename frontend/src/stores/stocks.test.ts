import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

import { useStocksStore } from '@/stores/stocks'
import { filterStocks, getStocks, searchStocks } from '@/api/stocks'

vi.mock('@/api/stocks', () => ({
  getStocks: vi.fn(),
  searchStocks: vi.fn(),
  filterStocks: vi.fn()
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

describe('stores/stocks', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('loads regular list mode', async () => {
    vi.mocked(getStocks).mockResolvedValue({
      data: [stock],
      total: 1,
      limit: 20,
      offset: 0
    })

    const store = useStocksStore()
    await store.fetchStocks()

    expect(store.mode).toBe('list')
    expect(store.items).toEqual([stock])
    expect(store.total).toBe(1)
    expect(getStocks).toHaveBeenCalledOnce()
  })

  it('loads search mode and resets offset', async () => {
    vi.mocked(searchStocks).mockResolvedValue({
      query: 'AAPL',
      data: [stock],
      total: 1
    })

    const store = useStocksStore()
    store.offset = 40
    store.setQuery('AAPL')

    await store.fetchStocks()

    expect(store.mode).toBe('search')
    expect(store.offset).toBe(0)
    expect(searchStocks).toHaveBeenCalledWith({ q: 'AAPL', limit: 20 })
  })

  it('loads filter mode with selected metadata filters', async () => {
    vi.mocked(filterStocks).mockResolvedValue({
      filters: { action: 'upgrade', rating: '' },
      data: [stock],
      total: 1,
      limit: 20,
      offset: 0
    })

    const store = useStocksStore()
    store.setActionFilter('upgrade')

    await store.fetchStocks()

    expect(store.mode).toBe('filter')
    expect(filterStocks).toHaveBeenCalledWith({
      action: 'upgrade',
      rating: undefined,
      limit: 20,
      offset: 0
    })
  })

  it('handles API failures and sets error message', async () => {
    vi.mocked(getStocks).mockRejectedValue(new Error('backend down'))

    const store = useStocksStore()
    await store.fetchStocks()

    expect(store.error).toBe('backend down')
    expect(store.loading).toBe(false)
  })

  it('navigates pages respecting boundaries', () => {
    const store = useStocksStore()
    store.total = 45
    store.limit = 20

    expect(store.canGoPrev).toBe(false)
    expect(store.canGoNext).toBe(true)

    store.nextPage()
    expect(store.offset).toBe(20)

    store.prevPage()
    expect(store.offset).toBe(0)
  })
})
