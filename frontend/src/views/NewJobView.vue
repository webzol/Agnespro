<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { useRouter } from "vue-router";
import { jobsApi, stylesApi, adminApi, scriptApi } from "@/api/endpoints";
import { useToast } from "@/composables/toast";
import type { StyleLibrary, ModelInfo, VisualStyle, Genre, AspectRatio } from "@/types";

const router = useRouter();
const toast = useToast();

// Form state
const mode = ref<"manual"|"ai">("manual");
const title = ref("");
const script = ref("");
const idea = ref("");
const aiTitle = ref("");
const aiGenre = ref("urban_power");
const aiLength = ref("medium");

// Style library
const LIB = ref<StyleLibrary>({ visual_styles: [], genres: [], aspect_ratios: [], episode_counts: [] });
const MODELS = ref<ModelInfo[]>([]);
const LIB_LOADED = ref(false);

async function loadLibrary() {
  try {
    LIB.value = await stylesApi.library();
    LIB_LOADED.value = true;
  } catch (e) { toast.err("加载风格库失败: " + (e as Error).message); }
}
async function loadModels() {
  try {
    const r = await adminApi.listModels();
    MODELS.value = r.models || [];
  } catch (e) { MODELS.value = []; }
}

// Style picker state
const visualStyle = ref("cinematic");
const visualStyleName = ref("");

// 4 chip-btn states
const aspectRatio = ref("9:16");
const episodeCount = ref(10);
const genreId = ref("urban_power");
const genreName = ref("都市爽文风");

// Model states
const scriptModel = ref("agnes-2.5-flash");
const imageModel = ref("agnes-image-2.1-flash");
const videoModel = ref("agnes-video-v2.0");

const chatModels = computed(() => MODELS.value.filter(m => m.type === "chat"));
const imgModels  = computed(() => MODELS.value.filter(m => m.type === "image"));
const vidModels  = computed(() => MODELS.value.filter(m => m.type === "video"));

// Parse preview
const parsed = ref<{ index: number; title: string; body: string }[]>([]);
const parsing = ref(false);
const creating = ref(false);

onMounted(async () => {
  await Promise.all([loadLibrary(), loadModels()]);
  // Default genre name
  const g = LIB.value.genres.find(x => x.id === genreId.value);
  if (g) genreName.value = g.name;
});

// === Style library modal ===
const LIB_MODAL_OPEN = ref(false);
const LIB_FILTER = ref<"all"|"2d"|"3d"|"真人">("all");
const LIB_SELECTED = ref<{ id: string; name: string } | null>(null);

const filteredStyles = computed(() => {
  const all = LIB.value.visual_styles;
  if (LIB_FILTER.value === "all") return all;
  return all.filter(s => s.category === LIB_FILTER.value);
});

function openLib() { LIB_MODAL_OPEN.value = true; LIB_FILTER.value = "all"; LIB_SELECTED.value = null; }
function closeLib() { LIB_MODAL_OPEN.value = false; LIB_SELECTED.value = null; }
function applyLib() {
  if (LIB_SELECTED.value) {
    visualStyle.value = LIB_SELECTED.value.id;
    visualStyleName.value = LIB_SELECTED.value.name;
    closeLib();
  }
}

// === AI script generation ===
const aiGenerating = ref(false);
const aiResult = ref("");
const aiResultTitle = ref("");

async function generateScript() {
  if (!idea.value.trim()) { toast.err("请先输入你的想法"); return; }
  aiGenerating.value = true;
  aiResult.value = "";
  try {
    const r = await scriptApi.generate({
      title: aiTitle.value || idea.value.slice(0, 20),
      style: visualStyleName.value,
      idea: idea.value,
      genre: aiGenre.value,
      length: aiLength.value,
      lang: "zh",
      visual_style: visualStyle.value,
      visual_style_name: visualStyleName.value,
      aspect_ratio: aspectRatio.value,
      episode_count: episodeCount.value,
      script_model: scriptModel.value,
    });
    aiResult.value = r.script;
    aiResultTitle.value = r.title || aiTitle.value || "未命名剧本";
    toast.ok("剧本生成完成");
  } catch (e) {
    toast.err("生成失败: " + (e as Error).message);
  } finally {
    aiGenerating.value = false;
  }
}

