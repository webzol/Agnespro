<script setup lang="ts">
import { RouterLink } from "vue-router";

const groups = [
  {
    title: "工作台",
    items: [
      { name: "home",     label: "首页",        badge: "" },
      { name: "jobs",     label: "作品管理",    badge: "" },
      { name: "new",      label: "新建短剧",    badge: "AI" },
    ],
  },
  {
    title: "资源",
    items: [
      { name: "docs",     label: "剧本模板",    badge: "" },
      { name: "settings", label: "账户设置",    badge: "" },
    ],
  },
  {
    title: "管理",
    items: [
      { name: "admin",    label: "后台管理",    badge: "" },
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
        <svg v-else-if="t.name === 'docs'" class="nav-icon" viewBox="0 0 24 24"><path d="M14 3H7a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V8z"/><path d="M14 3v5h5M9 13h6M9 17h6M9 9h2"/></svg>
        <svg v-else-if="t.name === 'settings'" class="nav-icon" viewBox="0 0 24 24"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>
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
