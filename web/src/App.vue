<script setup lang="ts"></script>

<template>
  <div class="min-h-screen">
    <header class="sticky top-0 z-10 border-b border-border bg-paper/85 backdrop-blur">
      <div class="mx-auto flex h-16 max-w-3xl items-center justify-between px-6">
        <RouterLink
          to="/"
          class="font-serif text-xl font-bold tracking-wide text-ink transition-colors hover:text-primary"
        >
          笔耕
        </RouterLink>
        <span class="font-mono text-xs uppercase tracking-widest text-ink-faint">Blog</span>
      </div>
    </header>

    <main class="mx-auto max-w-3xl px-6 py-12">
      <RouterView />
    </main>

    <footer class="border-t border-border">
      <div class="mx-auto max-w-3xl px-6 py-8 text-center font-mono text-xs tracking-wider text-ink-faint">
        笔耕 · 个人博客
      </div>
    </footer>
  </div>
</template>
