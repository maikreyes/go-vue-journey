<script setup lang="ts">
import { onMounted } from 'vue'
import { useStockStore } from '../stores/stock.store'
import BestStocks from '../components/BestStocks.vue'
import StockFilters from '../components/StockFilters.vue'
import StockTable from './StockTableView.vue'


const store = useStockStore()

onMounted(() => {
  store.setStocks()
  store.setTopStocks()
})
</script>

<template>
  <section class="space-y-10">
    <BestStocks />
    <StockFilters />
    <StockTable />
  </section>
</template>
