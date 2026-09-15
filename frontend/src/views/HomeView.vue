<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useRouter } from "vue-router";
import { jobsApi } from "@/api/endpoints";
import { useToast } from "@/composables/toast";
import type { Job } from "@/types";

import d1 from "@/assets/img/showcase/drama1.jpg?url";
import d2 from "@/assets/img/showcase/drama2.jpg?url";
import d3 from "@/assets/img/showcase/drama3.jpg?url";
import d4 from "@/assets/img/showcase/drama4.jpg?url";
import d5 from "@/assets/img/showcase/drama5.jpg?url";
import d6 from "@/assets/img/showcase/drama6.jpg?url";

const router = useRouter();
const toast = useToast();

// === Demo data: latest completed job ===
const job = ref<Job | null>(null);
const recent = ref<Job[]>([]);

const episodes = [
  { id: 1, title: "命运的开端", cover: d1, duration: "00:05", size: "12.4 MB", genre: "都市悬疑", status: "ready" },
  { id: 2, title: "雨夜的告白", cover: d2, duration: "00:05", size: "11.8 MB", genre: "都市悬疑", status: "ready" },
  { id: 3, title: "真相浮出",   cover: d3, duration: "00:05", size: "13.1 MB", genre: "都市悬疑", status: "ready" },
  { id: 4, title: "再次抉择",   cover: d4, duration: "00:05", size: "12.7 MB", genre: "都市悬疑", status: "ready" },
  { id: 5, title: "风暴前夕",   cover: d5, duration: "00:05", size: "12.2 MB", genre: "都市悬疑", status: "ready" },
  { id: 6, title: "终章·破晓",  cover: d6, duration: "00:05", size: "14.0 MB", genre: "都市悬疑", status: "ready" },
];

const selected = ref<number[]>([]);
const filterTab = ref<"all" | "ready" | "processing">("all");

const filteredEpisodes = computed(() => {
  if (filterTab.value === "all") return episodes;
  return episodes.filter(e => e.status === filterTab.value);
});

const allSelected = computed({
  get: () => selected.value.length === filteredEpisodes.value.length && filteredEpisodes.value.length > 0,
  set: (v: boolean) => {
    if (v) selected.value = filteredEpisodes.value.map(e => e.id);
    else selected.value = [];
  },
});

function toggle(id: number) {
  const i = selected.value.indexOf(id);
  if (i >= 0) selected.value.splice(i, 1);
  else selected.value.push(id);
}

const totalSize = computed(() => {
  const sel = episodes.filter(e => selected.value.includes(e.id));
  const mb = sel.reduce((s, e) => s + parseFloat(e.size), 0);
  return mb.toFixed(1);
});

// === Download settings ===
const resolution = ref("1080p");
const format = ref("mp4");
const quality = ref("high");
const fps = ref("30");
const includeSubtitle = ref(true);
const includeIntro = ref(false);
const includeWatermark = ref(false);
const packMode = ref<"zip" | "single">("zip");
const notify = ref(true);

function confirmDownload() {
  if (selected.value.length === 0) {
    toast.err("请先选择要下载的剧集");
    return;
  }
  toast.ok(`已提交下载任务,共 ${selected.value.length} 集 (${totalSize.value} MB)`);
}

function newJob() {
  router.push({ name: "new" });
}

onMounted(async () => {
  try {
    const all = await jobsApi.list();
    recent.value = all.slice(0, 3);
    job.value = all[0] ?? null;
  } catch (e) { /* ignore */ }
});
</script>

