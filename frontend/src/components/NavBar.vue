<script setup lang="ts">
import { RouterLink } from "vue-router";
const tabs = [
  { name: "home",     label: "首页" },
  { name: "new",      label: "新建任务" },
  { name: "jobs",     label: "任务列表" },
  { name: "settings", label: "设置" },
  { name: "admin",    label: "⚙ 后台" },
  { name: "docs",     label: "使用文档" },
];
</script>

<template>
  <nav class="tabs glass" id="nav">
    <RouterLink v-for="t in tabs" :key="t.name" :to="{ name: t.name }"
                class="tab" active-class="active" :exact-active-class="t.name === 'home' ? 'active' : ''">
      {{ t.label }}
    </RouterLink>
  </nav>
</template>
