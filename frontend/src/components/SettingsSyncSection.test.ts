import { createPinia, setActivePinia } from 'pinia'
import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { syncStocks } from '@/api/stocks'
import { useSyncStore } from '@/stores/sync'
import SettingsSyncSection from '@/components/SettingsSyncSection.vue'

vi.mock('@/api/stocks', () => ({
  syncStocks: vi.fn()
}))

describe('components/SettingsSyncSection', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.restoreAllMocks()
  })

  it('does not run sync when user cancels confirmation', async () => {
    vi.spyOn(window, 'confirm').mockReturnValue(false)

    const wrapper = mount(SettingsSyncSection, {
      global: {
        plugins: [createPinia()]
      }
    })

    await wrapper.find('button').trigger('click')

    expect(syncStocks).not.toHaveBeenCalled()
  })

  it('runs sync and shows result when confirmed', async () => {
    vi.spyOn(window, 'confirm').mockReturnValue(true)
    vi.mocked(syncStocks).mockResolvedValue({
      message: 'Data fetched successfully',
      total_new: 4,
      total_updated: 2,
      total_fetched: 6,
      duration_ms: 999
    })

    const wrapper = mount(SettingsSyncSection, {
      global: {
        plugins: [createPinia()]
      }
    })

    await wrapper.find('button').trigger('click')
    await flushPromises()

    const store = useSyncStore()
    expect(store.result?.total_fetched).toBe(6)
    expect(wrapper.text()).toContain('Data fetched successfully')
  })
})
