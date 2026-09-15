<script setup lang="ts">
import { docsApi } from "@/api/endpoints";
import { ref, onMounted } from "vue";

const md = ref("");
const loading = ref(true);

onMounted(async () => {
  try {
    md.value = await docsApi.md();
  } catch (e) {
    md.value = "# 文档加载失败\n\n请确认后端 `/api/docs` 端点正常工作。";
  } finally { loading.value = false; }
});

function render(text: string): string {
  // 简单 markdown -> HTML 渲染 (标题/列表/代码块/加粗/链接)
  if (!text) return "";
  const esc = (s: string) => s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
  const lines = text.split(/\r?\n/);
  const out: string[] = [];
  let inCode = false, inList = false;
  for (let line of lines) {
    if (line.startsWith("```")) {
      if (inCode) { out.push("</code></pre>"); inCode = false; }
      else { out.push("<pre><code>"); inCode = true; }
      continue;
    }
    if (inCode) { out.push(esc(line)); continue; }
    if (/^#{1,6}\s/.test(line)) {
      const m = line.match(/^(#{1,6})\s+(.*)$/);
      if (m) { if (inList) { out.push("</ul>"); inList = false; } out.push(`<h${m[1].length}>${inline(m[2])}</h${m[1].length}>`); continue; }
    }
    if (/^[-*]\s/.test(line)) {
      if (!inList) { out.push("<ul>"); inList = true; }
      out.push("<li>" + inline(line.replace(/^[-*]\s+/, "")) + "</li>");
      continue;
    }
    if (line.trim() === "") { if (inList) { out.push("</ul>"); inList = false; } out.push(""); continue; }
    out.push("<p>" + inline(line) + "</p>");
  }
  if (inList) out.push("</ul>");
  if (inCode) out.push("</code></pre>");
  return out.join("\n");

  function inline(s: string): string {
    return esc(s)
      .replace(/\*\*(.+?)\*\*/g, "<strong>$1</strong>")
      .replace(/`([^`]+)`/g, "<code>$1</code>")
      .replace(/\[([^\]]+)\]\(([^)]+)\)/g, `<a href="$2" target="_blank">$1</a>`);
  }
}
</script>

<template>
  <section class="card glass">
    <h2>使用文档</h2>
    <p class="muted">本系统核心操作手册:粘贴剧本、生成短剧、查看任务、后台配置。</p>
    <div v-if="loading" class="muted small" style="padding: 20px;">加载中…</div>
    <div v-else class="docs-body" v-html="render(md)"></div>
  </section>
</template>

<style scoped>
.card { padding: 28px; }
.docs-body { line-height: 1.75; font-size: 14px; color: var(--text-soft); }
.docs-body :deep(h1) { font-size: 24px; margin: 24px 0 12px; border-bottom: 1px solid var(--border); padding-bottom: 10px; color: var(--text); font-weight: 700; }
.docs-body :deep(h2) { font-size: 20px; margin: 22px 0 10px; color: var(--text); font-weight: 600; }
.docs-body :deep(h3) { font-size: 16px; margin: 16px 0 8px; color: var(--text); font-weight: 600; }
.docs-body :deep(h4) { font-size: 14px; margin: 14px 0 6px; color: var(--text); font-weight: 600; }
.docs-body :deep(p) { margin: 10px 0; color: var(--text-soft); }
.docs-body :deep(ul) { margin: 10px 0; padding-left: 24px; color: var(--text-soft); }
.docs-body :deep(li) { margin: 6px 0; }
.docs-body :deep(code) { background: var(--card-soft); border: 1px solid var(--border-soft); padding: 1px 7px; border-radius: 5px; font-size: 12px; color: var(--primary-deep); font-family: ui-monospace, monospace; }
.docs-body :deep(pre) { background: var(--card-soft); border: 1px solid var(--border-soft); padding: 14px 16px; border-radius: var(--r-sm); overflow-x: auto; margin: 12px 0; }
.docs-body :deep(pre code) { background: transparent; border: 0; padding: 0; color: var(--text); }
.docs-body :deep(strong) { color: var(--text); font-weight: 600; }
.docs-body :deep(a) { color: var(--primary-deep); text-decoration: none; border-bottom: 1px dashed rgba(20, 184, 166, .4); }
.docs-body :deep(a:hover) { color: var(--primary-deep-2); border-bottom-style: solid; }
</style>
