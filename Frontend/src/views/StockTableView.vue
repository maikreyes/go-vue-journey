<script setup lang="ts">
import { onMounted } from 'vue'
import { useStockStore } from '../stores/stock.store'

import StockTable from '../components/StockTable.vue'
import StockTableSkeleton from '../components/skeletons/StockTableSkeleton.vue'

import StockCard from '../components/StockCard.vue'
import StockCardSkeleton from '../components/skeletons/StockCardSkeleton.vue'

import { pillButton } from '../ui/classes'

const store = useStockStore()

onMounted(() => {
  store.setStocks()
})
</script>

<template>
  <div
    class="flex flex-col items-center space-y-8 m-4
           dark:bg-gray-800 bg-gray-50
           shadow-lg rounded-2xl p-8"
  >

    <!-- ================= LOADING ================= -->
    <template v-if="store.loading">
      <!-- Desktop -->
      <div class="hidden md:block w-full">
        <StockTableSkeleton />
      </div>

      <!-- Mobile -->
      <div class="md:hidden w-full grid grid-cols-1 gap-6">
        <StockCardSkeleton v-for="n in 5" :key="n" />
      </div>
    </template>

    <!-- ================= DATA ================= -->
    <template v-else-if="store.hasData">

      <!-- Desktop -->
      <div class="hidden md:block w-full">
        <StockTable />
      </div>

      <!-- Mobile -->
      <div class="md:hidden w-full grid grid-cols-1 gap-6">
        <StockCard
          v-for="card in store.paginatedStoreCards"
          :key="card.ticker"
          v-bind="card"
        />
      </div>

      <!-- Pagination -->
      <div class="flex items-center justify-center gap-6">
        <button
          @click="store.prevPage"
          :disabled="store.currentPage === 1"
          :class="pillButton"
        >
          Anterior
        </button>

        <span class="text-sm text-indigo-700 dark:text-indigo-300">
          Página {{ store.currentPage }} de {{ store.totalPages }}
        </span>

        <button
          @click="store.nextPage"
          :disabled="store.currentPage === store.totalPages"
          :class="pillButton"
        >
          Siguiente
        </button>
      </div>

    </template>

    <!-- ================= EMPTY STATE (OPCIONAL) ================= -->
    <template v-else>
      <p class="text-gray-500 dark:text-gray-400">
        No hay datos disponibles
      </p>
    </template>

  </div>
</template>