function applyGenerated() {
  script.value = aiResult.value;
  title.value = aiResultTitle.value;
  aiResult.value = "";
  mode.value = "manual";
  toast.ok("已应用到剧本");
}

function discardGenerated() {
  aiResult.value = "";
}

// === Parse-only ===
async function parseOnly() {
  if (!script.value.trim()) { toast.err("请先粘贴剧本"); return; }
  parsing.value = true;
  try {
    const r = await jobsApi.create({ title: title.value || "未命名", script: script.value, style: visualStyleName.value });
    parsed.value = r.episodes || [];
    toast.ok("检测到 " + parsed.value.length + " 集");
  } catch (e) { toast.err((e as Error).message); }
  finally { parsing.value = false; }
}

// === Create job ===
async function createJob() {
  if (!script.value.trim()) { toast.err("请先粘贴剧本"); return; }
  creating.value = true;
  try {
    const r = await jobsApi.create({
      title: title.value || "Untitled",
      script: script.value,
      style: visualStyleName.value,
      visual_style: visualStyle.value,
      visual_style_name: visualStyleName.value,
      aspect_ratio: aspectRatio.value,
      episode_count: episodeCount.value,
      genre_id: genreId.value,
      genre_name: genreName.value,
      script_model: scriptModel.value,
      image_model: imageModel.value,
      video_model: videoModel.value,
    });
    toast.ok("任务已创建");
    script.value = ""; title.value = ""; parsed.value = [];
    router.push({ name: "detail", params: { id: r.id } });
  } catch (e) { toast.err((e as Error).message); }
  finally { creating.value = false; }
}
</script>

