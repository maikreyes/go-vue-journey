import { defineStore } from 'pinia'
import type { Stock } from '../models/stock.model'
import { stockToCard } from '../mappers/stock.mapper'
import type { CardProps } from '../models/cardProps.model'
import { fetchStocks, fetchTopStocks } from '../api/stock.api'
import { parseMoney } from '../utils/paserMoney'
import type { StocksStats } from '../models/stocksResponse.model'

export type StockFilter = 'all' | 'up' | 'down' | 'equal'
export type SortDirection = 'asc' | 'desc'

export const useStockStore = defineStore('stock', {
  state: () => ({
    stock: [] as Stock[],
    topStocks: [] as Stock[],
    ticker: '' as string,

    serverStats: null as StocksStats | null,
    serverTotalPages: 0,
    nextCursor: null as string | null,
    pageCursors: { 1: null as string | null } as Record<number, string | null>,

    filter: 'all' as StockFilter,

    sortBy: 'ticker' as keyof EnrichedStock,
    sortDirection: 'asc' as SortDirection,

    currentPage: 1,
    pageSize: 10,

    loading: false,
    error: null as string | null,
  }),

  getters: {
    hasData: (state) => state.stock.length > 0,


    enrichedStocks(state): EnrichedStock[] {
      return state.stock.map((s) => {
        const from = parseMoney(s.target_from)
        const to = parseMoney(s.target_to)

        return {
          ...s,
          priceChange: to - from,
          percentageChange: ((to - from) / ((to + from) / 2)) * 100,
        }
      })
    },


    filteredStocks(): EnrichedStock[] {
      const byFilter = this.filter === 'all'
        ? this.enrichedStocks
        : this.enrichedStocks.filter((s) => {
          if (this.filter === 'up') return s.priceChange > 0
          if (this.filter === 'down') return s.priceChange < 0
          return s.priceChange === 0
        })

      return byFilter
    },


    sortedStocks(): EnrichedStock[] {
      const dir = this.sortDirection === 'asc' ? 1 : -1

      return [...this.filteredStocks].sort((a, b) => {
        const valA = a[this.sortBy]
        const valB = b[this.sortBy]

        if (typeof valA === 'number' && typeof valB === 'number') {
          return (valA - valB) * dir
        }

        return String(valA).localeCompare(String(valB)) * dir
      })
    },


    totalPages(): number {
      if (this.serverTotalPages > 0) return this.serverTotalPages
      return Math.ceil(this.sortedStocks.length / this.pageSize)
    },

    paginatedStocks(): EnrichedStock[] {
      if (this.serverTotalPages > 0) return this.sortedStocks

      const start = (this.currentPage - 1) * this.pageSize
      return this.sortedStocks.slice(start, start + this.pageSize)
    },


    paginatedStoreCards(): CardProps[] {
      return this.paginatedStocks.map(stockToCard)
    },

    topStockCards(state): CardProps[] {
      return state.topStocks.map(stockToCard)
    },


    totalCount(): number {
      return this.serverStats?.all_stocks ?? this.stock.length
    },
    upCount(): number {
      return this.serverStats?.up_stocks ?? this.enrichedStocks.filter(s => s.priceChange > 0).length
    },
    downCount(): number {
      return this.serverStats?.down_stocks ?? this.enrichedStocks.filter(s => s.priceChange < 0).length
    },
    noChangeCount(): number {
      return this.serverStats?.no_change ?? this.enrichedStocks.filter(s => s.priceChange === 0).length
    },
  },

  actions: {
    async setStocks() {
      this.loading = true
      this.error = null

      try {
        this.currentPage = 1
        this.pageCursors = { 1: null }

        const resp = await fetchStocks({
          nextPage: null,
          filter: this.filter,
          ticker: this.ticker ? this.ticker : null,
        })
        this.stock = resp.items
        this.serverStats = resp.stats
        this.serverTotalPages = resp.stats.pages
        this.nextCursor = resp.next_cursor && resp.next_cursor.length > 0 ? resp.next_cursor : null
      } catch {
        this.error = 'Error cargando stocks'
      } finally {
        this.loading = false
      }
    },

    async setTopStocks() {
      this.loading = true
      this.error = null

      try {
        this.topStocks = await fetchTopStocks()
      } catch {
        this.error = 'Error cargando top stocks'
      } finally {
        this.loading = false
      }
    },

    async setTicker(ticker: string) {
      this.loading = true
      this.error = null
      this.ticker = ticker.toUpperCase()

      try {
        this.currentPage = 1
        this.pageCursors = { 1: null }

        const resp = await fetchStocks({
          nextPage: null,
          filter: this.filter,
          ticker: this.ticker ? this.ticker : null,
        })

        this.stock = resp.items
        this.serverStats = resp.stats
        this.serverTotalPages = resp.stats.pages
        this.nextCursor = resp.next_cursor && resp.next_cursor.length > 0 ? resp.next_cursor : null
      } catch {
        this.error = 'Error cargando stock por ticker'
      } finally {
        this.loading = false
      }
    },

    setFilter(filter: StockFilter) {
      this.filter = filter
      this.currentPage = 1

      if (this.serverTotalPages > 0) {
        void this.setStocks()
      }
    },

    setSort(column: keyof EnrichedStock) {
      if (this.sortBy === column) {
        this.sortDirection = this.sortDirection === 'asc' ? 'desc' : 'asc'
      } else {
        this.sortBy = column
        this.sortDirection = 'asc'
      }
    },

    async nextPage() {
      if (this.currentPage >= this.totalPages) return
      if (this.loading) return

      const nextPageNum = this.currentPage + 1
      const cursor = this.pageCursors[nextPageNum] ?? this.nextCursor
      if (!cursor) return

      this.loading = true
      this.error = null
      try {
        const resp = await fetchStocks({
          nextPage: cursor,
          filter: this.filter,
          ticker: this.ticker ? this.ticker : null,
        })
        this.pageCursors[nextPageNum] = cursor
        this.currentPage = nextPageNum
        this.stock = resp.items
        this.serverStats = resp.stats
        this.serverTotalPages = resp.stats.pages
        this.nextCursor = resp.next_cursor && resp.next_cursor.length > 0 ? resp.next_cursor : null
      } catch {
        this.error = 'Error cargando stocks'
      } finally {
        this.loading = false
      }
    },

    async prevPage() {
      if (this.currentPage <= 1) return
      if (this.loading) return

      const prevPageNum = this.currentPage - 1
      const cursor = this.pageCursors[prevPageNum]

      this.loading = true
      this.error = null
      try {
        const resp = await fetchStocks({
          nextPage: cursor,
          filter: this.filter,
          ticker: this.ticker ? this.ticker : null,
        })
        this.currentPage = prevPageNum
        this.stock = resp.items
        this.serverStats = resp.stats
        this.serverTotalPages = resp.stats.pages
        this.nextCursor = resp.next_cursor && resp.next_cursor.length > 0 ? resp.next_cursor : null
      } catch {
        this.error = 'Error cargando stocks'
      } finally {
        this.loading = false
      }
    },
  },
})


export interface EnrichedStock extends Stock {
  priceChange: number
  percentageChange: number
}