<template>
  <header class="topbar">
    <div>
      <div class="topbar-title">我的短剧作品</div>
      <div class="topbar-sub">工作台 / 作品管理 / 《都市迷局》短剧系列</div>
    </div>
    <div class="topbar-right">
      <button class="btn ghost">
        <svg class="btn-icon" viewBox="0 0 24 24"><circle cx="11" cy="11" r="7"/><path d="M21 21l-4.35-4.35"/></svg>
        搜索
      </button>
      <button class="btn outline" @click="newJob">
        <svg class="btn-icon" viewBox="0 0 24 24"><path d="M12 5v14M5 12h14"/></svg>
        新建短剧
      </button>
      <button class="btn primary" :disabled="selected.length === 0" @click="confirmDownload">
        <svg class="btn-icon" viewBox="0 0 24 24"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
        批量下载{{ selected.length > 0 ? ` (${selected.length})` : "" }}
      </button>
    </div>
  </header>

  <div>
    <!-- Summary card -->
    <div class="card summary">
      <div class="summary-main">
        <div class="summary-cover">
          <img :src="d1" alt="cover" />
          <div class="summary-badge">已完成</div>
        </div>
        <div class="summary-info">
          <div class="summary-eyebrow">
            <span class="chip outline-grad"><span class="chip-dot"></span>都市悬疑</span>
            <span class="chip">高清 1080P</span>
            <span class="chip">5 秒/集</span>
          </div>
          <h1 class="summary-title">都市迷局 · 6 集短剧系列</h1>
          <p class="summary-desc">一个发生在都市雨夜的悬疑故事,从一份匿名邮件开始,引出尘封十年的真相。本剧由 agnes-2.5-flash 文本模型 + agnes-image-2.1-flash 图像模型 + agnes-video-v2.0 视频模型联合生成。</p>
          <div class="summary-meta">
            <div class="meta-item">
              <div class="meta-num">6</div>
              <div class="meta-label">剧集</div>
            </div>
            <div class="meta-divider"></div>
            <div class="meta-item">
              <div class="meta-num">4</div>
              <div class="meta-label">角色</div>
            </div>
            <div class="meta-divider"></div>
            <div class="meta-item">
              <div class="meta-num">12</div>
              <div class="meta-label">道具</div>
            </div>
            <div class="meta-divider"></div>
            <div class="meta-item">
              <div class="meta-num">76.2</div>
              <div class="meta-label">总大小(MB)</div>
            </div>
            <div class="meta-divider"></div>
            <div class="meta-item">
              <div class="meta-num">3'42"</div>
              <div class="meta-label">耗时</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Filter row -->
    <div class="filter-row">
      <div class="filter-tabs">
        <button class="filter-tab" :class="{ active: filterTab === 'all' }"        @click="filterTab = 'all'">全部剧集</button>
        <button class="filter-tab" :class="{ active: filterTab === 'ready' }"      @click="filterTab = 'ready'">已完成</button>
        <button class="filter-tab" :class="{ active: filterTab === 'processing' }" @click="filterTab = 'processing'">生成中</button>
      </div>
      <div class="filter-right">
        <label class="check">
          <input type="checkbox" v-model="allSelected" />
          全选当前列表
        </label>
        <span v-if="selected.length > 0" class="selected-pill">
          已选 <strong>{{ selected.length }}</strong> 集 · 约 <strong>{{ totalSize }}</strong> MB
        </span>
      </div>
    </div>

    <!-- Main work area: grid + settings panel -->
    <div class="work-area">
      <!-- Episode grid -->
      <div class="episode-grid">
        <div
          v-for="ep in filteredEpisodes"
          :key="ep.id"
          class="ep-card"
          :class="{ selected: selected.includes(ep.id) }"
        >
          <div class="ep-thumb">
            <img :src="ep.cover" :alt="ep.title" />
            <div class="ep-thumb-mask">
              <div class="ep-play">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <polygon points="5 3 19 12 5 21 5 3"/>
                </svg>
              </div>
            </div>
            <div class="ep-duration">{{ ep.duration }}</div>
            <label class="ep-check" @click.stop>
              <input
                type="checkbox"
                :checked="selected.includes(ep.id)"
                @change="toggle(ep.id)"
              />
            </label>
            <div class="ep-index">第 {{ ep.id }} 集</div>
          </div>
          <div class="ep-meta">
            <div class="ep-title">{{ ep.title }}</div>
            <div class="ep-stats">
              <span>{{ ep.size }}</span>
              <span class="dot-sep">·</span>
              <span>1080P</span>
              <span class="dot-sep">·</span>
              <span>30FPS</span>
            </div>
            <div class="ep-actions">
              <button class="btn sm" @click.stop>
                <svg class="btn-icon" viewBox="0 0 24 24"><polygon points="5 3 19 12 5 21 5 3"/></svg>
                预览
              </button>
              <button class="btn sm primary" @click.stop="toggle(ep.id); confirmDownload()">
                <svg class="btn-icon" viewBox="0 0 24 24"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
                下载
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Right: download settings panel -->
      <aside class="settings-panel">
        <div class="card card-pad">
          <div class="settings-head">
            <h2>下载设置</h2>
            <span class="chip outline-grad">默认方案</span>
          </div>

          <div class="form">
            <div class="field">
              <label class="field-label">分辨率</label>
              <select class="select" v-model="resolution">
                <option value="720p">720P · 高清</option>
                <option value="1080p">1080P · 全高清(推荐)</option>
                <option value="2k">2K · 超清</option>
                <option value="4k">4K · 影院级</option>
              </select>
            </div>

            <div class="field">
              <label class="field-label">输出格式</label>
              <select class="select" v-model="format">
                <option value="mp4">MP4 · 通用性最好</option>
                <option value="mov">MOV · Apple 生态</option>
                <option value="webm">WebM · 网页嵌入</option>
                <option value="gif">GIF · 动图预览</option>
              </select>
            </div>

            <div class="field-row">
              <div class="field">
                <label class="field-label">画质</label>
                <select class="select" v-model="quality">
                  <option value="origin">原画</option>
                  <option value="high">高码率</option>
                  <option value="medium">中码率</option>
                  <option value="low">压缩</option>
                </select>
              </div>
              <div class="field">
                <label class="field-label">帧率</label>
                <select class="select" v-model="fps">
                  <option value="24">24 FPS</option>
                  <option value="30">30 FPS</option>
                  <option value="60">60 FPS</option>
                </select>
              </div>
            </div>

            <div class="divider"><span>附加内容</span></div>

            <label class="check">
              <input type="checkbox" v-model="includeSubtitle" />
              <span>嵌入字幕文件 (.srt)</span>
            </label>
            <label class="check">
              <input type="checkbox" v-model="includeIntro" />
              <span>包含片头动画</span>
            </label>
            <label class="check">
              <input type="checkbox" v-model="includeWatermark" />
              <span>添加创作者水印</span>
            </label>
            <label class="check">
              <input type="checkbox" v-model="notify" />
              <span>下载完成后邮件通知</span>
            </label>

            <div class="divider"><span>打包方式</span></div>

            <label class="radio">
              <input type="radio" value="zip" v-model="packMode" />
              <span>合并为 ZIP 压缩包</span>
            </label>
            <label class="radio">
              <input type="radio" value="single" v-model="packMode" />
              <span>逐集独立下载</span>
            </label>

            <button class="btn primary lg" @click="confirmDownload">
              <svg class="btn-icon" viewBox="0 0 24 24"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
              确认下载
            </button>

            <p class="settings-tip">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" width="14" height="14"><circle cx="12" cy="12" r="10"/><path d="M12 8v4M12 16h.01"/></svg>
              文件将保留 7 天,过期前请及时转存
            </p>
          </div>
        </div>
      </aside>
    </div>
  </div>