<template>
  <section class="card glass">
    <div class="mode-tabs">
      <button :class="['mode-tab', mode === 'manual' && 'active']" @click="mode = 'manual'">
        <span class="ico">📝</span>上传剧本
      </button>
      <button :class="['mode-tab', mode === 'ai' && 'active']" @click="mode = 'ai'">
        <span class="ico">✨</span>AI 生成剧本
      </button>
    </div>

    <!-- Manual mode -->
    <div v-show="mode === 'manual'" class="mode-panel">
      <div class="form">
        <label>
          <span>标题</span>
          <input v-model="title" type="text" placeholder="例:咖啡店的奇妙夜晚">
        </label>
        <label>
          <span>视觉风格</span>
          <div v-if="LIB_LOADED && LIB.visual_styles.length === 0" class="muted small">未加载到风格</div>
          <div v-else class="style-picker" id="job-style-picker">
            <div v-if="!LIB_LOADED" class="style-picker-loading">加载风格中...</div>
            <div v-else class="style-grid">
              <div v-for="s in LIB.visual_styles" :key="s.id"
                   :class="['style-card', visualStyle === s.id && 'is-selected']"
                   :data-style-id="s.id"
                   @click="visualStyle = s.id; visualStyleName = s.name">
                <img :src="s.preview" :alt="s.name" loading="lazy">
                <div class="style-card-label">{{ s.name }}</div>
              </div>
            </div>
          </div>
          <input id="job-style" type="hidden" :value="visualStyle">
          <span class="hint">点击图片选择风格,影响图像生成与视频的视觉表现。</span>
        </label>
        <label>
          <span>剧本正文</span>
          <textarea v-model="script" rows="12" placeholder="第1集 雨夜...&#10;咖啡店里,林夏...&#10;&#10;第2集 重逢..."/>
          <span class="hint">剧集标记支持:第N集 / Episode N / EP N / Scene N。检测不到标记时,整段作为单集处理。</span>
        </label>

        <div class="chip-grid">
          <button class="chip-btn" id="btn-style-library" type="button" @click="openLib">
            <span class="chip-ico">✨</span>
            <span class="chip-lbl">风格库</span>
            <span class="chip-val" id="chip-style-library">{{ visualStyleName || visualStyle }}</span>
            <span class="chip-caret">▾</span>
          </button>
          <label class="chip-btn chip-select-wrap">
            <span class="chip-ico">▭</span>
            <span class="chip-lbl">比例</span>
            <select v-model="aspectRatio" class="chip-native">
              <option v-for="a in LIB.aspect_ratios" :key="a.id" :value="a.id">{{ a.label }}</option>
            </select>
            <span class="chip-caret">▾</span>
          </label>
          <label class="chip-btn chip-select-wrap">
            <span class="chip-ico">▶</span>
            <span class="chip-lbl">剧集</span>
            <select v-model.number="episodeCount" class="chip-native">
              <option v-for="ec in LIB.episode_counts" :key="ec.value" :value="ec.value">{{ ec.label }}</option>
            </select>
            <span class="chip-caret">▾</span>
          </label>
        </div>

        <div class="model-row">
          <label class="mini-label">AI 剧本 (chat)
            <select v-model="scriptModel" class="mini-select">
              <option v-for="m in chatModels" :key="m.id" :value="m.id">{{ m.id }}</option>
            </select>
          </label>
          <label class="mini-label">AI 绘图 (image)
            <select v-model="imageModel" class="mini-select">
              <option v-for="m in imgModels" :key="m.id" :value="m.id">{{ m.id }}</option>
            </select>
          </label>
          <label class="mini-label">AI 视频 (video)
            <select v-model="videoModel" class="mini-select">
              <option v-for="m in vidModels" :key="m.id" :value="m.id">{{ m.id }}</option>
            </select>
          </label>
        </div>

        <div v-if="parsed.length" class="parse-preview">
          <h3>检测到 {{ parsed.length }} 集</h3>
          <div v-for="ep in parsed" :key="ep.index" class="ep">
            <strong>第{{ ep.index }}集: {{ ep.title || "(无标题)" }}</strong>
            <div class="muted small">{{ ep.body.slice(0, 80) }}{{ ep.body.length > 80 ? "..." : "" }}</div>
          </div>
        </div>

        <div class="form-actions">
          <button class="btn glass" :disabled="parsing" @click="parseOnly">
            {{ parsing ? "解析中..." : "仅解析剧集" }}
          </button>
          <button class="btn primary" :disabled="creating" @click="createJob">
            {{ creating ? "创建中..." : "创建任务" }}
          </button>
        </div>
      </div>
    </div>

    <!-- AI mode -->
    <div v-show="mode === 'ai'" class="mode-panel">
      <div class="form">
        <label>
          <span>标题(可选,留空则用你的想法前 20 字)</span>
          <input v-model="aiTitle" type="text" placeholder="例:雨夜咖啡店">
        </label>
        <label>
          <span>题材</span>
          <select v-model="aiGenre">
            <option v-for="g in LIB.genres" :key="g.id" :value="g.id" @change="genreName = g.name">{{ g.name }}</option>
          </select>
        </label>
        <label>
          <span>视觉风格</span>
          <input v-model="visualStyleName" type="text" placeholder="选填,如:韩漫厚涂">
        </label>
        <label>
          <span>长度</span>
          <select v-model="aiLength">
            <option value="short">短 · 约 3 场戏 / 400-600 字</option>
            <option value="medium" selected>中 · 约 5 场戏 / 800-1200 字</option>
            <option value="long">长 · 约 8 场戏 / 1500-2000 字</option>
          </select>
        </label>
        <label>
          <span>你的想法</span>
          <textarea v-model="idea" rows="6" placeholder="例:雨夜,一位独自在咖啡店加班的女孩,推门进来一位五年前不告而别的初恋男友..."/>
          <span class="hint">越具体越好:人物、场景、核心冲突、情绪氛围。</span>
        </label>
        <div class="form-actions">
          <button class="btn primary" :disabled="aiGenerating" @click="generateScript">
            {{ aiGenerating ? "生成中..." : "✨ 生成剧本" }}
          </button>
        </div>
        <div v-if="aiResult" class="ai-result">
          <h3>{{ aiResultTitle }}</h3>
          <pre>{{ aiResult }}</pre>
          <div class="form-actions">
            <button class="btn glass" @click="discardGenerated">丢弃</button>
            <button class="btn primary" @click="applyGenerated">应用到剧本并继续 →</button>
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- Style library modal -->
  <Teleport to="body">
    <div v-if="LIB_MODAL_OPEN" class="modal">
      <div class="modal-mask" @click="closeLib"></div>
      <div class="modal-card glass">
        <div class="modal-head">
          <h3>风格库</h3>
          <button class="modal-close" @click="closeLib">✕</button>
        </div>
        <div class="modal-tabs">
          <button v-for="c in ['all','2d','3d','真人']" :key="c"
                  :class="['modal-tab', LIB_FILTER === c && 'is-active']"
                  @click="LIB_FILTER = c as any">{{ c === 'all' ? '全部' : c }}</button>
        </div>
        <div class="modal-body">
          <div v-if="filteredStyles.length === 0" class="muted small" style="padding: 40px; text-align: center;">未加载到风格</div>
          <div v-else class="lib-grid">
            <div v-for="s in filteredStyles" :key="s.id"
                 :class="['lib-card', LIB_SELECTED?.id === s.id && 'is-selected']"
                 @click="LIB_SELECTED = { id: s.id, name: s.name }">
              <div class="lib-img-wrap"><img :src="s.preview" :alt="s.name" loading="lazy"></div>
              <div class="lib-name">{{ s.name }}</div>
              <div class="lib-desc">{{ s.desc }}</div>
            </div>
          </div>
        </div>
        <div class="modal-foot">
          <span class="muted">{{ LIB_SELECTED ? LIB_SELECTED.name : '尚未选择' }}</span>
          <button class="btn primary" :disabled="!LIB_SELECTED" @click="applyLib">适用</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.mode-tabs { display: inline-flex; gap: 4px; padding: 5px; border-radius: 999px; margin-bottom: 20px; }
