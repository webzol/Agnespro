<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { jobsApi } from "@/api/endpoints";
import type { Job } from "@/types";

const router = useRouter();
const jobs = ref<Job[]>([]);
const search = ref("");
const loading = ref(true);

async function load() {
  loading.value = true;
  try {
    jobs.value = await jobsApi.list();
  } finally { loading.value = false; }
}
onMounted(load);

const filtered = (() => {
  const q = search.value.trim().toLowerCase();
  if (!q) return jobs.value;
  return jobs.value.filter(j => j.title.toLowerCase().includes(q));
})();

const STATUS_LABEL: Record<string,string> = { pending:"待处理", planning:"规划中", characters:"生成角色", props:"生成道具", scenes:"生成场景", video:"生成视频", done:"已完成", failed:"失败", cancelled:"已取消" };

async function del(id: string) {
  if (!confirm("删除任务？")) return;
  try { await jobsApi.del(id); await load(); } catch (e) { alert((e as Error).message); }
}

function fmt(d: string) { return d ? new Date(d).toLocaleString() : ""; }
</script>

<template>
  <section class="card glass">
    <div class="card-head">
      <h2>所有任务</h2>
      <div class="row">
        <input v-model="search" type="search" placeholder="搜索标题..." class="search">
        <button class="btn primary" @click="router.push({ name: 'new' })">+ 新建</button>
      </div>
    </div>
    <div v-if="loading" class="muted small" style="padding:24px;text-align:center;">加载中…</div>
    <div v-else-if="jobs.length === 0" class="muted small" style="padding:24px;text-align:center;">还没有任务,点击右上角新建。</div>
    <div v-else class="job-list">
      <div v-for="j in jobs" :key="j.id" class="job-row glass" @click="router.push({ name: 'detail', params: { id: j.id } })">
        <div class="job-row-main">
          <div class="job-title">{{ j.title }}</div>
          <div class="muted small">
            {{ j.num_episodes }} 集 · {{ j.num_characters }} 角色 · {{ j.num_scenes }} 场景 · {{ j.num_videos }} 视频 · {{ fmt(j.updated_at) }}
          </div>
        </div>
        <div class="job-row-meta">
          <span :class="['badge', `badge-${j.status}`]">{{ STATUS_LABEL[j.status] || j.status }}</span>
          <button class="btn ghost small" @click.stop="del(j.id)">删除</button>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.row { display: flex; gap: 8px; align-items: center; }
.search { background: rgba(255,255,255,.04); border: 1px solid var(--border); color: var(--text); border-radius: 10px; padding: 9px 12px; outline: none; min-width: 200px; }
.search:focus { border-color: var(--primary); }
.job-list { display: flex; flex-direction: column; gap: 10px; }
.job-row { display: flex; align-items: center; justify-content: space-between; gap: 16px; padding: 14px 18px; cursor: pointer; transition: all .15s; border-radius: 12px; }
.job-row:hover { transform: translateY(-1px); border-color: var(--border-strong); }
.job-row-main { flex: 1; }
.job-title { font-weight: 600; margin-bottom: 4px; }
.job-row-meta { display: flex; align-items: center; gap: 10px; }
.badge { padding: 3px 10px; border-radius: 999px; font-size: 11px; font-weight: 500; }
.badge-pending, .badge-planning { background: rgba(120,120,150,.18); color: #ccc; }
.badge-characters, .badge-props, .badge-scenes { background: rgba(99,102,241,.18); color: #a5b4fc; }
.badge-video { background: rgba(168,85,247,.18); color: #c4b5fd; }
.badge-done { background: rgba(74,222,128,.18); color: #86efac; }
.badge-failed, .badge-cancelled { background: rgba(248,113,113,.18); color: #fca5a5; }
.btn.small { padding: 4px 10px; font-size: 12px; }
</style>