</template>

<style scoped>
/* === Summary === */
.summary { padding: 24px; margin-bottom: 20px; }
.summary-main { display: flex; gap: 24px; align-items: stretch; }
.summary-cover {
  position: relative;
  width: 200px;
  flex-shrink: 0;
  border-radius: var(--r-md);
  overflow: hidden;
  box-shadow: var(--shadow-md);
}
.summary-cover img { width: 100%; height: 100%; object-fit: cover; display: block; aspect-ratio: 3/4; }
.summary-badge {
  position: absolute;
  top: 10px; left: 10px;
  padding: 3px 10px;
  border-radius: 999px;
  font-size: 11px;
  font-weight: 600;
  background: var(--grad);
  color: #fff;
  letter-spacing: .04em;
}
.summary-info { flex: 1; display: flex; flex-direction: column; gap: 10px; min-width: 0; }
.summary-eyebrow { display: flex; gap: 8px; flex-wrap: wrap; }
.summary-title { font-size: 22px; font-weight: 700; color: var(--text); margin-top: 2px; }
.summary-desc {
  font-size: 13px;
  color: var(--text-soft);
  line-height: 1.7;
  max-width: 720px;
}
.summary-meta {
  margin-top: auto;
  display: flex;
  align-items: center;
  gap: 18px;
  padding-top: 14px;
  border-top: 1px dashed var(--border);
}
.meta-item { display: flex; flex-direction: column; gap: 2px; }
.meta-num {
  font-size: 18px;
  font-weight: 700;
  color: var(--text);
  background: var(--grad-text);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}
