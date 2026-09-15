<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { jobsApi } from "@/api/endpoints";
import type { Job } from "@/types";

// Use Vite ?url to import images as URL strings (not ES modules)
import d1 from "@/assets/img/showcase/drama1.jpg?url";
import d2 from "@/assets/img/showcase/drama2.jpg?url";
import d3 from "@/assets/img/showcase/drama3.jpg?url";
import d4 from "@/assets/img/showcase/drama4.jpg?url";
import d5 from "@/assets/img/showcase/drama5.jpg?url";
import d6 from "@/assets/img/showcase/drama6.jpg?url";
const showcase = [d1, d2, d3, d4, d5, d6];

const router = useRouter();
const recent = ref<Job[]>([]);
const loading = ref(true);

onMounted(async () => {
  try {
    const all = await jobsApi.list();
    recent.value = all.slice(0, 4);
  } catch (e) { /* ignore */ }
  finally { loading.value = false; }
});

function fmt(d: string) {
  if (!d) return "";
  const dt = new Date(d);
  return dt.toLocaleString();
}
</script>

<template>
  <section class="hero glass">
    <h1 class="hero-title">从一段脚本到一段视频 <span class="dim">就这几步。</span></h1>
    <p class="hero-sub">粘贴剧本,系统会先拆出剧集,识别角色与道具,为每一集生成场景图,最后合成短剧。</p>
    <div class="hero-cta">
      <button class="btn primary" @click="router.push({ name: 'new' })">开始新创作</button>
      <button class="btn glass" @click="router.push({ name: 'jobs' })">查看我的作品</button>
    </div>
    <div class="model-strip">
      <span>文本 <code>agnes-2.5-flash</code></span>
      <span>图像 <code>agnes-image-2.1-flash</code></span>
      <span>视频 <code>agnes-video-v2.0</code></span>
    </div>
  </section>

  <section class="steps">
    <article class="step glass">
      <div class="step-num">1</div>
      <h3>粘贴剧本</h3>
      <p>支持中文 / 英文。可用「第N集」「Episode N」「EP N」标记剧集,系统会自动拆分。</p>
    </article>
    <article class="step glass">
      <div class="step-num">2</div>
      <h3>自动资产</h3>
      <p>提取主角、配角、关键道具与场景;为每场戏生成首帧图。</p>
    </article>
    <article class="step glass">
      <div class="step-num">3</div>
      <h3>逐集生成</h3>
      <p>单次生成全部剧集,或指定只生成某一集;支持取消、排队、断点续传。</p>
    </article>
  </section>

  <section class="showcase">
    <header class="showcase-head">
      <h2>精选作品</h2>
      <p class="muted">从都市短剧到神话史诗,从东方意境到好莱坞级视觉——这些示例展示了本平台能生成的多样化视觉风格。</p>
    </header>
    <h3 class="showcase-title">高端短剧系列 <span class="muted small">都市 / 悬疑 / 古风 / 职场</span></h3>
    <div class="showcase-grid">
      <div v-for="(src, i) in showcase" :key="i" class="showcase-card">
        <img :src="src" :alt="`showcase-${i}`" loading="lazy" />
      </div>
    </div>
  </section>

  <section class="recent glass">
    <header class="card-head">
      <h2>最近任务</h2>
      <button class="btn ghost" @click="router.push({ name: 'jobs' })">查看全部 →</button>
    </header>
    <div v-if="loading" class="muted small">加载中…</div>
    <div v-else-if="recent.length === 0" class="empty">暂无任务。</div>
    <div v-else class="recent-list">
      <div v-for="j in recent" :key="j.id" class="recent-row" @click="router.push({ name: 'detail', params: { id: j.id } })">
        <div class="recent-title">{{ j.title }}</div>
        <div class="muted small">{{ j.num_episodes }} 集 · {{ j.num_characters }} 角色 · {{ j.num_videos }} 视频 · {{ fmt(j.updated_at) }}</div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.hero { padding: 36px 32px; margin-bottom: 20px; }
.hero-title { font-size: 30px; font-weight: 800; letter-spacing: .01em; }
.hero-title .dim { color: var(--text-muted); font-weight: 600; }
.hero-sub { color: var(--text-soft); margin-top: 8px; max-width: 640px; }
.hero-cta { display: flex; gap: 12px; margin-top: 18px; }
.model-strip { display: flex; gap: 18px; margin-top: 18px; font-size: 12px; color: var(--text-muted); flex-wrap: wrap; }
.model-strip code { background: rgba(255,255,255,.06); padding: 2px 7px; border-radius: 6px; font-size: 11px; }

.steps { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 14px; margin: 24px 0; }
.step { padding: 20px; }
.step-num { width: 30px; height: 30px; border-radius: 8px; background: var(--grad); color: #fff; display: flex; align-items: center; justify-content: center; font-weight: 700; margin-bottom: 10px; }
.step h3 { font-size: 15px; margin-bottom: 6px; }
.step p { font-size: 12px; color: var(--text-soft); margin: 0; }

.showcase { margin: 24px 0; }
.showcase-head h2 { font-size: 20px; margin-bottom: 4px; }
.showcase-head p { font-size: 13px; color: var(--text-muted); max-width: 720px; }
.showcase-title { font-size: 14px; margin: 16px 0 10px; }
.showcase-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(170px, 1fr)); gap: 12px; }
.showcase-card { aspect-ratio: 3/4; border-radius: 12px; overflow: hidden; background: rgba(255,255,255,.04); border: 1px solid var(--border); }
.showcase-card img { width: 100%; height: 100%; object-fit: cover; transition: transform .3s; }
.showcase-card:hover img { transform: scale(1.04); }

.recent { padding: 20px; margin-top: 20px; }
.recent-list { display: flex; flex-direction: column; gap: 10px; }
.recent-row { padding: 12px 14px; border-radius: 10px; background: rgba(255,255,255,.03); border: 1px solid var(--border); cursor: pointer; transition: all .15s; }
.recent-row:hover { background: rgba(255,255,255,.06); border-color: var(--border-strong); }
.recent-title { font-weight: 600; margin-bottom: 4px; }
.empty { color: var(--text-muted); padding: 24px; text-align: center; font-size: 13px; }
</style>
