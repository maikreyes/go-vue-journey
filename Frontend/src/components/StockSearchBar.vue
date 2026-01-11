<script setup lang="ts">
    import { ref, watch } from 'vue';
    import { useStockStore } from '../stores/stock.store';  
    const store = useStockStore();
    const searchQuery = ref('');
    watch(searchQuery, (newQuery) => {
        store.setTicker(newQuery);
    });

</script>
<template>
  <div class="search-bar">
    <input
      v-model="searchQuery"
      type="text"
      placeholder="Search by ticker..."
      class="search-input rounded-2xl bg-indigo-50 text-gray-900 px-4 py-2 focus:border-indigo-400 focus:outline-none focus:ring-2 focus:ring-indigo-200 dark:bg-gray-800 dark:text-indigo-200 dark:focus:border-gray-600 dark:focus:ring-gray-700 w-full transition"
    />
  </div>
</template>