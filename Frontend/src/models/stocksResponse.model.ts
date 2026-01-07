import type { Stock } from './stock.model'

export interface StocksStats {
  total: number
  up: number
  down: number
}

export interface StocksResponse {
  items: Stock[]
  stats: StocksStats
  total_pages: number
  next_cursor?: string | null
}
