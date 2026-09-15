<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { jobsApi } from "@/api/endpoints";
import { useToast } from "@/composables/toast";
import type { Job } from "@/types";

const route = useRoute();
const router = useRouter();
const toast = useToast();
const id = computed(() => route.params.id as string);

const job = ref<Job | null>(null);
const loading = ref(true);
const pollHandle = ref<number | null>(null);
const STYLE_OPTIONS = [
  { id: "modern", name: "现代白" }, { id: "neon", name: "霓虹" },
  { id: "typewriter", name: "打字机" }, { id: "cinema", name: "影院金" },
  { id: "kawaii", name: "软萌粉" }, { id: "minimal", name: "极简黑" },
  { id: "news", name: "新闻蓝带" }, { id: "burn", name: "火焰" },
  { id: "ink", name: "水墨" }, { id: "anime", name: "二次元" },
];
const selectedStyle = ref("modern");

async function load() {
  try {
    job.value = await jobsApi.get(id.value);
  } catch (e) {
    toast.err((e as Error).message);
    job.value = null;
  } finally {
    loading.value = false;
  }
}

function startPolling() {
  stopPolling();
  pollHandle.value = window.setInterval(() => {
    if (!job.value) return;
    if (["done","failed","cancelled"].includes(job.value.status)) { stopPolling(); return; }
    load();
  }, 3000);
}
function stopPolling() {
  if (pollHandle.value) { clearInterval(pollHandle.value); pollHandle.value = null; }
}
onMounted(() => { load(); startPolling(); });
onUnmounted(stopPolling);
watch(() => route.params.id, () => { load(); startPolling(); });

const isRunning = computed(() => job.value && !["done","failed","cancelled"].includes(job.value.status));

async function runJob() {
  if (!job.value) return;
  try { await jobsApi.run(job.value.id); toast.ok("已提交运行"); await load(); }
  catch (e) { toast.err((e as Error).message); }
}
async function cancelJob() {
  if (!job.value || !confirm("确认取消?")) return;
  try { await jobsApi.cancel(job.value.id); toast.ok("已取消"); await load(); }
  catch (e) { toast.err((e as Error).message); }
}
async function delJob() {
  if (!job.value || !confirm("确认删除?")) return;
  try { await jobsApi.del(job.value.id); toast.ok("已删除"); router.push({ name: "jobs" }); }
  catch (e) { toast.err((e as Error).message); }
}

function fmt(d?: string) { return d ? new Date(d).toLocaleString() : "—"; }
function statusLabel(s: string) { return ({pending:"待处理", planning:"规划中", characters:"生成角色", props:"生成道具", scenes:"生成场景", video:"生成视频", done:"已完成", failed:"失败", cancelled:"已取消"})[s] || s; }
function epStateLabel(s: string) { return ({pending:"待处理", running:"进行中", characters:"生成角色", props:"生成道具", scenes:"生成场景", video:"生成视频", done:"已完成", failed:"失败", skipped:"已跳过"})[s] || s; }
</script>

<template>
  <section v-if="loading" class="card glass"><h2>加载中…</h2></section>
  <section v-else-if="!job" class="card glass">
    <h2>任务不存在或已删除</h2>
    <button class="btn primary" @click="router.push({ name: 'jobs' })">返回任务列表</button>
  </section>
  <section v-else class="job-detail">
    <div class="card glass">
      <div class="card-head">
        <div>
          <h2>{{ job.title }}</h2>
          <div class="muted small">ID: {{ job.id }} · {{ statusLabel(job.status) }}</div>
        </div>
        <div class="row">
          <button v-if="isRunning" class="btn ghost" @click="cancelJob">取消</button>
          <button class="btn primary" @click="runJob">{{ isRunning ? "重新运行" : "▶ 运行" }}</button>
          <button class="btn ghost" @click="delJob">删除</button>
        </div>
      </div>
      <div class="progress-wrap">
        <div class="progress-bar"><div :style="{ width: (job.progress || 0) + '%' }"></div></div>
        <span class="muted small">{{ job.progress || 0 }}%</span>
      </div>
      <div v-if="job.error" class="err-banner">{{ job.error }}</div>
      <div class="stats">
        <div class="stat"><span class="num">{{ job.num_episodes }}</span><span class="lbl">剧集</span></div>
        <div class="stat"><span class="num">{{ job.num_characters }}</span><span class="lbl">角色</span></div>
        <div class="stat"><span class="num">{{ job.num_props }}</span><span class="lbl">道具</span></div>
        <div class="stat"><span class="num">{{ job.num_scenes }}</span><span class="lbl">场景</span></div>
        <div class="stat"><span class="num">{{ job.num_videos }}</span><span class="lbl">视频</span></div>
      </div>
      <div class="model-strip">
        <span>视觉风格: <code>{{ job.style || "(默认)" }}</code></span>
        <span>剧本模型: <code>{{ job.script_model || "(默认)" }}</code></span>
        <span>图像模型: <code>{{ job.image_model || "(默认)" }}</code></span>
        <span>视频模型: <code>{{ job.video_model || "(默认)" }}</code></span>
        <span>比例: <code>{{ job.aspect_ratio || "(默认)" }}</code></span>
        <span v-if="job.genre_name">题材: <code>{{ job.genre_name }}</code></span>
      </div>
    </div>

    <div v-if="job.characters.length" class="card glass">
      <h3>角色</h3>
      <div class="asset-grid">
        <div v-for="c in job.characters" :key="c.id" class="asset">
          <img v-if="c.image_url" :src="c.image_url" :alt="c.name" loading="lazy">
          <div v-else class="pending">未生成</div>
          <div class="asset-name">{{ c.name }}</div>
        </div>
      </div>
    </div>

    <div v-if="job.props.length" class="card glass">
      <h3>道具</h3>
      <div class="asset-grid">
        <div v-for="p in job.props" :key="p.id" class="asset">
          <img v-if="p.image_url" :src="p.image_url" :alt="p.name" loading="lazy">
          <div v-else class="pending">未生成</div>
          <div class="asset-name">{{ p.name }}</div>
        </div>
      </div>
    </div>

    <div v-for="ep in job.episodes" :key="ep.id" class="card glass episode">
      <div class="episode-head">
        <h4>第{{ ep.index }}集: {{ ep.title }}</h4>
        <span :class="['badge', `badge-${ep.state}`]">{{ epStateLabel(ep.state) }} · {{ ep.progress }}%</span>
      </div>
      <div v-if="ep.body" class="ep-body muted small">{{ ep.body }}</div>
      <div v-if="ep.video_url" class="ep-video">
        <video :src="ep.video_url" controls></video>
        <a class="btn glass small" :href="`/api/jobs/${job.id}/subtitles?style=${selectedStyle}`" target="_blank">下载字幕 ({{ STYLE_OPTIONS.find(s => s.id === selectedStyle)?.name }})</a>
      </div>
      <div v-if="ep.scenes.length" class="scene-grid">
        <div v-for="s in ep.scenes" :key="s.id" class="scene">
          <img v-if="s.image_url" :src="s.image_url" :alt="s.heading" loading="lazy">
          <div v-else class="pending">场景 {{ s.index }}</div>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.card-head { display: flex; align-items: center; justify-content: space-between; gap: 12px; flex-wrap: wrap; margin-bottom: 16px; }
