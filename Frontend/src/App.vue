<script setup lang="ts">
  
  import StocksView from './views/StocksView.vue';
  import StockFooter from './components/StockFooter.vue';
  import AppNavbar from './components/AppNavbar.vue';
  
</script>
<template>
  <div
    class="
      min-h-screen
      bg-white text-indigo-900
      dark:bg-gray-950 dark:text-indigo-200
      transition-colors duration-300
    "
  >
    <AppNavbar />

    <main class="relative min-h-svh p-10">
      <!-- Background Pattern -->
      <div class="absolute inset-0 -z-10">
        <div
          class="
            h-full w-full
            bg-[radial-gradient(#e5e7eb_1px,transparent_1px)]
            bg-size-[16px_16px]
            dark:bg-[radial-gradient(125%_125%_at_50%_10%,#000_40%,#6366f1_100%)]
          "
        ></div>
      </div>

      <div class="relative z-10 flex h-full justify-center px-4">
        <div class="w-full max-w-7xl">
          <StocksView />
        </div>
      </div>
    </main>

    <StockFooter />
  </div>
</template>