.meta-label { font-size: 11px; color: var(--text-muted); }
.meta-divider { width: 1px; height: 28px; background: var(--border); }

/* === Filter row === */
.filter-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin: 20px 0 14px;
  flex-wrap: wrap;
}
.filter-tabs {
  display: inline-flex;
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: var(--r-md);
  padding: 4px;
  gap: 2px;
  box-shadow: var(--shadow-xs);
}
.filter-tab {
  border: 0;
  background: transparent;
  color: var(--text-soft);
  padding: 7px 16px;
  font-size: 13px;
  font-weight: 500;
  border-radius: 10px;
  cursor: pointer;
  transition: all .15s var(--ease);
}
.filter-tab:hover { color: var(--text); }
.filter-tab.active {
  background: var(--grad);
  color: #fff;
  font-weight: 600;
  box-shadow: 0 4px 10px rgba(20, 184, 166, .25);
}
.filter-right { display: flex; align-items: center; gap: 14px; }
.selected-pill {
  padding: 6px 12px;
  border-radius: 999px;
  background: var(--grad-soft);
  color: var(--primary-deep);
  font-size: 12px;
  font-weight: 500;
  border: 1px solid rgba(20, 184, 166, .25);
}
.selected-pill strong { font-weight: 700; }

/* === Work area === */
.work-area {
  display: grid;
  grid-template-columns: 1fr 320px;
  gap: 20px;
  align-items: start;
}

