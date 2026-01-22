<script setup lang="ts">
import { useStockStore } from '../stores/stock.store'
import { parseMoney } from '../utils/paserMoney'

const store = useStockStore()

const sortIcon = (col: string) => {
  if (store.sortBy !== col) return ''
  return store.sortDirection === 'asc' ? '▲' : '▼'
}

const priceDiff = (stock: { target_from?: string | number | null; target_to?: string | number | null }) => {
  const from = parseMoney(stock.target_from)
  const to = parseMoney(stock.target_to)
  return to - from
}

const priceClass = (stock: { target_from?: string | number | null; target_to?: string | number | null }) => {
  const diff = priceDiff(stock)
  if (diff > 0) return 'text-green-600'
  if (diff < 0) return 'text-red-500'
  return 'text-yellow-500'
}

const percentChange = (stock: { target_from?: string | number | null; target_to?: string | number | null }) => {
  const from = parseMoney(stock.target_from)
  const to = parseMoney(stock.target_to)
  const avg = (to + from) / 2
  if (avg === 0) return 0
  return ((to - from) / avg) * 100
}
</script>

<template>
  <!-- DESKTOP -->
  <div class="hidden md:block w-full overflow-x-auto">
    <table class="w-full border-collapse text-gray-900 dark:text-indigo-200">
      <thead>
        <tr class="border-b dark:border-gray-700">

          <!-- Ticker -->
          <th
            @click="store.setSort('ticker')"
            class="
              py-3 px-2 text-left
              cursor-pointer select-none
              transition-all duration-200
              hover:bg-indigo-100 dark:hover:bg-gray-700
              hover:shadow-sm
              group
            "
            :class="store.sortBy === 'ticker'
              ? 'bg-indigo-100 dark:bg-gray-700'
              : ''"
          >
            <span class="inline-flex items-center gap-1">
              Ticker
              <span class="text-xs opacity-0 group-hover:opacity-100 transition">
                {{ sortIcon('ticker') || '⇅' }}
              </span>
            </span>
          </th>

          <!-- Company -->
          <th
            @click="store.setSort('company')"
            class="
              py-3 px-2 text-left
              cursor-pointer select-none
              transition-all duration-200
              hover:bg-indigo-100 dark:hover:bg-gray-700
              hover:shadow-sm
              group
            "
            :class="store.sortBy === 'company'
              ? 'bg-indigo-100 dark:bg-gray-700'
              : ''"
          >
            <span class="inline-flex items-center gap-1">
              Company
              <span class="text-xs opacity-0 group-hover:opacity-100 transition">
                {{ sortIcon('company') || '⇅' }}
              </span>
            </span>
          </th>

          <!-- Current Price -->
          <th
            @click="store.setSort('target_from')"
            class="
              py-3 px-2 text-right
              cursor-pointer select-none
              transition-all duration-200
              hover:bg-indigo-100 dark:hover:bg-gray-700
              hover:shadow-sm
              group
            "
            :class="store.sortBy === 'target_from'
              ? 'bg-indigo-100 dark:bg-gray-700'
              : ''"
          >
            <span class="inline-flex items-center gap-1 justify-end w-full">
              Current Price
              <span class="text-xs opacity-0 group-hover:opacity-100 transition">
                {{ sortIcon('target_from') || '⇅' }}
              </span>
            </span>
          </th>

          <!-- Next Price -->
          <th
            @click="store.setSort('target_to')"
            class="
              py-3 px-2 text-right
              cursor-pointer select-none
              transition-all duration-200
              hover:bg-indigo-100 dark:hover:bg-gray-700
              hover:shadow-sm
              group
            "
            :class="store.sortBy === 'target_to'
              ? 'bg-indigo-100 dark:bg-gray-700'
              : ''"
          >
            <span class="inline-flex items-center gap-1 justify-end w-full">
              Next Price
              <span class="text-xs opacity-0 group-hover:opacity-100 transition">
                {{ sortIcon('target_to') || '⇅' }}
              </span>
            </span>
          </th>

          <!-- Best Action -->
          <th
            @click="store.setSort('rating_to')"
            class="
              py-3 px-2 text-right
              cursor-pointer select-none
              transition-all duration-200
              hover:bg-indigo-100 dark:hover:bg-gray-700
              hover:shadow-sm
              group
            "
            :class="store.sortBy === 'rating_to'
              ? 'bg-indigo-100 dark:bg-gray-700'
              : ''"
          >
            <span class="inline-flex items-center gap-1 justify-end w-full">
              Best Action
              <span class="text-xs opacity-0 group-hover:opacity-100 transition">
                {{ sortIcon('rating_to') || '⇅' }}
              </span>
            </span>
          </th>

          <!-- Price Change -->
          <th
            @click="store.setSort('priceChange')"
            class="
              py-3 px-2 text-right
              cursor-pointer select-none
              transition-all duration-200
              hover:bg-indigo-100 dark:hover:bg-gray-700
              hover:shadow-sm
              group
            "
            :class="store.sortBy === 'priceChange'
              ? 'bg-indigo-100 dark:bg-gray-700'
              : ''"
          >
            <span class="inline-flex items-center gap-1 justify-end w-full">
              Price Change
              <span class="text-xs opacity-0 group-hover:opacity-100 transition">
                {{ sortIcon('priceChange') || '⇅' }}
              </span>
            </span>
          </th>

          <!-- % Change -->
          <th
            @click="store.setSort('percentageChange')"
            class="
              py-3 px-2 text-right
              cursor-pointer select-none
              transition-all duration-200
              hover:bg-indigo-100 dark:hover:bg-gray-700
              hover:shadow-sm
              group
            "
            :class="store.sortBy === 'percentageChange'
              ? 'bg-indigo-100 dark:bg-gray-700'
              : ''"
          >
            <span class="inline-flex items-center gap-1 justify-end w-full">
              % Change
              <span class="text-xs opacity-0 group-hover:opacity-100 transition">
                {{ sortIcon('percentageChange') || '⇅' }}
              </span>
            </span>
          </th>

        </tr>
      </thead>

      <tbody>
        <tr
          v-for="stock in store.paginatedStocks"
          :key="stock.ticker"
          class="border-b dark:border-gray-800 hover:bg-indigo-50 dark:hover:bg-gray-800"
        >
          <td class="py-3 px-2">{{ stock.ticker }}</td>
          <td class="py-3 px-2">{{ stock.company }}</td>
          <td class="py-3 px-2 text-right">{{ stock.target_from }}</td>
          <td class="py-3 px-2 text-right">{{ stock.target_to }}</td>
          <td class="py-3 px-2 text-right">{{ stock.rating_to }}</td>

          <td
            class="py-3 px-2 text-right font-semibold"
            :class="priceClass(stock)"
          >
            {{
              `$${priceDiff(stock).toFixed(2)}`
            }}
          </td>

          <td
            class="py-3 px-2 text-right font-semibold"
            :class="priceClass(stock)"
          >
            {{
              `${percentChange(stock).toFixed(2)}%`
            }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
