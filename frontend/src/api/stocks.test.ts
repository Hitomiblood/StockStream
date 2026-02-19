import { beforeEach, describe, expect, it, vi } from 'vitest'

import {
  filterStocks,
  getLatestStocks,
  getRecommendations,
  getStockById,
  getStocks,
  getStocksByTicker,
  searchStocks,
  syncStocks
} from '@/api/stocks'
import { client } from '@/api/http'

vi.mock('@/api/http', () => ({
  client: {
    get: vi.fn(),
    post: vi.fn()
  }
}))

describe('api/stocks', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('clamps pagination and sorting params for getStocks', async () => {
    vi.mocked(client.get).mockResolvedValue({ data: { data: [], total: 0, limit: 200, offset: 0 } })

    await getStocks({ limit: 999, offset: -5, sort: 'time', order: 'desc' })

    expect(client.get).toHaveBeenCalledWith('/stocks', {
      params: {
        limit: 200,
        offset: 0,
        sort: 'time',
        order: 'desc'
      }
    })
  })

  it('uses search endpoint with bounded limit', async () => {
    vi.mocked(client.get).mockResolvedValue({ data: { query: 'AAPL', data: [], total: 0 } })

    await searchStocks({ q: 'AAPL', limit: 0 })

    expect(client.get).toHaveBeenCalledWith('/stocks/search', {
      params: {
        q: 'AAPL',
        limit: 1
      }
    })
  })

  it('throws when filter is called without action and rating', async () => {
    await expect(filterStocks({ limit: 20, offset: 0 })).rejects.toThrow(
      'At least one filter (action or rating) is required'
    )
    expect(client.get).not.toHaveBeenCalled()
  })

  it('encodes id/ticker for detail endpoints', async () => {
    vi.mocked(client.get).mockResolvedValue({ data: {} })

    await getStockById('100/ABC')
    await getStocksByTicker('BRK/B')

    expect(client.get).toHaveBeenNthCalledWith(1, '/stocks/100%2FABC')
    expect(client.get).toHaveBeenNthCalledWith(2, '/stocks/ticker/BRK%2FB')
  })

  it('bounds limits for latest and recommendations', async () => {
    vi.mocked(client.get).mockResolvedValue({ data: { data: [], total: 0 } })

    await getLatestStocks(999)
    await getRecommendations(999)

    expect(client.get).toHaveBeenNthCalledWith(1, '/stocks/latest', {
      params: { limit: 100 }
    })
    expect(client.get).toHaveBeenNthCalledWith(2, '/recommendations', {
      params: { limit: 50 }
    })
  })

  it('executes sync with unlimited timeout', async () => {
    vi.mocked(client.post).mockResolvedValue({ data: { message: 'ok' } })

    await syncStocks()

    expect(client.post).toHaveBeenCalledWith('/stocks/fetch', undefined, {
      timeout: 0
    })
  })
})
