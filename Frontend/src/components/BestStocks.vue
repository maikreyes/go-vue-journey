<script setup lang="ts">
    import { useStockStore } from '../stores/stock.store';
    import StockCard from './StockCard.vue'

    const props = useStockStore()

</script>
<template>
    <header class="text-center">
      <h1 class="text-4xl font-bold text-gray-900 dark:text-indigo-100">
        Top Stocks Today
      </h1>
      <p class="mt-2 text-gray-800 dark:text-indigo-200 ">
        Best opportunities based on target growth
      </p>
    </header>
    <section
        class="
            flex gap-4
            overflow-x-auto
            md:overflow-x-visible
            md:flex-wrap
            scrollbar-thin scrollbar-thumb-indigo-400
        "
    >
       <div
            v-for="stock in props.topStockCards"
            :key="stock.ticker"
            class="shrink-0"
        >
            <StockCard
                :ticker="stock.ticker"
                :targetFrom="stock.targetFrom"
                :targetTo="stock.targetTo"
                :company="stock.company"
                :action="stock.action"
            />
        </div>

    </section>
</template>