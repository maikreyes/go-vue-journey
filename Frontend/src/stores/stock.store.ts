import { defineStore } from 'pinia'
import type { Stock } from '../models/stock.model'
import { stockToCard } from '../mappers/stock.mapper'
import type { CardProps } from '../models/cardProps.model'
import { fetchStocks, fetchTopStocks } from '../api/stock.api'
import { parseMoney } from '../utils/Money'
import type { StocksStats } from '../models/stocksResponse.model'

export type StockFilter = 'all' | 'up' | 'down'
export type SortDirection = 'asc' | 'desc'

export const useStockStore = defineStore('stock', {
  state: () => ({
    stock: [] as Stock[],
    topStocks: [] as Stock[],

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

    /* 🔹 Datos enriquecidos */
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

    /* 🔹 Filtro */
    filteredStocks(): EnrichedStock[] {
      if (this.filter === 'all') return this.enrichedStocks

      return this.enrichedStocks.filter((s) =>
        this.filter === 'up'
          ? s.priceChange > 0
          : s.priceChange < 0
      )
    },

    /* 🔹 Ordenamiento */
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

    /* 🔹 Paginación */
    totalPages(): number {
      if (this.serverTotalPages > 0) return this.serverTotalPages
      return Math.ceil(this.sortedStocks.length / this.pageSize)
    },

    paginatedStocks(): EnrichedStock[] {
      // Con paginación por cursor, el backend ya entrega “una página” (limit).
      // Mantenemos el mismo contrato con los componentes.
      if (this.serverTotalPages > 0) return this.sortedStocks

      const start = (this.currentPage - 1) * this.pageSize
      return this.sortedStocks.slice(start, start + this.pageSize)
    },

    /* 🔹 Cards */
    paginatedStoreCards(): CardProps[] {
      return this.paginatedStocks.map(stockToCard)
    },

    topStockCards(state): CardProps[] {
      return state.topStocks.map(stockToCard)
    },

    /* 🔹 Contadores */
    totalCount(): number {
      return this.serverStats?.total ?? this.stock.length
    },
    upCount(): number {
      return this.serverStats?.up ?? this.enrichedStocks.filter(s => s.priceChange > 0).length
    },
    downCount(): number {
      return this.serverStats?.down ?? this.enrichedStocks.filter(s => s.priceChange < 0).length
    },
  },

  actions: {
    async setStocks() {
      this.loading = true
      this.error = null

      try {
        this.currentPage = 1
        this.pageCursors = { 1: null }

        const resp = await fetchStocks(this.pageSize, null, this.filter)
        this.stock = resp.items
        this.serverStats = resp.stats
        this.serverTotalPages = resp.total_pages
        this.nextCursor = resp.next_cursor ?? null
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

    setFilter(filter: StockFilter) {
      this.filter = filter
      this.currentPage = 1

    // Cuando estamos usando paginación/stats del backend,
    // el filtro debe aplicarse server-side.
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
        const resp = await fetchStocks(this.pageSize, cursor, this.filter)
        this.pageCursors[nextPageNum] = cursor
        this.currentPage = nextPageNum
        this.stock = resp.items
        this.serverStats = resp.stats
        this.serverTotalPages = resp.total_pages
        this.nextCursor = resp.next_cursor ?? null
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
        const resp = await fetchStocks(this.pageSize, cursor, this.filter)
        this.currentPage = prevPageNum
        this.stock = resp.items
        this.serverStats = resp.stats
        this.serverTotalPages = resp.total_pages
        this.nextCursor = resp.next_cursor ?? null
      } catch {
        this.error = 'Error cargando stocks'
      } finally {
        this.loading = false
      }
    },
  },
})

/* 🔹 Tipo enriquecido */
export interface EnrichedStock extends Stock {
  priceChange: number
  percentageChange: number
}
