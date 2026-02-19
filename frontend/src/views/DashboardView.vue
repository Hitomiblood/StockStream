<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'

import PaginationControls from '@/components/PaginationControls.vue'
import RecommendationsPanel from '@/components/RecommendationsPanel.vue'
import SettingsSyncSection from '@/components/SettingsSyncSection.vue'
import StockTable from '@/components/StockTable.vue'
import { useMetadataStore } from '@/stores/metadata'
import { SORT_FIELDS, useStocksStore } from '@/stores/stocks'
import { useRecommendationsStore } from '@/stores/recommendations'
import type { Stock } from '@/types/domain'

const router = useRouter()

const stocksStore = useStocksStore()
const metadataStore = useMetadataStore()
const recommendationsStore = useRecommendationsStore()

const { items, loading, error, total, limit, offset, sort, order, query, actionFilter, ratingFilter, currentPage, totalPages, canGoPrev, canGoNext, mode } =
  storeToRefs(stocksStore)
const { actions, ratings, loading: metadataLoading, error: metadataError } = storeToRefs(metadataStore)
const {
  items: recommendationItems,
  loading: recommendationsLoading,
  error: recommendationsError,
  generatedAt
} = storeToRefs(recommendationsStore)

const isSearchMode = computed(() => mode.value === 'search')

onMounted(async () => {
  await Promise.all([
    stocksStore.fetchStocks(),
    metadataStore.fetchMetadata(),
    recommendationsStore.fetchRecommendations(10)
  ])
})

watch([offset, sort, order], async () => {
  if (!stocksStore.hasSearch) {
    await stocksStore.fetchStocks()
  }
})

async function onSearch(): Promise<void> {
  await stocksStore.fetchStocks()
}

async function onClearSearch(): Promise<void> {
  stocksStore.setQuery('')
  await stocksStore.fetchStocks()
}

async function onApplyFilters(): Promise<void> {
  stocksStore.setQuery('')
  await stocksStore.fetchStocks()
}

async function onClearFilters(): Promise<void> {
  stocksStore.clearFilters()
  await stocksStore.fetchStocks()
}

async function onPrevPage(): Promise<void> {
  stocksStore.prevPage()
  await stocksStore.fetchStocks()
}

async function onNextPage(): Promise<void> {
  stocksStore.nextPage()
  await stocksStore.fetchStocks()
}

function openDetail(stock: Stock): void {
  router.push({ name: 'stock-detail', params: { id: stock.id } })
}
</script>

<template>
  <div class="space-y-6">
    <section class="rounded-lg border border-slate-200 bg-white p-4">
      <h1 class="text-lg font-semibold text-slate-900">Dashboard</h1>
      <p class="mt-1 text-sm text-slate-600">Browse, search, sort, filter and inspect stock updates.</p>

      <div class="mt-4 grid gap-3 md:grid-cols-2 xl:grid-cols-4">
        <label class="flex flex-col gap-1 text-sm">
          <span class="text-slate-600">Search</span>
          <input
            :value="query"
            type="text"
            class="rounded border border-slate-300 px-3 py-2"
            placeholder="Ticker or company"
            @input="stocksStore.setQuery(($event.target as HTMLInputElement).value)"
            @keyup.enter="onSearch"
          />
        </label>

        <label class="flex flex-col gap-1 text-sm">
          <span class="text-slate-600">Sort by</span>
          <select
            :value="sort"
            class="rounded border border-slate-300 px-3 py-2"
            @change="stocksStore.setSort(($event.target as HTMLSelectElement).value as (typeof SORT_FIELDS)[number])"
          >
            <option v-for="field in SORT_FIELDS" :key="field" :value="field">{{ field }}</option>
          </select>
        </label>

        <label class="flex flex-col gap-1 text-sm">
          <span class="text-slate-600">Order</span>
          <select
            :value="order"
            class="rounded border border-slate-300 px-3 py-2"
            @change="stocksStore.setOrder(($event.target as HTMLSelectElement).value as 'asc' | 'desc')"
          >
            <option value="desc">desc</option>
            <option value="asc">asc</option>
          </select>
        </label>

        <label class="flex flex-col gap-1 text-sm">
          <span class="text-slate-600">Limit</span>
          <select
            :value="limit"
            class="rounded border border-slate-300 px-3 py-2"
            @change="stocksStore.limit = Number(($event.target as HTMLSelectElement).value); stocksStore.offset = 0; stocksStore.fetchStocks()"
          >
            <option :value="10">10</option>
            <option :value="20">20</option>
            <option :value="50">50</option>
            <option :value="100">100</option>
          </select>
        </label>
      </div>

      <div class="mt-3 flex flex-wrap items-center gap-2">
        <button class="rounded bg-slate-900 px-3 py-2 text-sm text-white hover:bg-slate-800" :disabled="loading" @click="onSearch">
          Search
        </button>
        <button class="rounded border border-slate-300 px-3 py-2 text-sm text-slate-700 hover:bg-slate-50" :disabled="loading || !query" @click="onClearSearch">
          Clear search
        </button>
        <span v-if="isSearchMode" class="text-xs text-slate-500">Search mode disables pagination and uses endpoint `/stocks/search`.</span>
      </div>

      <div class="mt-5 grid gap-3 md:grid-cols-2">
        <label class="flex flex-col gap-1 text-sm">
          <span class="text-slate-600">Action filter</span>
          <select
            :value="actionFilter"
            class="rounded border border-slate-300 px-3 py-2"
            :disabled="metadataLoading"
            @change="stocksStore.setActionFilter(($event.target as HTMLSelectElement).value)"
          >
            <option value="">All actions</option>
            <option v-for="action in actions" :key="action" :value="action">{{ action }}</option>
          </select>
        </label>

        <label class="flex flex-col gap-1 text-sm">
          <span class="text-slate-600">Rating filter</span>
          <select
            :value="ratingFilter"
            class="rounded border border-slate-300 px-3 py-2"
            :disabled="metadataLoading"
            @change="stocksStore.setRatingFilter(($event.target as HTMLSelectElement).value)"
          >
            <option value="">All ratings</option>
            <option v-for="rating in ratings" :key="rating" :value="rating">{{ rating }}</option>
          </select>
        </label>
      </div>

      <div class="mt-3 flex flex-wrap items-center gap-2">
        <button class="rounded bg-slate-900 px-3 py-2 text-sm text-white hover:bg-slate-800" :disabled="loading || (!actionFilter && !ratingFilter)" @click="onApplyFilters">
          Apply filters
        </button>
        <button class="rounded border border-slate-300 px-3 py-2 text-sm text-slate-700 hover:bg-slate-50" :disabled="loading" @click="onClearFilters">
          Clear filters
        </button>
      </div>

      <p v-if="metadataError" class="mt-2 text-sm text-red-600">{{ metadataError }}</p>
    </section>

    <StockTable :items="items" :loading="loading" :error="error" @row-click="openDetail" />

    <PaginationControls
      :total="total"
      :limit="limit"
      :offset="offset"
      :current-page="currentPage"
      :total-pages="totalPages"
      :can-go-prev="canGoPrev"
      :can-go-next="canGoNext"
      :disabled="loading || isSearchMode"
      @prev="onPrevPage"
      @next="onNextPage"
    />

    <div class="grid gap-6 xl:grid-cols-2">
      <RecommendationsPanel
        :items="recommendationItems"
        :loading="recommendationsLoading"
        :error="recommendationsError"
        :generated-at="generatedAt"
      />
      <SettingsSyncSection />
    </div>
  </div>
</template>
