import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

import { getMetadata } from '@/api/stocks'
import { useMetadataStore } from '@/stores/metadata'

vi.mock('@/api/stocks', () => ({
  getMetadata: vi.fn()
}))

describe('stores/metadata', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('loads actions and ratings', async () => {
    vi.mocked(getMetadata).mockResolvedValue({
      actions: ['upgrade'],
      ratings: ['Buy']
    })

    const store = useMetadataStore()
    await store.fetchMetadata()

    expect(store.actions).toEqual(['upgrade'])
    expect(store.ratings).toEqual(['Buy'])
  })

  it('handles metadata fetch error', async () => {
    vi.mocked(getMetadata).mockRejectedValue(new Error('metadata error'))

    const store = useMetadataStore()
    await store.fetchMetadata()

    expect(store.error).toBe('metadata error')
    expect(store.loading).toBe(false)
  })
})
