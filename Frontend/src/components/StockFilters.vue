<script setup lang="ts">
    import { useStockStore } from '../stores/stock.store';

    const store = useStockStore()

    const baseClass =
    'px-6 py-2 cursor-pointer select-none rounded-full font-medium transition-all duration-200 ' +
    'bg-gray-100 text-gray-900 ' +
    'dark:bg-gray-900 dark:text-gray-200 ' +
    'hover:scale-105'

    const activeClass =
    'ring-2 ring-indigo-500 dark:ring-indigo-400'

    const activeUpClass =
    'bg-green-600 text-white dark:bg-green-600 ring-2 ring-green-400'

    const activeDownClass =
    'bg-red-500 text-white dark:bg-red-600 ring-2 ring-red-400'

</script>

<template>
  <div class="flex gap-4 mb-6 justify-left m-3 ">
    
    <!-- Todas -->
    <button
      @click="store.setFilter('all')"
      :class="[
        baseClass,
        store.filter === 'all' && activeClass
      ]"
    >
      All: {{ store.totalCount }}
    </button>

    <!-- Suben -->
    <button
      @click="store.setFilter('up')"
      :class="[
        baseClass,
        store.filter === 'up' && activeUpClass
      ]"
    >
      Rise: {{ store.upCount }}
    </button>

    <!-- Bajan -->
    <button
      @click="store.setFilter('down')"
      :class="[
        baseClass,
        store.filter === 'down' && activeDownClass
      ]"
    >
      Down: {{ store.downCount }}
    </button>

  </div>
</template>
