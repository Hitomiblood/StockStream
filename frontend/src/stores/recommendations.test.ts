import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

import { getRecommendations } from '@/api/stocks'
import { useRecommendationsStore } from '@/stores/recommendations'

vi.mock('@/api/stocks', () => ({
  getRecommendations: vi.fn()
}))

describe('stores/recommendations', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('loads recommendations payload', async () => {
    vi.mocked(getRecommendations).mockResolvedValue({
      recommendations: [
        {
          stock: {
            id: '1',
            ticker: 'AAPL',
            target_from: '100',
            target_to: '120',
            company: 'Apple',
            action: 'upgrade',
            brokerage: 'Broker',
            rating_from: 'Hold',
            rating_to: 'Buy',
            time: '2026-02-01T00:00:00Z',
            created_at: '2026-02-01T00:00:00Z',
            updated_at: '2026-02-01T00:00:00Z'
          },
          score: 95,
          reason: 'Strong upside',
          confidence: 'high'
        }
      ],
      generated_at: '2026-02-01T00:00:00Z',
      count: 1,
      criteria: {}
    })

    const store = useRecommendationsStore()
    await store.fetchRecommendations(5)

    expect(getRecommendations).toHaveBeenCalledWith(5)
    expect(store.items).toHaveLength(1)
    expect(store.count).toBe(1)
    expect(store.generatedAt).toBe('2026-02-01T00:00:00Z')
  })

  it('sets error when recommendations fail', async () => {
    vi.mocked(getRecommendations).mockRejectedValue(new Error('recommendations error'))

    const store = useRecommendationsStore()
    await store.fetchRecommendations()

    expect(store.error).toBe('recommendations error')
    expect(store.loading).toBe(false)
  })
})
