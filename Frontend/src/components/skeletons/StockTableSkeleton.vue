<template>
  <div class="overflow-x-auto">
    <table class="w-full border-collapse animate-pulse">
      <thead>
        <tr>
          <th
            v-for="n in 6"
            :key="n"
            class="px-4 py-3"
          >
            <div class="h-4 w-20 bg-gray-300 dark:bg-gray-700 rounded"></div>
          </th>
        </tr>
      </thead>

      <tbody>
        <tr v-for="row in 10" :key="row" class="border-t">
          <td v-for="col in 6" :key="col" class="px-4 py-3">
            <div class="h-4 w-full bg-gray-200 dark:bg-gray-800 rounded"></div>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
