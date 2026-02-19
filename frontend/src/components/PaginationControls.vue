<script setup lang="ts">
defineProps<{
  total: number
  limit: number
  offset: number
  currentPage: number
  totalPages: number
  canGoPrev: boolean
  canGoNext: boolean
  disabled?: boolean
}>()

const emit = defineEmits<{
  prev: []
  next: []
}>()
</script>

<template>
  <div class="flex flex-wrap items-center justify-between gap-3 rounded-lg border border-slate-200 bg-white px-4 py-3 text-sm">
    <p class="text-slate-600">
      Showing {{ total === 0 ? 0 : offset + 1 }}–{{ Math.min(offset + limit, total) }} of {{ total }}
    </p>

    <div class="flex items-center gap-2">
      <span class="text-slate-600">Page {{ currentPage }} / {{ totalPages }}</span>
      <button
        class="rounded border border-slate-300 px-3 py-1 text-slate-700 hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-50"
        :disabled="disabled || !canGoPrev"
        @click="emit('prev')"
      >
        Prev
      </button>
      <button
        class="rounded border border-slate-300 px-3 py-1 text-slate-700 hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-50"
        :disabled="disabled || !canGoNext"
        @click="emit('next')"
      >
        Next
      </button>
    </div>
  </div>
</template>
