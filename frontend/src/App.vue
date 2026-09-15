<script setup lang="ts">
import { RouterView } from "vue-router";
import NavBar from "./components/NavBar.vue";
import ToastStack from "./components/ToastStack.vue";
</script>

<template>
  <div class="aurora">
    <div class="aurora-blob a1"></div>
    <div class="aurora-blob a2"></div>
    <div class="aurora-blob a3"></div>
  </div>
  <header class="topbar">
    <div class="brand">
      <div class="brand-mark">A</div>
      <div>
        <div class="brand-name">Agnes AI Studio</div>
        <div class="brand-sub">把剧本变成视频,简单几步</div>
      </div>
    </div>
    <NavBar />
    <div class="topbar-right">
      <div id="queue-pill" class="pill glass hidden">
        <span class="dot"></span><span id="queue-text">队列</span>
      </div>
      <button id="theme-toggle" class="theme-toggle" aria-label="Toggle theme" title="Toggle theme">☀</button>
    </div>
  </header>

  <main>
    <RouterView />
  </main>

  <ToastStack />
</template>
