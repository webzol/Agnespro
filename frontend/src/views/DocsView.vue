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
.card { padding: 24px; }
.docs-body { line-height: 1.7; font-size: 14px; }
.docs-body :deep(h1) { font-size: 24px; margin: 20px 0 12px; border-bottom: 1px solid var(--border); padding-bottom: 8px; }
.docs-body :deep(h2) { font-size: 20px; margin: 18px 0 10px; color: var(--text); }
.docs-body :deep(h3) { font-size: 16px; margin: 14px 0 8px; }
.docs-body :deep(p) { margin: 10px 0; }
.docs-body :deep(ul) { margin: 10px 0; padding-left: 24px; }
.docs-body :deep(li) { margin: 4px 0; }
.docs-body :deep(code) { background: rgba(255,255,255,.08); padding: 1px 6px; border-radius: 4px; font-size: 12px; }
.docs-body :deep(pre) { background: rgba(0,0,0,.3); padding: 14px; border-radius: 10px; overflow-x: auto; margin: 12px 0; }
.docs-body :deep(pre code) { background: transparent; padding: 0; }
.docs-body :deep(strong) { color: var(--text); }
.docs-body :deep(a) { color: var(--accent); }
</style>
