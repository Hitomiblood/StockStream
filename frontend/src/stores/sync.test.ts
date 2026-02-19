import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

import { syncStocks } from '@/api/stocks'
import { useSyncStore } from '@/stores/sync'

vi.mock('@/api/stocks', () => ({
  syncStocks: vi.fn()
}))

describe('stores/sync', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('stores sync result on success', async () => {
    vi.mocked(syncStocks).mockResolvedValue({
      message: 'done',
      total_new: 10,
      total_updated: 5,
      total_fetched: 15,
      duration_ms: 1234
    })

    const store = useSyncStore()
    await store.runSync()

    expect(store.result?.message).toBe('done')
    expect(store.error).toBeNull()
    expect(store.loading).toBe(false)
  })

  it('captures sync error', async () => {
    vi.mocked(syncStocks).mockRejectedValue(new Error('sync failed'))

    const store = useSyncStore()
    await store.runSync()

    expect(store.error).toBe('sync failed')
    expect(store.loading).toBe(false)
  })
})
