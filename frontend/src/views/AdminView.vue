<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { adminApi } from "@/api/endpoints";
import { useToast } from "@/composables/toast";
import type { AdminConfig, ModelInfo } from "@/types";

const toast = useToast();
const CFG = ref<AdminConfig | null>(null);
const apiKeyInput = ref("");
const showKey = ref(false);
const testResult = ref<{ ok: boolean; message?: string; duration_ms?: number; base_url?: string } | null>(null);
const saving = ref(false);
const savingModels = ref(false);

async function load() {
  try {
    CFG.value = await adminApi.getConfig();
  } catch (e) {
    toast.err("加载后台配置失败: " + (e as Error).message);
  }
}
onMounted(load);

async function saveKey() {
  const k = apiKeyInput.value.trim();
  if (!k) { toast.err("请先输入 API Key"); return; }
  if (!k.startsWith("sk-")) { toast.err("Key 格式应以 sk- 开头"); return; }
  saving.value = true;
  try {
    await adminApi.updateConfig({ api_key: k });
    apiKeyInput.value = "";
    toast.ok("Key 已保存(加密)");
    await load();
  } catch (e) { toast.err((e as Error).message); }
  finally { saving.value = false; }
}

async function clearKey() {
  if (!confirm("确定要清除已保存的 API Key 吗?")) return;
  try {
    await adminApi.updateConfig({ clear_api_key: true });
    toast.ok("Key 已清除");
    await load();
  } catch (e) { toast.err((e as Error).message); }
}

async function testKey() {
  const k = apiKeyInput.value.trim();
  testResult.value = null;
  try {
    testResult.value = await adminApi.testKey(k || undefined);
  } catch (e) { toast.err((e as Error).message); }
}

async function saveRoute() {
  if (!CFG.value) return;
  try {
    const r = await adminApi.updateConfig({ route: CFG.value.api_route });
    if (r.restart_required) toast.info("路由已保存,需重启服务生效");
    else toast.ok("路由已保存");
    await load();
  } catch (e) { toast.err((e as Error).message); }
}

async function saveModels() {
  if (!CFG.value) return;
  savingModels.value = true;
  try {
    await adminApi.updateConfig({
      script_model: CFG.value.script_model,
      image_model: CFG.value.image_model,
      video_model: CFG.value.video_model,
    });
    toast.ok("默认模型已保存");
    await load();
  } catch (e) { toast.err((e as Error).message); }
  finally { savingModels.value = false; }
}

const chatModels = computed<ModelInfo[]>(() => (CFG.value?.available_models || []).filter(m => m.type === "chat"));
const imgModels  = computed<ModelInfo[]>(() => (CFG.value?.available_models || []).filter(m => m.type === "image"));
const vidModels  = computed<ModelInfo[]>(() => (CFG.value?.available_models || []).filter(m => m.type === "video"));

function pickChip(modelId: string, type: string) {
  if (!CFG.value) return;
  if (type === "chat") CFG.value.script_model = modelId;
  else if (type === "image") CFG.value.image_model = modelId;
  else if (type === "video") CFG.value.video_model = modelId;
}
</script>

