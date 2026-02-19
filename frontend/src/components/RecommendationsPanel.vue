<script setup lang="ts">
import type { Recommendation } from '@/types/domain'

defineProps<{
  items: Recommendation[]
  loading: boolean
  error: string | null
  generatedAt: string
}>()

function confidenceBadge(confidence: Recommendation['confidence']): string {
  if (confidence === 'high') return 'bg-emerald-100 text-emerald-700'
  if (confidence === 'medium') return 'bg-amber-100 text-amber-700'
  return 'bg-slate-100 text-slate-700'
}
</script>

<template>
  <section class="rounded-lg border border-slate-200 bg-white">
    <header class="border-b border-slate-200 px-4 py-3">
      <h2 class="text-sm font-semibold text-slate-900">Recommendations</h2>
      <p v-if="generatedAt" class="mt-1 text-xs text-slate-500">Generated at: {{ new Date(generatedAt).toLocaleString() }}</p>
    </header>

    <div v-if="loading" class="space-y-2 p-4">
      <div v-for="idx in 3" :key="idx" class="h-20 animate-pulse rounded bg-slate-100"></div>
    </div>

    <div v-else-if="error" class="p-4 text-sm text-red-600">{{ error }}</div>

    <div v-else-if="items.length === 0" class="p-4 text-sm text-slate-500">No recommendations available.</div>

    <ul v-else class="divide-y divide-slate-100">
      <li v-for="item in items" :key="`${item.stock.id}-${item.score}`" class="space-y-2 px-4 py-3">
        <div class="flex items-center justify-between gap-2">
          <p class="font-medium text-slate-900">{{ item.stock.ticker }} · {{ item.stock.company }}</p>
          <span class="rounded px-2 py-1 text-xs font-medium" :class="confidenceBadge(item.confidence)">
            {{ item.confidence }}
          </span>
        </div>
        <p class="text-sm text-slate-700">Score: {{ item.score.toFixed(2) }}</p>
        <p class="text-sm text-slate-600">{{ item.reason }}</p>
      </li>
    </ul>
  </section>
</template>
