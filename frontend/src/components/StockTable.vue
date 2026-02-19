<script setup lang="ts">
import type { Stock } from '@/types/domain'

defineProps<{
  items: Stock[]
  loading: boolean
  error: string | null
}>()

const emit = defineEmits<{
  rowClick: [stock: Stock]
}>()

function onRowClick(stock: Stock): void {
  emit('rowClick', stock)
}
</script>

<template>
  <section class="rounded-lg border border-slate-200 bg-white">
    <header class="border-b border-slate-200 px-4 py-3">
      <h2 class="text-sm font-semibold text-slate-900">Stocks</h2>
    </header>

    <div v-if="loading" class="space-y-2 p-4">
      <div v-for="idx in 6" :key="idx" class="h-10 animate-pulse rounded bg-slate-100"></div>
    </div>

    <div v-else-if="error" class="p-4 text-sm text-red-600">
      {{ error }}
    </div>

    <div v-else-if="items.length === 0" class="p-4 text-sm text-slate-500">
      No stocks found for the current criteria.
    </div>

    <div v-else class="overflow-x-auto">
      <table class="min-w-full divide-y divide-slate-200 text-sm">
        <thead class="bg-slate-50 text-left text-xs uppercase tracking-wide text-slate-500">
          <tr>
            <th class="px-3 py-2">Ticker</th>
            <th class="px-3 py-2">Company</th>
            <th class="px-3 py-2">Action</th>
            <th class="px-3 py-2">Brokerage</th>
            <th class="px-3 py-2">Rating</th>
            <th class="px-3 py-2">Target</th>
            <th class="px-3 py-2">Time</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100 bg-white">
          <tr
            v-for="stock in items"
            :key="stock.id"
            class="cursor-pointer hover:bg-slate-50"
            @click="onRowClick(stock)"
          >
            <td class="px-3 py-2 font-medium text-slate-900">{{ stock.ticker }}</td>
            <td class="px-3 py-2 text-slate-700">{{ stock.company }}</td>
            <td class="px-3 py-2 text-slate-700">{{ stock.action }}</td>
            <td class="px-3 py-2 text-slate-700">{{ stock.brokerage }}</td>
            <td class="px-3 py-2 text-slate-700">{{ stock.rating_from }} → {{ stock.rating_to }}</td>
            <td class="px-3 py-2 text-slate-700">{{ stock.target_from }} → {{ stock.target_to }}</td>
            <td class="px-3 py-2 text-slate-700">{{ new Date(stock.time).toLocaleString() }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>
