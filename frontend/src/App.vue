<script setup lang="ts">
import { RouterView } from "vue-router";
import NavBar from "./components/NavBar.vue";
import ToastStack from "./components/ToastStack.vue";
</script>

<template>
  <div class="bg-deco"></div>
  <div class="app-shell">
    <NavBar />
    <div class="main">
      <RouterView v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </RouterView>
    </div>
  </div>
  <ToastStack />
</template>
