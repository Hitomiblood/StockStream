import { defineStore } from 'pinia'

import { syncStocks } from '@/api/stocks'
import type { SyncResponse } from '@/types/domain'

interface SyncState {
  loading: boolean
  error: string | null
  result: SyncResponse | null
}

export const useSyncStore = defineStore('sync', {
  state: (): SyncState => ({
    loading: false,
    error: null,
    result: null
  }),
  actions: {
    async runSync() {
      this.loading = true
      this.error = null

      try {
        this.result = await syncStocks()
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Failed to sync stocks'
      } finally {
        this.loading = false
      }
    }
  }
})
