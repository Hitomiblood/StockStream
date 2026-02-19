import { defineStore } from 'pinia'

import { getRecommendations } from '@/api/stocks'
import type { Recommendation } from '@/types/domain'

interface RecommendationsState {
  items: Recommendation[]
  generatedAt: string
  count: number
  loading: boolean
  error: string | null
}

export const useRecommendationsStore = defineStore('recommendations', {
  state: (): RecommendationsState => ({
    items: [],
    generatedAt: '',
    count: 0,
    loading: false,
    error: null
  }),
  actions: {
    async fetchRecommendations(limit = 10) {
      this.loading = true
      this.error = null

      try {
        const response = await getRecommendations(limit)
        this.items = response.recommendations
        this.generatedAt = response.generated_at
        this.count = response.count
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Failed to load recommendations'
      } finally {
        this.loading = false
      }
    }
  }
})