/* === Episode grid === */
.episode-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 16px;
}
.ep-card {
  background: var(--card);
  border-radius: var(--r-md);
  overflow: hidden;
  border: 1.5px solid var(--border-soft);
  transition: all .2s var(--ease);
  cursor: pointer;
  box-shadow: var(--shadow-xs);
}
.ep-card:hover {
  transform: translateY(-3px);
  box-shadow: var(--shadow-md);
  border-color: var(--border-strong);
}
.ep-card.selected {
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(20, 184, 166, .12), var(--shadow-md);
}
.ep-thumb {
  position: relative;
  aspect-ratio: 16/9;
  overflow: hidden;
  background: var(--card-soft);
}
.ep-thumb img {
  width: 100%; height: 100%; object-fit: cover;
  transition: transform .35s var(--ease);
}
.ep-card:hover .ep-thumb img { transform: scale(1.05); }
.ep-thumb-mask {
  position: absolute; inset: 0;
  background: linear-gradient(to top, rgba(0,0,0,.35) 0%, transparent 50%);
  display: flex; align-items: center; justify-content: center;
  opacity: 0;
  transition: opacity .2s;
}
.ep-card:hover .ep-thumb-mask { opacity: 1; }
.ep-play {
  width: 44px; height: 44px;
  border-radius: 50%;
  background: rgba(255,255,255,.92);
  color: var(--primary-deep);
  display: flex; align-items: center; justify-content: center;
  box-shadow: 0 4px 14px rgba(0,0,0,.2);
}
.ep-play svg { width: 20px; height: 20px; fill: currentColor; stroke: currentColor; }
.ep-duration {
  position: absolute;
  right: 8px; bottom: 8px;
  padding: 2px 8px;
  background: rgba(0,0,0,.65);
  color: #fff;
  font-size: 11px;
  border-radius: 6px;
  font-weight: 500;
  font-feature-settings: "tnum" 1;
}
.ep-check {
  position: absolute;
  top: 8px; left: 8px;
  width: 22px; height: 22px;
  border-radius: 6px;
  background: rgba(255,255,255,.92);
  display: flex; align-items: center; justify-content: center;
  cursor: pointer;
  transition: all .15s var(--ease);
  box-shadow: 0 2px 6px rgba(0,0,0,.12);
}
.ep-check input {
  appearance: none; -webkit-appearance: none;
  width: 14px; height: 14px;
  border: 1.5px solid var(--border-strong);
  background: transparent;
  border-radius: 3px;
  position: relative;
  cursor: pointer;
  margin: 0;
}
.ep-check input:checked {
  background: var(--grad);
  border-color: transparent;
}
.ep-check input:checked::after {
  content: "";
  position: absolute;
  left: 3px; top: 0px;
  width: 4px; height: 8px;
  border: solid #fff;
  border-width: 0 2px 2px 0;
  transform: rotate(45deg);
}
.ep-card.selected .ep-check {
  background: var(--grad);
}
.ep-card.selected .ep-check input { border-color: transparent; }
.ep-card.selected .ep-check::after {
  content: "";
  position: absolute;
  left: 6px; top: 3px;
  width: 5px; height: 9px;
  border: solid #fff;
  border-width: 0 2px 2px 0;
  transform: rotate(45deg);
}
.ep-card.selected .ep-check input { display: none; }
.ep-index {
  position: absolute;
  top: 8px; right: 8px;
  padding: 2px 10px;
  background: rgba(255,255,255,.92);
  color: var(--text);
  font-size: 11px;
  font-weight: 600;
  border-radius: 999px;
  box-shadow: 0 2px 6px rgba(0,0,0,.08);
}

.ep-meta { padding: 14px; display: flex; flex-direction: column; gap: 8px; }
.ep-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text);
  overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
}
.ep-stats {
  display: flex; align-items: center; gap: 6px;
  font-size: 12px; color: var(--text-muted);
  font-feature-settings: "tnum" 1;
}
.dot-sep { color: var(--text-faint); }
.ep-actions { display: flex; gap: 8px; margin-top: 4px; }
.ep-actions .btn { flex: 1; }

/* === Settings panel === */
.settings-panel { position: sticky; top: 84px; }
.settings-head {
  display: flex; align-items: center; justify-content: space-between;
  margin-bottom: 18px;
}
.settings-head h2 { font-size: 15px; font-weight: 600; }
.field-row { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.divider {
  display: flex; align-items: center; gap: 10px;
  margin: 4px 0 2px;
  font-size: 11px;
  color: var(--text-muted);
  font-weight: 500;
  letter-spacing: .04em;
}
.divider::before, .divider::after {
  content: ""; flex: 1; height: 1px; background: var(--border);
}
.settings-tip {
  display: flex; align-items: center; gap: 6px;
  font-size: 11.5px;
  color: var(--text-muted);
  margin: 0;
}
.settings-tip svg { flex-shrink: 0; color: var(--primary); }

/* === Responsive === */
@media (max-width: 1100px) {
  .work-area { grid-template-columns: 1fr; }
  .settings-panel { position: static; }
}
@media (max-width: 720px) {
  .summary-main { flex-direction: column; }
  .summary-cover { width: 100%; aspect-ratio: 16/9; }
  .summary-title { font-size: 18px; }
  .summary-meta { flex-wrap: wrap; gap: 12px; }
  .meta-divider { display: none; }
}
</style>
