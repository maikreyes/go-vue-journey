<script setup lang="ts">
    import { GitHub, LinkedIn } from '../config/config';
</script>

<template>
  <nav
    class="w-full flex items-center justify-between px-6 py-4
           bg-gray-100 dark:bg-gray-900
           text-gray-900 dark:text-indigo-50 shadow"
  >
    <!-- Title -->
    <h1 class="text-xl font-bold">Stock Dashboard</h1>

    <!-- Social links -->
    <div class="flex items-center gap-3">
      <!-- GitHub -->
      <a
        :href= "GitHub"
        target="_blank"
        rel="noopener noreferrer"
        class="
          flex items-center gap-2 px-3 py-2 rounded
          bg-indigo-100 hover:bg-indigo-200
          dark:bg-gray-700 dark:hover:bg-gray-600
          transition
        "
      >
        <svg
          class="w-5 h-5"
          fill="currentColor"
          viewBox="0 0 24 24"
          aria-hidden="true"
        >
          <path
            fill-rule="evenodd"
            clip-rule="evenodd"
            d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.009-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.004.07 1.532 1.032 1.532 1.032.892 1.53 2.341 1.088 2.91.832.091-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.03-2.688-.103-.253-.447-1.27.098-2.645 0 0 .84-.269 2.75 1.026A9.564 9.564 0 0112 6.844a9.55 9.55 0 012.504.337c1.909-1.295 2.747-1.026 2.747-1.026.546 1.375.202 2.392.1 2.645.64.7 1.028 1.595 1.028 2.688 0 3.848-2.338 4.695-4.566 4.943.36.31.68.921.68 1.856 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.523 2 12 2z"
          />
        </svg>
        <span class="hidden sm:inline">GitHub</span>
      </a>

      <!-- LinkedIn -->
      <a
        :href="LinkedIn"
        target="_blank"
        rel="noopener noreferrer"
        class="
          flex items-center gap-2 px-3 py-2 rounded
          bg-indigo-100 hover:bg-indigo-200
          dark:bg-gray-700 dark:hover:bg-gray-600
          transition
        "
      >
        <svg
          class="w-5 h-5"
          fill="currentColor"
          viewBox="0 0 24 24"
          aria-hidden="true"
        >
          <path
            d="M19 0h-14c-2.761 0-5 2.239-5 5v14c0 
               2.761 2.239 5 5 5h14c2.761 
               0 5-2.239 5-5v-14c0-2.761-2.239-5-5-5zm-11 
               19h-3v-11h3v11zm-1.5-12.268c-.966 
               0-1.75-.79-1.75-1.764s.784-1.764 
               1.75-1.764 1.75.79 1.75 
               1.764-.784 1.764-1.75 1.764zm13.5 
               12.268h-3v-5.604c0-1.337-.026-3.059-1.865-3.059-1.867 
               0-2.154 1.459-2.154 2.967v5.696h-3v-11h2.881v1.5h.041c.401-.762 
               1.381-1.566 2.843-1.566 3.039 0 3.6 2.002 3.6 4.604v6.462z"
          />
        </svg>
        <span class="hidden sm:inline">LinkedIn</span>
      </a>
    </div>
  </nav>
</template>
