export type SortOrder = 'asc' | 'desc'

export type SortField =
  | 'id'
  | 'ticker'
  | 'target_from'
  | 'target_to'
  | 'company'
  | 'action'
  | 'brokerage'
  | 'rating_from'
  | 'rating_to'
  | 'time'
  | 'created_at'
  | 'updated_at'

export interface Stock {
  id: string
  ticker: string
  target_from: string
  target_to: string
  company: string
  action: string
  brokerage: string
  rating_from: string
  rating_to: string
  time: string
  created_at: string
  updated_at: string
}

export interface StocksListResponse {
  data: Stock[]
  total: number
  limit: number
  offset: number
}

export interface SearchStocksResponse {
  query: string
  data: Stock[]
  total: number
}

export interface FilterStocksResponse {
  filters: {
    action: string
    rating: string
  }
  data: Stock[]
  total: number
  limit: number
  offset: number
}

export interface TickerHistoryResponse {
  ticker: string
  company: string
  history: Stock[]
  total: number
}

export interface MetadataResponse {
  actions: string[]
  ratings: string[]
}

export interface Recommendation {
  stock: Stock
  score: number
  reason: string
  confidence: 'low' | 'medium' | 'high'
}

export interface RecommendationsResponse {
  recommendations: Recommendation[]
  generated_at: string
  count: number
  criteria: Record<string, number>
}

export interface SyncResponse {
  message: string
  total_new: number
  total_updated: number
  total_fetched: number
  duration_ms: number
}
