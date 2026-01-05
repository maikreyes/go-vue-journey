<script setup lang="ts">
    const year = new Date().getFullYear()
    import { GitHub, LinkedIn } from '../config/config';
</script>

<template>
  <footer 
    class="w-full border-t
           border-gray-200 dark:border-gray-800
           bg-gray-100 dark:bg-gray-900"
  >
    <div
      class="max-w-7xl mx-auto
             px-6 py-6
             flex flex-col md:flex-row
             items-center justify-between
             gap-4"
    >
      <!-- Left -->
      <p class="text-sm text-gray-600 dark:text-indigo-100">
        © {{ year }} · Stock Tracker
      </p>

      <!-- Center -->
      <p class="text-xs text-gray-500 dark:text-indigo-200 text-center">
        Built with Vue 3 · Pinia · Tailwind · Go
      </p>

      <!-- Right -->
      <div class="flex items-center gap-4">
        <a
          :href="GitHub"
          class="text-sm text-gray-600 dark:text-indigo-100
                 hover:text-indigo-500 transition"
        >
          GitHub
        </a>
        <a
          :href="LinkedIn"
          class="text-sm text-gray-600 dark:text-indigo-100
                 hover:text-indigo-500 transition"
        >
          LinkedIn
        </a>
      </div>
    </div>
  </footer>
</template>