.card-head h2 { font-size: 18px; }
.row { display: flex; gap: 8px; flex-wrap: wrap; }
.progress-wrap { display: flex; align-items: center; gap: 12px; margin-bottom: 14px; }
.progress-bar { flex: 1; height: 6px; background: rgba(255,255,255,.06); border-radius: 999px; overflow: hidden; }
.progress-bar > div { height: 100%; background: var(--grad); transition: width .4s; }
.err-banner { padding: 12px; background: rgba(248,113,113,.12); border: 1px solid rgba(248,113,113,.3); border-radius: 10px; color: #fca5a5; font-size: 13px; margin-bottom: 14px; }
.stats { display: grid; grid-template-columns: repeat(5, 1fr); gap: 10px; margin-bottom: 14px; }
.stat { text-align: center; padding: 12px; background: rgba(255,255,255,.04); border-radius: 10px; }
.stat .num { display: block; font-size: 24px; font-weight: 700; }
.stat .lbl { font-size: 11px; color: var(--text-muted); }
.model-strip { display: flex; gap: 14px; flex-wrap: wrap; font-size: 12px; color: var(--text-muted); padding-top: 8px; border-top: 1px solid var(--border); }
.model-strip code { background: rgba(255,255,255,.06); padding: 2px 6px; border-radius: 5px; font-size: 11px; }
h3 { font-size: 15px; margin-bottom: 12px; }
.asset-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(120px, 1fr)); gap: 10px; }
.asset { aspect-ratio: 3/4; border-radius: 10px; overflow: hidden; background: rgba(255,255,255,.04); border: 1px solid var(--border); position: relative; }
.asset img { width: 100%; height: 100%; object-fit: cover; }
.asset .pending { height: 100%; display: flex; align-items: center; justify-content: center; color: var(--text-muted); font-size: 12px; }
.asset-name { position: absolute; bottom: 0; left: 0; right: 0; padding: 6px 8px; background: linear-gradient(to top, rgba(0,0,0,.85), transparent); font-size: 11px; }
.episode { margin-top: 14px; }
.episode-head { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; }
.episode-head h4 { font-size: 14px; }
.ep-body { margin-bottom: 10px; max-height: 80px; overflow: hidden; }
.badge { padding: 3px 10px; border-radius: 999px; font-size: 11px; font-weight: 500; }
.badge-running, .badge-video { background: rgba(168,85,247,.18); color: #c4b5fd; }
.badge-done { background: rgba(74,222,128,.18); color: #86efac; }
.badge-failed, .badge-skipped { background: rgba(248,113,113,.18); color: #fca5a5; }
.badge-pending, .badge-characters, .badge-props, .badge-scenes { background: rgba(99,102,241,.18); color: #a5b4fc; }
.ep-video { margin: 12px 0; display: flex; flex-direction: column; gap: 10px; }
.ep-video video { width: 100%; max-width: 480px; border-radius: 10px; }
.btn.small { padding: 4px 10px; font-size: 12px; align-self: flex-start; }
.scene-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(150px, 1fr)); gap: 8px; margin-top: 10px; }
.scene { aspect-ratio: 16/9; border-radius: 8px; overflow: hidden; background: rgba(255,255,255,.04); border: 1px solid var(--border); }
.scene img { width: 100%; height: 100%; object-fit: cover; }
.scene .pending { height: 100%; display: flex; align-items: center; justify-content: center; color: var(--text-muted); font-size: 11px; }
.job-detail { display: flex; flex-direction: column; gap: 14px; }
.card { padding: 20px; }
@media (max-width: 720px) { .stats { grid-template-columns: repeat(2, 1fr); } }
</style>
