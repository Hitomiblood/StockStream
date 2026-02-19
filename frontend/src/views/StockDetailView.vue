<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'

import { getStockById, getStocksByTicker } from '@/api/stocks'
import type { Stock } from '@/types/domain'

const route = useRoute()

const stock = ref<Stock | null>(null)
const history = ref<Stock[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const historyLoading = ref(false)
const historyError = ref<string | null>(null)

const stockId = computed(() => String(route.params.id || ''))

onMounted(async () => {
  await loadStock()
})

async function loadStock(): Promise<void> {
  loading.value = true
  error.value = null

  try {
    const data = await getStockById(stockId.value)
    stock.value = data
    await loadHistory(data.ticker)
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load stock detail'
  } finally {
    loading.value = false
  }
}

async function loadHistory(ticker: string): Promise<void> {
  historyLoading.value = true
  historyError.value = null

  try {
    const response = await getStocksByTicker(ticker)
    history.value = response.history
  } catch (err) {
    historyError.value = err instanceof Error ? err.message : 'Failed to load ticker history'
    history.value = []
  } finally {
    historyLoading.value = false
  }
}
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-lg font-semibold text-slate-900">Stock detail</h1>
      <RouterLink to="/" class="text-sm text-slate-600 hover:text-slate-900">← Back to dashboard</RouterLink>
    </div>

    <section v-if="loading" class="rounded-lg border border-slate-200 bg-white p-4">
      <div class="h-8 w-1/3 animate-pulse rounded bg-slate-100"></div>
      <div class="mt-3 h-20 animate-pulse rounded bg-slate-100"></div>
    </section>

    <section v-else-if="error" class="rounded-lg border border-slate-200 bg-white p-4 text-sm text-red-600">
      {{ error }}
    </section>

    <section v-else-if="stock" class="rounded-lg border border-slate-200 bg-white p-4">
      <h2 class="text-base font-semibold text-slate-900">{{ stock.ticker }} · {{ stock.company }}</h2>
      <div class="mt-4 grid gap-3 text-sm md:grid-cols-2">
        <p><span class="font-medium text-slate-800">Action:</span> {{ stock.action }}</p>
        <p><span class="font-medium text-slate-800">Brokerage:</span> {{ stock.brokerage }}</p>
        <p><span class="font-medium text-slate-800">Rating:</span> {{ stock.rating_from }} → {{ stock.rating_to }}</p>
        <p><span class="font-medium text-slate-800">Target:</span> {{ stock.target_from }} → {{ stock.target_to }}</p>
        <p><span class="font-medium text-slate-800">Time:</span> {{ new Date(stock.time).toLocaleString() }}</p>
        <p><span class="font-medium text-slate-800">ID:</span> {{ stock.id }}</p>
      </div>
    </section>

    <section class="rounded-lg border border-slate-200 bg-white">
      <header class="border-b border-slate-200 px-4 py-3">
        <h2 class="text-sm font-semibold text-slate-900">Ticker history</h2>
      </header>

      <div v-if="historyLoading" class="space-y-2 p-4">
        <div v-for="idx in 4" :key="idx" class="h-10 animate-pulse rounded bg-slate-100"></div>
      </div>

      <p v-else-if="historyError" class="p-4 text-sm text-red-600">{{ historyError }}</p>

      <p v-else-if="history.length === 0" class="p-4 text-sm text-slate-500">No historical data found.</p>

      <div v-else class="overflow-x-auto">
        <table class="min-w-full divide-y divide-slate-200 text-sm">
          <thead class="bg-slate-50 text-left text-xs uppercase tracking-wide text-slate-500">
            <tr>
              <th class="px-3 py-2">Time</th>
              <th class="px-3 py-2">Action</th>
              <th class="px-3 py-2">Rating</th>
              <th class="px-3 py-2">Target</th>
              <th class="px-3 py-2">Brokerage</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr v-for="item in history" :key="item.id">
              <td class="px-3 py-2 text-slate-700">{{ new Date(item.time).toLocaleString() }}</td>
              <td class="px-3 py-2 text-slate-700">{{ item.action }}</td>
              <td class="px-3 py-2 text-slate-700">{{ item.rating_from }} → {{ item.rating_to }}</td>
              <td class="px-3 py-2 text-slate-700">{{ item.target_from }} → {{ item.target_to }}</td>
              <td class="px-3 py-2 text-slate-700">{{ item.brokerage }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>
  </div>
</template>
