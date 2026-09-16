<script setup lang="ts">
import { RouterLink } from "vue-router";

const groups = [
  {
    title: "工作台",
    items: [
      { name: "home",      label: "首页",        badge: "" },
      { name: "jobs",      label: "作品管理",    badge: "" },
      { name: "new",       label: "新建短剧",    badge: "AI" },
    ],
  },
  {
    title: "资源",
    items: [
      { name: "templates", label: "剧本模板",    badge: "8" },
      { name: "docs",      label: "使用文档",    badge: "" },
    ],
  },
  {
    title: "管理",
    items: [
      { name: "admin",     label: "后台管理",    badge: "" },
    ],
  },
];
</script>

<template>
  <aside class="sidebar">
    <div class="sidebar-brand">
      <div class="brand-mark">A</div>
      <div>
        <div class="brand-name">Agnes AI Studio</div>
        <div class="brand-sub">短剧智能生成工作台</div>
      </div>
    </div>

    <nav class="sidebar-nav" v-for="g in groups" :key="g.title">
      <div class="sidebar-section">{{ g.title }}</div>
      <RouterLink
        v-for="t in g.items"
        :key="t.name"
        :to="{ name: t.name }"
        class="nav-item"
        active-class="active"
        :exact-active-class="t.name === 'home' ? 'active' : ''"
      >
        <!-- Icons (minimal flat stroke) -->
        <svg v-if="t.name === 'home'"     class="nav-icon" viewBox="0 0 24 24"><path d="M3 11l9-8 9 8v10a2 2 0 0 1-2 2h-4v-7h-6v7H5a2 2 0 0 1-2-2z"/></svg>
        <svg v-else-if="t.name === 'jobs'" class="nav-icon" viewBox="0 0 24 24"><rect x="3" y="4" width="18" height="16" rx="2"/><path d="M3 9h18M8 4v16"/></svg>
        <svg v-else-if="t.name === 'new'"  class="nav-icon" viewBox="0 0 24 24"><circle cx="12" cy="12" r="9"/><path d="M12 8v8M8 12h8"/></svg>
        <svg v-else-if="t.name === 'templates'" class="nav-icon" viewBox="0 0 24 24"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="9" y1="13" x2="15" y2="13"/><line x1="9" y1="17" x2="13" y2="17"/></svg>
        <svg v-else-if="t.name === 'docs'" class="nav-icon" viewBox="0 0 24 24"><path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"/><path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"/></svg>
        <svg v-else-if="t.name === 'admin'" class="nav-icon" viewBox="0 0 24 24"><path d="M12 2l9 4v6c0 5-3.5 9.4-9 10-5.5-.6-9-5-9-10V6z"/><path d="M9 12l2 2 4-4"/></svg>

        <span>{{ t.label }}</span>
        <span v-if="t.badge" class="nav-badge">{{ t.badge }}</span>
      </RouterLink>
    </nav>

    <div class="sidebar-footer">
      <div class="sidebar-user">
        <div class="avatar">创</div>
        <div>
          <div class="u-name">创作者</div>
          <div class="u-plan">免费版 · 升级</div>
        </div>
      </div>
    </div>
  </aside>
</template>
