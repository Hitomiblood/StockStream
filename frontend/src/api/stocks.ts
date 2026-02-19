import { client } from '@/api/http'
import type {
  FilterStocksResponse,
  MetadataResponse,
  RecommendationsResponse,
  SearchStocksResponse,
  SortField,
  SortOrder,
  Stock,
  StocksListResponse,
  SyncResponse,
  TickerHistoryResponse
} from '@/types/domain'

export interface StocksQuery {
  limit: number
  offset: number
  sort: SortField
  order: SortOrder
}

export interface SearchQuery {
  q: string
  limit: number
}

export interface FilterQuery {
  action?: string
  rating?: string
  limit: number
  offset: number
}

export async function getStocks(query: StocksQuery): Promise<StocksListResponse> {
  const { data } = await client.get<StocksListResponse>('/stocks', {
    params: {
      limit: clamp(query.limit, 1, 200),
      offset: clamp(query.offset, 0),
      sort: query.sort,
      order: query.order
    }
  })

  return data
}

export async function searchStocks(query: SearchQuery): Promise<SearchStocksResponse> {
  const { data } = await client.get<SearchStocksResponse>('/stocks/search', {
    params: {
      q: query.q,
      limit: clamp(query.limit, 1, 200)
    }
  })

  return data
}

export async function filterStocks(query: FilterQuery): Promise<FilterStocksResponse> {
  const action = query.action?.trim()
  const rating = query.rating?.trim()

  if (!action && !rating) {
    throw new Error('At least one filter (action or rating) is required')
  }

  const { data } = await client.get<FilterStocksResponse>('/stocks/filter', {
    params: {
      action,
      rating,
      limit: clamp(query.limit, 1, 200),
      offset: clamp(query.offset, 0)
    }
  })

  return data
}

export async function getLatestStocks(limit = 20): Promise<{ data: Stock[]; total: number }> {
  const { data } = await client.get<{ data: Stock[]; total: number }>('/stocks/latest', {
    params: {
      limit: clamp(limit, 1, 100)
    }
  })

  return data
}

export async function getStockById(id: string): Promise<Stock> {
  const { data } = await client.get<Stock>(`/stocks/${encodeURIComponent(id)}`)
  return data
}

export async function getStocksByTicker(ticker: string): Promise<TickerHistoryResponse> {
  const { data } = await client.get<TickerHistoryResponse>(`/stocks/ticker/${encodeURIComponent(ticker)}`)
  return data
}

export async function getRecommendations(limit = 10): Promise<RecommendationsResponse> {
  const { data } = await client.get<RecommendationsResponse>('/recommendations', {
    params: {
      limit: clamp(limit, 1, 50)
    }
  })

  return data
}

export async function getMetadata(): Promise<MetadataResponse> {
  const { data } = await client.get<MetadataResponse>('/metadata')
  return data
}

export async function syncStocks(): Promise<SyncResponse> {
  const { data } = await client.post<SyncResponse>('/stocks/fetch', undefined, {
    timeout: 0
  })
  return data
}

function clamp(value: number, min: number, max?: number): number {
  if (Number.isNaN(value)) {
    return min
  }

  if (max !== undefined) {
    return Math.max(min, Math.min(max, value))
  }

  return Math.max(min, value)
}
