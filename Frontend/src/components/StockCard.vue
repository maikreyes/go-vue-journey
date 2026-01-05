<script setup lang="ts">
import { computed } from 'vue'
import TrendingUp from '../../public/trending_up.svg?url'
import TrendingDown from '../../public/trending_down.svg?url'
import type { CardProps } from '../models/cardProps.model';
import { parseMoney } from '../utils/Money';

const props = defineProps<CardProps>()

const valueDifference = computed(() =>
  parseMoney(props.targetTo) - parseMoney(props.targetFrom)
)

const valueString = computed(() => {
  const diff = valueDifference.value
  return diff >= 0
    ? `+$${diff.toFixed(2)}`
    : `-$${Math.abs(diff).toFixed(2)}`
})

const trendingIcon = computed(() =>
  valueDifference.value < 0 ? TrendingDown : TrendingUp
)

const actionColor = computed(() => {
  if (props.action === 'Buy') return 'text-green-500'
  return 'text-gray-600'
})
</script>

<template>
  <section
   class="bg-gray-50 dark:bg-gray-900 shadow-lg rounded-2xl
         w-50 max-w-sm h-40 m-4
         flex flex-col justify-between
         hover:scale-[1.02] transition-transform"
  >
    <section class="flex justify-between px-4 py-3 h-22">
      <div class="text-left">
        <div class="flex items-center gap-2">
          <h1 class="font-bold text-lg truncate dark:text-indigo-100">{{ ticker }}</h1>
          <img :src="trendingIcon" class="w-5 h-5" alt="trending" />
        </div>
        <h2 class="text-gray-600 dark:text-indigo-100 text-sm ">{{ company }}</h2>
      </div>

      <div class="font-semibold text-lg dark:text-indigo-100">
        {{ targetTo }}
      </div>
    </section>

    <section class="flex justify-between px-4 py-4">
      <div class="text-left">
        <h4 class="text-xs text-gray-500 dark:text-indigo-50">Balance</h4>
        <h3 :class="valueDifference >= 0 ? 'text-green-500' : 'text-red-500'">
          {{ valueString }}
        </h3>
      </div>

      <div class="text-right">
        <h4 class="text-xs text-gray-500 dark:text-indigo-50">Action</h4>
        <h3 :class="actionColor">{{ action }}</h3>
      </div>
    </section>
  </section>
</template>