<template>
  <section class="card glass">
    <div class="card-head">
      <h2>⚙ 后台配置</h2>
      <span class="muted small">{{ CFG ? "已加载 · " + new Date(CFG.fetched_at).toLocaleTimeString() : "加载中…" }}</span>
    </div>
    <p class="muted">在这里集中管理 Agnes AI API Key、API 路由、默认 AI 模型。模型改动立即生效(route 切换除外)。</p>

    <div v-if="CFG">
      <!-- API Key -->
      <div class="admin-section">
        <h3 class="admin-section-title">1. API Key</h3>
        <p class="muted small">密钥经 AES-256-GCM 加密后存到 <code>data/.apikey</code>,不会以明文出现在任何日志或代码中。</p>
        <div class="form">
          <label>
            <span>当前状态</span>
            <span :class="['key-status', CFG.api_key_set ? 'set' : 'unset']">
              {{ CFG.api_key_set ? "✓ 已设置(加密存储)" : "✗ 未设置" }}
            </span>
          </label>
          <label>
            <span>设置新的 API Key</span>
            <div class="row">
              <input v-model="apiKeyInput" :type="showKey ? 'text' : 'password'" placeholder="sk-..." autocomplete="off">
              <button class="btn glass" type="button" @click="showKey = !showKey">{{ showKey ? "隐藏" : "显示" }}</button>
            </div>
            <span class="hint">留空表示不修改。Key 仅在保存时使用,不显示在 UI。</span>
          </label>
          <div class="form-actions">
            <button class="btn glass" :disabled="saving" @click="testKey">测试连接</button>
            <button class="btn ghost" @click="clearKey">清除已保存的 Key</button>
            <button class="btn primary" :disabled="saving" @click="saveKey">保存 Key</button>
          </div>
          <div v-if="testResult" :class="['test-result', testResult.ok ? 'ok' : 'err']">
            {{ testResult.ok ? "✓ " + (testResult.message || "ok") + " (" + testResult.duration_ms + "ms) · " + testResult.base_url : "✗ " + testResult.message }}
          </div>
        </div>
      </div>

      <!-- API Route -->
      <div class="admin-section">
        <h3 class="admin-section-title">2. API 路由</h3>
        <p class="muted small">选择 Agnes AI 网关地址。切换后需要重启服务才能生效。</p>
        <div class="route-options">
          <label v-for="opt in [
            { value: 'international', label: '国际版', url: 'https://apihub.agnes-ai.com/v1', def: true },
            { value: 'international-alt', label: '国际备用', url: 'https://apihub.agnes-ai.cn/v1', def: false },
            { value: 'china', label: '国内版', url: 'https://api.agnes-ai.cn/v1', def: false },
          ]" :key="opt.value" :class="['route-option', 'glass', CFG.api_route === opt.value && 'is-selected']">
            <input type="radio" name="admin-route" :value="opt.value" :checked="CFG.api_route === opt.value" @change="saveRoute">
            <div class="route-option-body">
              <div class="route-option-name">
                {{ opt.label }}
                <span v-if="opt.def" class="badge">默认</span>
                <span v-if="CFG.api_route === opt.value" class="badge active">当前</span>
              </div>
              <div class="route-option-url">{{ opt.url }}</div>
            </div>
          </label>
        </div>
      </div>

      <!-- Default Models -->
      <div class="admin-section">
        <h3 class="admin-section-title">3. 默认 AI 模型</h3>
        <p class="muted small">新建任务时用户可以临时覆盖这里的设置。</p>
        <div class="model-row">
          <label class="mini-label">AI 剧本 (chat)
            <select v-model="CFG.script_model" class="mini-select">
              <option v-for="m in chatModels" :key="m.id" :value="m.id">{{ m.id }}</option>
            </select>
          </label>
          <label class="mini-label">AI 绘图 (image)
            <select v-model="CFG.image_model" class="mini-select">
              <option v-for="m in imgModels" :key="m.id" :value="m.id">{{ m.id }}</option>
            </select>
          </label>
          <label class="mini-label">AI 视频 (video)
            <select v-model="CFG.video_model" class="mini-select">
              <option v-for="m in vidModels" :key="m.id" :value="m.id">{{ m.id }}</option>
            </select>
          </label>
        </div>
        <div class="form-actions">
          <button class="btn primary" :disabled="savingModels" @click="saveModels">{{ savingModels ? "保存中…" : "保存默认模型" }}</button>
        </div>
      </div>

      <!-- Available Models -->
      <div class="admin-section">
        <h3 class="admin-section-title">4. 可用模型列表</h3>
        <p class="muted small">从当前 API Key 拉取的 Agnes 网关可用模型(只读,按类型分组)。点击 chip 自动填入对应 select。</p>
        <div v-if="CFG.models_error" class="err-banner">{{ CFG.models_error }}</div>
        <div v-else class="models-list">
          <div class="model-group">
            <div class="model-group-label">AI 剧本 (chat) · {{ chatModels.length }} 个</div>
            <div class="model-group-items">
              <span v-for="m in chatModels" :key="m.id" class="model-chip" @click="pickChip(m.id, 'chat')">{{ m.id }}</span>
            </div>
          </div>
          <div class="model-group">
            <div class="model-group-label">AI 绘图 (image) · {{ imgModels.length }} 个</div>
            <div class="model-group-items">
              <span v-for="m in imgModels" :key="m.id" class="model-chip" @click="pickChip(m.id, 'image')">{{ m.id }}</span>
            </div>
          </div>
          <div class="model-group">
            <div class="model-group-label">AI 视频 (video) · {{ vidModels.length }} 个</div>
            <div class="model-group-items">
              <span v-for="m in vidModels" :key="m.id" class="model-chip" @click="pickChip(m.id, 'video')">{{ m.id }}</span>
            </div>
          </div>
        </div>
        <div class="form-actions">
          <button class="btn glass" @click="load">🔄 刷新模型列表</button>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.card-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 14px; flex-wrap: wrap; gap: 8px; }
