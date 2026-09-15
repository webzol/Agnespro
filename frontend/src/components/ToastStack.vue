<script setup lang="ts">
import { useToast } from "@/composables/toast";
const { state } = useToast();
</script>

<template>
  <div class="toast-stack">
    <div v-for="t in state.items" :key="t.id" :class="['toast', t.kind]">{{ t.text }}</div>
  </div>
</template>
