import type { Stock } from './stock.model'

export interface StocksStats {
  all_stocks: number
  up_stocks: number
  down_stocks: number
  no_change: number
  pages: number
}

export interface StocksResponse {
  items: Stock[]
  stats: StocksStats
  next_cursor?: string | null
}
