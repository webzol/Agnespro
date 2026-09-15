<script setup lang="ts">
import { onMounted } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();
// 用户反馈: 设置页面合并到后台, 自动跳转
onMounted(() => {
  router.replace({ name: "admin" });
});
</script>

<template>
  <section class="card glass" style="text-align:center; padding: 60px 20px;">
    <div class="spinner" style="margin: 0 auto 16px;"></div>
    <p class="muted">跳转到后台配置…</p>
  </section>
</template>
