<template>
  <div class="min-h-screen">
    <header class="border-b border-slate-200 bg-white">
      <div class="mx-auto flex max-w-7xl items-center justify-between px-4 py-4">
        <RouterLink to="/" class="text-lg font-semibold text-slate-900">Stock Analysis</RouterLink>
        <nav class="flex items-center gap-4 text-sm">
          <RouterLink to="/" class="text-slate-600 hover:text-slate-900">Dashboard</RouterLink>
        </nav>
      </div>
    </header>

    <main class="mx-auto max-w-7xl px-4 py-6">
      <RouterView />
    </main>
  </div>
</template>
