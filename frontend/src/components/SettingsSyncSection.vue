<script setup lang="ts">
import { storeToRefs } from 'pinia'

import { useSyncStore } from '@/stores/sync'

const syncStore = useSyncStore()
const { loading, error, result } = storeToRefs(syncStore)

async function onSync(): Promise<void> {
  const confirmed = window.confirm('This will fetch and synchronize stocks from external API. Continue?')
  if (!confirmed) {
    return
  }

  await syncStore.runSync()
}
</script>

<template>
  <section class="rounded-lg border border-slate-200 bg-white p-4">
    <h2 class="text-sm font-semibold text-slate-900">Settings</h2>
    <p class="mt-1 text-sm text-slate-600">Run a manual data synchronization. This operation can take time.</p>

    <div class="mt-4">
      <button
        class="rounded bg-slate-900 px-4 py-2 text-sm font-medium text-white hover:bg-slate-800 disabled:cursor-not-allowed disabled:opacity-60"
        :disabled="loading"
        @click="onSync"
      >
        {{ loading ? 'Synchronizing...' : 'Sync Stocks Now' }}
      </button>
    </div>

    <p v-if="error" class="mt-3 text-sm text-red-600">{{ error }}</p>

    <div v-if="result" class="mt-3 rounded border border-slate-200 bg-slate-50 p-3 text-sm text-slate-700">
      <p class="font-medium">{{ result.message }}</p>
      <p>New: {{ result.total_new }} · Updated: {{ result.total_updated }} · Total: {{ result.total_fetched }}</p>
      <p>Duration: {{ result.duration_ms }} ms</p>
    </div>
  </section>
</template>