.card-head h2 { font-size: 18px; }
.admin-section { margin: 18px 0; padding: 16px; border-radius: 12px; background: rgba(255,255,255,.03); border: 1px solid var(--border); }
.admin-section h3 { font-size: 14px; margin-bottom: 8px; }
.admin-section .muted.small { font-size: 12px; margin-bottom: 10px; }
.admin-section code { background: rgba(255,255,255,.06); padding: 1px 5px; border-radius: 4px; font-size: 11px; }

.form { display: flex; flex-direction: column; gap: 12px; }
.form label { display: flex; flex-direction: column; gap: 6px; }
.form label > span { font-size: 12px; color: var(--text-soft); font-weight: 500; }
.form input, .form select {
  background: rgba(255,255,255,.04);
  border: 1px solid var(--border);
  color: var(--text);
  border-radius: 10px;
  padding: 10px 12px;
  font-size: 13px;
  outline: none;
}
.form input:focus, .form select:focus { border-color: var(--primary); }
.row { display: flex; gap: 8px; align-items: center; }
.row input { flex: 1; }
.form-actions { display: flex; gap: 8px; justify-content: flex-end; margin-top: 6px; flex-wrap: wrap; }

.key-status { display: inline-block; padding: 4px 12px; border-radius: 999px; font-size: 13px; font-weight: 500; }
.key-status.set { background: rgba(74,222,128,.18); color: #86efac; border: 1px solid rgba(74,222,128,.35); }
.key-status.unset { background: rgba(248,113,113,.18); color: #fca5a5; border: 1px solid rgba(248,113,113,.35); }
.test-result { padding: 10px 12px; border-radius: 10px; font-size: 12px; font-family: ui-monospace, monospace; word-break: break-all; }
.test-result.ok { background: rgba(74,222,128,.18); color: #86efac; border: 1px solid rgba(74,222,128,.35); }
.test-result.err { background: rgba(248,113,113,.18); color: #fca5a5; border: 1px solid rgba(248,113,113,.35); }

.route-options { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 10px; margin-top: 10px; }
.route-option { display: flex; align-items: flex-start; gap: 10px; padding: 10px 12px; border-radius: 12px; cursor: pointer; background: rgba(255,255,255,.04); border: 1px solid var(--border); }
.route-option:hover { border-color: rgba(99,102,241,.4); }
.route-option.is-selected { border-color: rgba(99,102,241,.7); background: rgba(99,102,241,.10); }
.route-option input[type="radio"] { margin-top: 4px; accent-color: #6366f1; cursor: pointer; }
.route-option-name { font-size: 13px; font-weight: 600; margin-bottom: 4px; }
.route-option-url { font-size: 11px; color: var(--text-muted); font-family: ui-monospace, monospace; word-break: break-all; }
.badge { display: inline-block; margin-left: 6px; padding: 1px 7px; font-size: 10px; background: rgba(99,102,241,.25); border: 1px solid rgba(99,102,241,.5); border-radius: 999px; color: #c7d2fe; }
.badge.active { background: rgba(74,222,128,.25); border-color: rgba(74,222,128,.5); color: #86efac; }

.model-row { display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px; }
.mini-label { display: flex; flex-direction: column; gap: 6px; font-size: 12px; opacity: .8; }
.mini-select { padding: 8px 10px; border-radius: 8px; background: rgba(255,255,255,.04); border: 1px solid var(--border); color: var(--text); font-size: 13px; }

.models-list { display: flex; flex-direction: column; gap: 12px; max-height: 360px; overflow-y: auto; }
.model-group-label { font-size: 11px; opacity: .7; margin-bottom: 6px; letter-spacing: .04em; text-transform: uppercase; }
.model-group-items { display: flex; flex-wrap: wrap; gap: 6px; }
.model-chip { display: inline-block; padding: 5px 11px; background: rgba(255,255,255,.06); border: 1px solid var(--border); border-radius: 999px; font-size: 12px; cursor: pointer; font-family: ui-monospace, monospace; }
.model-chip:hover { background: rgba(99,102,241,.20); border-color: rgba(99,102,241,.55); }
.err-banner { padding: 12px; background: rgba(248,113,113,.12); border: 1px solid rgba(248,113,113,.3); border-radius: 10px; color: #fca5a5; font-size: 13px; margin-bottom: 14px; }

.card { padding: 24px; }
@media (max-width: 720px) {
  .model-row { grid-template-columns: 1fr; }
  .route-options { grid-template-columns: 1fr; }
}
</style>
