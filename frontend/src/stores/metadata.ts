import { defineStore } from 'pinia'

import { getMetadata } from '@/api/stocks'

interface MetadataState {
  actions: string[]
  ratings: string[]
  loading: boolean
  error: string | null
}

export const useMetadataStore = defineStore('metadata', {
  state: (): MetadataState => ({
    actions: [],
    ratings: [],
    loading: false,
    error: null
  }),
  actions: {
    async fetchMetadata() {
      this.loading = true
      this.error = null

      try {
        const response = await getMetadata()
        this.actions = response.actions
        this.ratings = response.ratings
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Failed to load metadata'
      } finally {
        this.loading = false
      }
    }
  }
})
