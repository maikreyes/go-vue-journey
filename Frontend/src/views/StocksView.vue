<script setup lang="ts">
import { onMounted } from 'vue'
import { useStockStore } from '../stores/stock.store'
import BestStocks from '../components/BestStocks.vue'
import StockFilters from '../components/StockFilters.vue'
import SearchBar from '../components/StockSearchBar.vue'

import StockTableView from './StockTableView.vue'

const store = useStockStore()

onMounted(() => {
  store.setTopStocks()
})
</script>

<template>
  <section class="space-y-10">
    <BestStocks />
    <SearchBar />
    <StockFilters />
    <StockTableView />
  </section>
</template>
