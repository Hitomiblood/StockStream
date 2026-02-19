import { defineStore } from 'pinia'

import { filterStocks, getStocks, searchStocks } from '@/api/stocks'
import type { SortField, SortOrder, Stock } from '@/types/domain'

export const SORT_FIELDS: SortField[] = [
  'time',
  'ticker',
  'company',
  'action',
  'brokerage',
  'rating_to',
  'target_to',
  'created_at',
  'updated_at',
  'id',
  'target_from',
  'rating_from'
]

type ViewMode = 'list' | 'search' | 'filter'

interface StocksState {
  items: Stock[]
  total: number
  limit: number
  offset: number
  sort: SortField
  order: SortOrder
  query: string
  actionFilter: string
  ratingFilter: string
  loading: boolean
  error: string | null
  mode: ViewMode
}

export const useStocksStore = defineStore('stocks', {
  state: (): StocksState => ({
    items: [],
    total: 0,
    limit: 20,
    offset: 0,
    sort: 'time',
    order: 'desc',
    query: '',
    actionFilter: '',
    ratingFilter: '',
    loading: false,
    error: null,
    mode: 'list'
  }),
  getters: {
    hasActiveFilters: (state) => Boolean(state.actionFilter || state.ratingFilter),
    hasSearch: (state) => Boolean(state.query.trim()),
    canGoPrev: (state) => state.offset > 0 && state.mode !== 'search',
    canGoNext: (state) => state.offset + state.limit < state.total && state.mode !== 'search',
    currentPage: (state) => Math.floor(state.offset / state.limit) + 1,
    totalPages: (state) => Math.max(1, Math.ceil(state.total / state.limit))
  },
  actions: {
    setQuery(query: string) {
      this.query = query
      this.offset = 0
    },
    setSort(sort: SortField) {
      this.sort = sort
      this.offset = 0
    },
    setOrder(order: SortOrder) {
      this.order = order
      this.offset = 0
    },
    setActionFilter(action: string) {
      this.actionFilter = action
      this.offset = 0
    },
    setRatingFilter(rating: string) {
      this.ratingFilter = rating
      this.offset = 0
    },
    clearFilters() {
      this.actionFilter = ''
      this.ratingFilter = ''
      this.offset = 0
    },
    nextPage() {
      if (this.canGoNext) {
        this.offset += this.limit
      }
    },
    prevPage() {
      if (this.canGoPrev) {
        this.offset = Math.max(0, this.offset - this.limit)
      }
    },
    async fetchStocks() {
      this.loading = true
      this.error = null

      try {
        if (this.query.trim()) {
          this.mode = 'search'
          const response = await searchStocks({ q: this.query.trim(), limit: this.limit })
          this.items = response.data
          this.total = response.total
          this.offset = 0
          return
        }

        if (this.actionFilter || this.ratingFilter) {
          this.mode = 'filter'
          const response = await filterStocks({
            action: this.actionFilter || undefined,
            rating: this.ratingFilter || undefined,
            limit: this.limit,
            offset: this.offset
          })
          this.items = response.data
          this.total = response.total
          return
        }

        this.mode = 'list'
        const response = await getStocks({
          limit: this.limit,
          offset: this.offset,
          sort: this.sort,
          order: this.order
        })
        this.items = response.data
        this.total = response.total
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Failed to load stocks'
      } finally {
        this.loading = false
      }
    }
  }
})