.mode-tab { border: 0; background: transparent; color: var(--text-soft); padding: 9px 18px; border-radius: 999px; cursor: pointer; font-size: 13px; }
.mode-tab:hover { color: var(--text); }
.mode-tab.active { background: var(--grad-soft); color: var(--text); }
.mode-tab .ico { margin-right: 6px; }

.style-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(110px, 1fr)); gap: 10px; }
.style-card { border-radius: 10px; overflow: hidden; border: 2px solid transparent; cursor: pointer; transition: all .15s; background: rgba(255,255,255,.03); }
.style-card:hover { transform: translateY(-2px); border-color: rgba(99,102,241,.4); }
.style-card.is-selected { border-color: var(--primary); box-shadow: 0 0 0 2px rgba(99,102,241,.3); }
.style-card img { width: 100%; aspect-ratio: 4/3; object-fit: cover; display: block; }
.style-card-label { padding: 6px 8px; font-size: 12px; text-align: center; }
.style-picker-loading { padding: 24px; text-align: center; color: var(--text-muted); font-size: 13px; }

.chip-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px; margin-top: 4px; }
.chip-btn { display: flex; align-items: center; gap: 6px; padding: 10px 12px; border-radius: 10px; background: rgba(255,255,255,.04); border: 1px solid var(--border); color: var(--text); cursor: pointer; font-size: 13px; min-height: 42px; transition: all .15s; }
.chip-btn:hover { background: rgba(255,255,255,.08); border-color: var(--border-strong); }
.chip-ico { opacity: .8; }
.chip-lbl { opacity: .65; font-size: 12px; }
.chip-val { margin-left: auto; background: rgba(99,102,241,.18); padding: 2px 8px; border-radius: 6px; font-size: 12px; }
.chip-caret { opacity: .5; }
.chip-select-wrap { position: relative; padding: 0; }
.chip-select-wrap .chip-ico, .chip-select-wrap .chip-lbl, .chip-select-wrap .chip-caret { position: absolute; pointer-events: none; }
.chip-select-wrap .chip-ico { left: 12px; top: 50%; transform: translateY(-50%); }
.chip-select-wrap .chip-lbl { left: 36px; top: 50%; transform: translateY(-50%); }
.chip-select-wrap .chip-caret { right: 12px; top: 50%; transform: translateY(-50%); }
.chip-native { appearance: none; background: transparent; border: 0; color: transparent; padding: 10px 12px 10px 70px; width: 100%; height: 100%; cursor: pointer; font-family: inherit; font-size: 13px; }
.chip-native option { color: #1a1a1a; background: #fff; }

.model-row { display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px; margin-top: 4px; }
.mini-label { display: flex; flex-direction: column; gap: 6px; font-size: 12px; opacity: .8; }
.mini-select { padding: 8px 10px; border-radius: 8px; background: rgba(255,255,255,.04); border: 1px solid var(--border); color: var(--text); font-size: 13px; }

.parse-preview { background: rgba(255,255,255,.03); border: 1px solid var(--border); border-radius: 12px; padding: 16px; }
.parse-preview h3 { font-size: 14px; margin-bottom: 10px; }
.parse-preview .ep { padding: 8px 10px; background: rgba(255,255,255,.03); border-radius: 8px; margin-bottom: 6px; }

.ai-result { background: rgba(255,255,255,.03); border: 1px solid var(--border); border-radius: 12px; padding: 16px; max-height: 400px; overflow-y: auto; }
.ai-result h3 { font-size: 14px; margin-bottom: 8px; }
.ai-result pre { white-space: pre-wrap; font-size: 12px; line-height: 1.6; font-family: inherit; margin: 0; }

/* === Style Library Modal === */
.modal { position: fixed; inset: 0; z-index: 100; display: flex; align-items: center; justify-content: center; padding: 24px; }
.modal-mask { position: absolute; inset: 0; background: rgba(0,0,0,.55); backdrop-filter: blur(6px); }
.modal-card { position: relative; width: min(880px, 100%); max-height: calc(100vh - 48px); display: flex; flex-direction: column; border-radius: 18px; overflow: hidden; }
.modal-head { display: flex; align-items: center; justify-content: space-between; padding: 18px 22px 12px; border-bottom: 1px solid var(--border); }
.modal-close { background: transparent; border: 0; color: #aaa; font-size: 18px; cursor: pointer; padding: 4px 10px; border-radius: 8px; }
.modal-close:hover { background: rgba(255,255,255,.08); color: #fff; }
.modal-tabs { display: flex; gap: 8px; padding: 12px 22px 8px; border-bottom: 1px solid var(--border); }
.modal-tab { background: transparent; border: 1px solid var(--border); color: #ccc; padding: 6px 16px; border-radius: 999px; cursor: pointer; font-size: 13px; }
.modal-tab.is-active { background: rgba(99,102,241,.25); border-color: rgba(99,102,241,.6); color: #fff; }
.modal-body { flex: 1; overflow-y: auto; padding: 18px 22px; }
.lib-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(140px, 1fr)); gap: 14px; }
.lib-card { position: relative; display: flex; flex-direction: column; border-radius: 12px; overflow: hidden; background: rgba(255,255,255,.04); border: 1px solid var(--border); cursor: pointer; transition: all .15s; }
.lib-card:hover { border-color: rgba(99,102,241,.4); transform: translateY(-2px); }
.lib-card.is-selected { border-color: rgba(99,102,241,.85); box-shadow: 0 0 0 2px rgba(99,102,241,.4); }
.lib-img-wrap { aspect-ratio: 4/3; background: linear-gradient(135deg, rgba(99,102,241,.18), rgba(168,85,247,.18)); display: flex; align-items: center; justify-content: center; overflow: hidden; }
.lib-img-wrap img { width: 100%; height: 100%; object-fit: cover; }
.lib-name { padding: 8px 10px 2px; font-size: 13px; font-weight: 600; }
.lib-desc { padding: 0 10px 10px; font-size: 11px; color: #aaa; line-height: 1.4; }
.modal-foot { display: flex; align-items: center; justify-content: space-between; padding: 14px 22px; border-top: 1px solid var(--border); }

@media (max-width: 720px) {
  .chip-grid, .model-row { grid-template-columns: 1fr; }
  .lib-grid { grid-template-columns: repeat(auto-fill, minmax(110px, 1fr)); }
}
</style>
