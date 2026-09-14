// === Agnes AI Studio — frontend app ===
// Vanilla JS. No build step.

(function () {
  "use strict";

  const $ = (sel, root) => (root || document).querySelector(sel);
  const $$ = (sel, root) => Array.from((root || document).querySelectorAll(sel));
  const el = (tag, attrs, children) => {
    const e = document.createElement(tag);
    if (attrs) Object.keys(attrs).forEach((k) => {
      if (k === "class") e.className = attrs[k];
      else if (k === "html") e.innerHTML = attrs[k];
      else if (k.startsWith("on") && typeof attrs[k] === "function") e.addEventListener(k.slice(2), attrs[k]);
      else if (k === "data") Object.keys(attrs.data).forEach((dk) => e.dataset[dk] = attrs.data[dk]);
      else e.setAttribute(k, attrs[k]);
    });
    if (children) (Array.isArray(children) ? children : [children]).forEach((c) => {
      if (c == null) return;
      if (c instanceof Node) { e.appendChild(c); }
      else { e.appendChild(document.createTextNode(String(c))); }
    });
    return e;
  };

  const state = {
    view: "home",
    currentJobId: null,
    pollHandle: null,
    queuePollHandle: null,
    search: "",
    detailCache: {},
  };

  // ---------- API ----------
  const api = {
    base: "",
    async req(method, path, body) {
      const opts = { method, headers: { "Content-Type": "application/json" } };
      if (body !== undefined) opts.body = JSON.stringify(body);
      const r = await fetch(this.base + path, opts);
      if (!r.ok) {
        let msg = r.status + " " + r.statusText;
        try { const j = await r.json(); if (j.error) msg = j.error; } catch (e) {}
        throw new Error(msg);
      }
      if (r.status === 204) return null;
      const ct = r.headers.get("Content-Type") || "";
      if (ct.indexOf("application/json") >= 0) return r.json();
      return r.text();
    },
    get(p) { return this.req("GET", p); },
    post(p, b) { return this.req("POST", p, b); },
    del(p) { return this.req("DELETE", p); },
    put(p, b) { return this.req("PUT", p, b); },
  };

  // ---------- Toast ----------
  const toast = (msg, kind) => {
    const t = $("#toast");
    t.textContent = msg;
    t.className = "toast show " + (kind || "");
    clearTimeout(toast._h);
    toast._h = setTimeout(() => { t.className = "toast"; }, 3000);
  };

  // ---------- Modal ----------
  function showModal(content, opts) {
    opts = opts || {};
    const m = $("#modal");
    const c = $("#modal-card");
    c.innerHTML = "";
    if (typeof content === "string") c.innerHTML = content;
    else c.appendChild(content);
    m.classList.remove("hidden");
    if (opts.onClose) m._onClose = opts.onClose; else m._onClose = null;
  }
  function hideModal() {
    const m = $("#modal");
    m.classList.add("hidden");
    if (m._onClose) try { m._onClose(); } catch (e) {}
    m._onClose = null;
  }
  $("#modal").addEventListener("click", (e) => { if (e.target.classList.contains("modal-backdrop")) hideModal(); });

  // ---------- View routing ----------
  function go(view, params) {
    state.view = view;
    if (params && params.id) state.currentJobId = params.id;
    $$(".view").forEach((v) => v.classList.toggle("active", v.dataset.view === view));
    $$("#nav .tab").forEach((t) => t.classList.toggle("active", t.dataset.view === view));
    if (view === "home") loadHome();
    if (view === "new") $("#job-script").focus();
    if (view === "jobs") loadJobs();
    if (view === "detail") loadJobDetail(state.currentJobId);
    if (view === "settings") loadSettings();
    if (view === "docs") loadDocs();
    // scroll to top
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  $$("[data-goto]").forEach((b) => b.addEventListener("click", () => go(b.dataset.goto)));
  $$("#nav .tab").forEach((t) => t.addEventListener("click", () => go(t.dataset.view)));

  // ---------- Home ----------
  async function loadHome() {
    const card = $("#recent-jobs");
    card.innerHTML = "<div class=\"empty\"><div class=\"ico\">...</div><div>Loading</div></div>";
    try {
      const r = await api.get("/api/jobs");
      const list = (r.jobs || []).slice(0, 4);
      if (list.length === 0) {
        card.innerHTML = "<div class=\"empty\"><div class=\"ico\">🎬</div><div>No jobs yet. Create your first one to get started.</div></div>";
        return;
      }
      card.innerHTML = "";
      list.forEach((j) => card.appendChild(jobRow(j)));
    } catch (e) {
      card.innerHTML = "<div class=\"empty\">" + escapeHTML(e.message) + "</div>";
    }
  }

  // ---------- Jobs list ----------
  async function loadJobs() {
    const list = $("#job-list");
    list.innerHTML = "<div class=\"empty\"><div class=\"ico\">...</div><div>Loading</div></div>";
    try {
      const r = await api.get("/api/jobs");
      let jobs = r.jobs || [];
      const q = state.search.trim().toLowerCase();
      if (q) jobs = jobs.filter((j) => (j.title || "").toLowerCase().indexOf(q) >= 0);
      if (jobs.length === 0) {
        list.innerHTML = "<div class=\"empty\"><div class=\"ico\">🎬</div><div>No jobs match. Try a new one.</div></div>";
        return;
      }
      list.innerHTML = "";
      jobs.forEach((j) => list.appendChild(jobRow(j)));
    } catch (e) {
      list.innerHTML = "<div class=\"empty\">" + escapeHTML(e.message) + "</div>";
    }
  }

  function jobRow(j) {
    const row = el("div", { class: "job-row", onclick: () => go("detail", { id: j.id }) });
    row.appendChild(el("div", { class: "grow" }, [
      el("div", { class: "title" }, j.title || "Untitled"),
      el("div", { class: "meta" }, [
        "" + (j.num_episodes || 0) + " episodes · " + (j.num_characters || 0) + " characters · " + (j.num_videos || 0) + " videos · " + formatDate(j.created_at),
      ]),
    ]));
    const status = el("div", { class: "status " + (j.status || "pending") }, j.status || "pending");
    row.appendChild(status);
    return row;
  }

  $("#job-search").addEventListener("input", (e) => {
    state.search = e.target.value;
    loadJobs();
  });

  // ---------- New job ----------
  $("#parse-btn").addEventListener("click", () => parseOnly());
  $("#create-btn").addEventListener("click", () => createJob());

  async function parseOnly() {
    const script = $("#job-script").value.trim();
    if (!script) { toast("Please paste a script first", "err"); return; }
    try {
      const r = await api.post("/api/jobs", { title: $("#job-title").value || "Untitled", script, style: $("#job-style").value });
      state.currentJobId = r.id;
      renderParsePreview(r.episodes || []);
    } catch (e) { toast(e.message, "err"); }
  }

  function renderParsePreview(eps) {
    const p = $("#parse-preview");
    p.classList.remove("hidden");
    p.innerHTML = "<h3>Detected " + eps.length + " episode" + (eps.length === 1 ? "" : "s") + "</h3>";
    eps.forEach((ep) => {
      const e = el("div", { class: "ep" });
      e.appendChild(el("b", null, "Episode " + ep.index));
      e.appendChild(document.createTextNode(ep.title || ""));
      e.appendChild(el("div", { class: "body" }, ep.body || "(empty)"));
      p.appendChild(e);
    });
  }

  async function createJob() {
    const title = $("#job-title").value.trim();
    const script = $("#job-script").value.trim();
    const style = $("#job-style").value.trim();
    if (!script) { toast("Script is required", "err"); return; }
    try {
      const r = await api.post("/api/jobs", { title: title || "Untitled", script, style });
      toast("Job created", "ok");
      $("#job-title").value = "";
      $("#job-script").value = "";
      $("#job-style").value = "";
      $("#parse-preview").classList.add("hidden");
      go("detail", { id: r.id });
    } catch (e) { toast(e.message, "err"); }
  }

  // ---------- Job detail ----------
  async function loadJobDetail(id) {
    if (!id) { go("jobs"); return; }
    const root = $("#job-detail");
    root.innerHTML = "<div class=\"card glass\"><div class=\"empty\"><div class=\"ico\">...</div><div>Loading</div></div></div>";
    stopPolling();
    try {
      const j = await api.get("/api/jobs/" + id);
      state.detailCache[id] = j;
      renderJobDetail(j);
      if (j.status !== "done" && j.status !== "failed" && j.status !== "cancelled") startPolling(id);
    } catch (e) {
      root.innerHTML = "<div class=\"card glass\"><div class=\"empty\">" + escapeHTML(e.message) + "</div></div>";
    }
  }

  function startPolling(id) {
    stopPolling();
    state.pollHandle = setInterval(async () => {
      try {
        const j = await api.get("/api/jobs/" + id);
        state.detailCache[id] = j;
        renderJobDetail(j);
        if (j.status === "done" || j.status === "failed" || j.status === "cancelled") stopPolling();
      } catch (e) { /* keep polling */ }
    }, 3000);
  }
  function stopPolling() {
    if (state.pollHandle) { clearInterval(state.pollHandle); state.pollHandle = null; }
  }

  function renderJobDetail(j) {
    const root = $("#job-detail");
    root.innerHTML = "";

    // Head
    const head = el("div", { class: "card glass" });
    const headRow = el("div", { class: "detail-head" }, [
      el("div", null, [
        el("h2", null, j.title || "Untitled"),
        el("div", { class: "muted" }, "Job " + j.id + " · created " + formatDate(j.created_at)),
      ]),
    ]);
    const actions = el("div", { class: "detail-actions" });
    if (j.status === "pending" || j.status === "failed" || j.status === "cancelled") {
      actions.appendChild(el("button", { class: "btn primary", onclick: () => runJob(j.id) }, [iconPlay(), "Run all"]));
    }
    if (j.status === "pending") {
      actions.appendChild(el("button", { class: "btn glass", onclick: () => deleteJob(j.id) }, "Delete"));
    }
    if (j.status !== "done" && j.status !== "pending" && j.status !== "cancelled") {
      actions.appendChild(el("button", { class: "btn glass", onclick: () => cancelJob(j.id) }, "Cancel"));
    }
    headRow.appendChild(actions);
    head.appendChild(headRow);

    // Progress
    const wrap = el("div", { class: "progress-wrap" });
    wrap.appendChild(el("div", { class: "progress-bar" }, [el("div", { style: "width:" + (j.progress || 0) + "%" })]));
    const meta = el("div", { class: "progress-meta" });
    meta.appendChild(el("span", null, "Status: " + (j.status || "pending")));
    meta.appendChild(el("span", null, (j.progress || 0) + "%"));
    wrap.appendChild(meta);
    if (j.error) wrap.appendChild(el("div", { class: "test-result err", style: "margin-top:12px" }, "Error: " + j.error));
    head.appendChild(wrap);
    root.appendChild(head);

    // Characters
    if (j.characters && j.characters.length > 0) {
      const sec = el("div", null, [el("div", { class: "section-title" }, ["Characters ", el("span", { class: "count" }, "(" + j.characters.length + ")")])]);
      const grid = el("div", { class: "char-grid" });
      j.characters.forEach((c) => grid.appendChild(characterCard(c)));
      sec.appendChild(grid);
      root.appendChild(sec);
    }
    // Props
    if (j.props && j.props.length > 0) {
      const sec = el("div", null, [el("div", { class: "section-title" }, ["Props ", el("span", { class: "count" }, "(" + j.props.length + ")")])]);
      const grid = el("div", { class: "char-grid" });
      j.props.forEach((p) => grid.appendChild(propCard(p)));
      sec.appendChild(grid);
      root.appendChild(sec);
    }
    // Episodes
    if (j.episodes && j.episodes.length > 0) {
      const sec = el("div", null, [el("div", { class: "section-title" }, ["Episodes ", el("span", { class: "count" }, "(" + j.episodes.length + ")")])]);
      j.episodes.forEach((ep) => sec.appendChild(episodeCard(j, ep)));
      root.appendChild(sec);
    } else {
      root.appendChild(el("div", { class: "card glass" }, [el("div", { class: "empty" }, "No episodes detected. Add markers like 第1集 or Episode 1 to your script, then re-create the job.")]));
    }
  }

  function characterCard(c) {
    const card = el("div", { class: "char-card glass" });
    const img = el("div", { class: "img" });
    if (c.image_url) img.appendChild(el("img", { src: c.image_url, alt: c.name, loading: "lazy" }));
    else img.appendChild(el("div", { class: "pending" }, "No image yet"));
    card.appendChild(img);
    card.appendChild(el("div", { class: "name" }, c.name));
    card.appendChild(el("div", { class: "role" }, c.role || "character"));
    if (c.appearance) card.appendChild(el("div", { class: "app" }, c.appearance));
    return card;
  }
  function propCard(p) {
    const card = el("div", { class: "char-card glass" });
    const img = el("div", { class: "img" });
    if (p.image_url) img.appendChild(el("img", { src: p.image_url, alt: p.name, loading: "lazy" }));
    else img.appendChild(el("div", { class: "pending" }, "No image yet"));
    card.appendChild(img);
    card.appendChild(el("div", { class: "name" }, p.name));
    card.appendChild(el("div", { class: "role" }, p.kind || "prop"));
    if (p.description) card.appendChild(el("div", { class: "app" }, p.description));
    return card;
  }

  function episodeCard(j, ep) {
    const card = el("div", { class: "episode glass" });
    const head = el("h4", null, [
      "Episode " + ep.index + ": " + (ep.title || ""),
      el("span", { class: "ep-state " + (ep.state || "pending") }, ep.state || "pending"),
    ]);
    card.appendChild(head);
    if (ep.body) card.appendChild(el("div", { class: "summary" }, truncate(ep.body, 200)));
    // progress
    if (ep.state && ep.state !== "done" && ep.state !== "failed" && ep.state !== "skipped" && ep.state !== "pending") {
      card.appendChild(el("div", { class: "progress-wrap" }, [
        el("div", { class: "progress-bar" }, [el("div", { style: "width:" + (ep.progress || 0) + "%" })]),
      ]));
    }
    if (ep.error) card.appendChild(el("div", { class: "test-result err" }, ep.error));
    // scenes
    if (ep.scenes && ep.scenes.length > 0) {
      const sl = el("div", { class: "scene-list" });
      ep.scenes.forEach((s) => {
        const sc = el("div", { class: "scene" });
        if (s.image_url) sc.appendChild(el("img", { src: s.image_url, alt: s.heading, loading: "lazy" }));
        else sc.appendChild(el("div", { class: "pending" }, "Scene " + s.index));
        if (s.heading) sc.appendChild(el("div", { class: "heading" }, s.heading));
        sl.appendChild(sc);
      });
      card.appendChild(sl);
    }
    // video
    if (ep.video_url) {
      const v = el("video", { controls: true, preload: "metadata", src: ep.video_url });
      v.style.maxHeight = "420px";
      card.appendChild(el("div", { class: "video-wrap" }, [v]));
    }
    // actions
    const acts = el("div", { class: "ep-actions" });
    if (j.status === "pending" || j.status === "failed" || j.status === "done" || j.status === "cancelled") {
      acts.appendChild(el("button", { class: "btn glass", onclick: () => runEpisode(j.id, ep.index) }, [iconPlay(), "Generate this episode"]));
    }
    card.appendChild(acts);
    return card;
  }

  function runJob(id) { api.post("/api/jobs/" + id + "/run").then(() => { toast("Queued", "ok"); loadJobDetail(id); }).catch((e) => toast(e.message, "err")); }
  function runEpisode(id, n) { api.post("/api/jobs/" + id + "/episodes/" + n + "/run").then(() => { toast("Queued episode " + n, "ok"); loadJobDetail(id); }).catch((e) => toast(e.message, "err")); }
  function cancelJob(id) { api.post("/api/jobs/" + id + "/cancel").then(() => { toast("Cancellation requested", "ok"); loadJobDetail(id); }).catch((e) => toast(e.message, "err")); }
  function deleteJob(id) {
    showModal("<h3>Delete this job?</h3><p class=\"muted\">This removes the job and all its generated assets. This cannot be undone.</p><div class=\"actions\"><button class=\"btn ghost\" onclick=\"hideModal()\">Cancel</button><button class=\"btn primary\" id=\"confirm-del\">Delete</button></div>");
    $("#confirm-del").addEventListener("click", () => { api.del("/api/jobs/" + id).then(() => { hideModal(); toast("Deleted", "ok"); go("jobs"); }).catch((e) => toast(e.message, "err")); });
  }

  // ---------- Settings ----------
  async function loadSettings() {
    try {
      const s = await api.get("/api/settings");
      $("#key-status").textContent = s.api_key_set ? "Configured · " + s.api_route + " · " + s.base_url : "Not set";
      const route = s.api_route || "international";
      const r = document.querySelector("input[name=route][value=" + route + "]");
      if (r) r.checked = true;
      const kv = $("#runtime-info");
      kv.innerHTML = "";
      appendKV(kv, "Base URL", s.base_url || "—");
      appendKV(kv, "Route", s.api_route || "—");
      appendKV(kv, "API key", s.api_key_set ? "Set (encrypted at rest)" : "Not set");
      appendKV(kv, "Concurrency", s.concurrency || 1);
      appendKV(kv, "Poll interval", (s.poll_interval_s || 5) + "s");
      appendKV(kv, "Server", "agnesai-studio 1.0.0");
    } catch (e) { toast(e.message, "err"); }
  }
  function appendKV(parent, k, v) {
    parent.appendChild(el("div", { class: "k" }, k));
    parent.appendChild(el("div", { class: "v" }, v));
  }

  $("#save-key").addEventListener("click", async () => {
    const k = $("#api-key").value.trim();
    if (!k) { toast("API key is empty", "err"); return; }
    try {
      await api.post("/api/settings/key", { api_key: k });
      $("#api-key").value = "";
      toast("Saved", "ok");
      loadSettings();
    } catch (e) { toast(e.message, "err"); }
  });
  $("#test-key").addEventListener("click", async () => {
    const k = $("#api-key").value.trim();
    const tr = $("#test-result");
    tr.className = "test-result"; tr.classList.remove("hidden");
    tr.textContent = "Testing...";
    try {
      const r = await api.post("/api/settings/key/test", { api_key: k });
      if (r.ok) { tr.className = "test-result ok"; tr.textContent = "Connected (" + r.duration + "ms) via " + r.base_url; }
      else { tr.className = "test-result err"; tr.textContent = "Failed: " + (r.message || "unknown error"); }
    } catch (e) { tr.className = "test-result err"; tr.textContent = e.message; }
  });
  $("#clear-key").addEventListener("click", async () => {
    if (!confirm("Clear the stored API key?")) return;
    try { await api.del("/api/settings/key"); toast("Cleared", "ok"); loadSettings(); }
    catch (e) { toast(e.message, "err"); }
  });
  $("#toggle-key").addEventListener("click", () => {
    const i = $("#api-key");
    if (i.type === "password") { i.type = "text"; $("#toggle-key").textContent = "Hide"; }
    else { i.type = "password"; $("#toggle-key").textContent = "Show"; }
  });
  $("#save-route").addEventListener("click", async () => {
    const r = document.querySelector("input[name=route]:checked");
    if (!r) return;
    const res = $("#route-result");
    res.className = "test-result"; res.classList.remove("hidden");
    res.textContent = "Saving...";
    try {
      const x = await api.post("/api/settings/route", { route: r.value });
      res.className = "test-result ok";
      res.textContent = "Saved. Restart the server to apply the new route.";
      loadSettings();
    } catch (e) { res.className = "test-result err"; res.textContent = e.message; }
  });

  // ---------- Docs ----------
  async function loadDocs() {
    const root = $("#docs-content");
    if (root.dataset.loaded) { return; }
    try {
      const r = await fetch("/api/docs");
      if (!r.ok) throw new Error("HTTP " + r.status);
      const md = await r.text();
      root.innerHTML = window.MiniMark.render(md);
      root.dataset.loaded = "1";
    } catch (e) {
      root.innerHTML = "<p class=\"muted\">Failed to load docs: " + escapeHTML(e.message) + "</p>";
    }
  }

  // ---------- Queue pill ----------
  function startQueuePolling() {
    if (state.queuePollHandle) return;
    const tick = async () => {
      try {
        const r = await api.get("/api/queue");
        const q = r.queue || {};
        const pill = $("#queue-pill");
        const text = $("#queue-text");
        if (q.running > 0 || q.queued > 0) {
          pill.classList.remove("hidden");
          pill.classList.add("busy");
          text.textContent = q.running + " running · " + q.queued + " queued";
        } else {
          pill.classList.add("hidden");
          pill.classList.remove("busy");
        }
      } catch (e) { /* keep silent */ }
    };
    tick();
    state.queuePollHandle = setInterval(tick, 4000);
  }

  // ---------- Theme ----------
  $("#theme-toggle").addEventListener("click", () => {
    const b = document.body;
    if (b.classList.contains("theme-dark")) { b.classList.remove("theme-dark"); b.classList.add("theme-light"); localStorage.setItem("theme", "light"); }
    else { b.classList.remove("theme-light"); b.classList.add("theme-dark"); localStorage.setItem("theme", "dark"); }
  });
  (function () {
    const t = localStorage.getItem("theme");
    if (t === "light") { document.body.classList.remove("theme-dark"); document.body.classList.add("theme-light"); }
  })();

  // ---------- Liquid glass: mouse-following highlight ----------
  function bindLiquid() {
    $$(".glass").forEach((node) => {
      node.addEventListener("mousemove", (e) => {
        const r = node.getBoundingClientRect();
        const x = ((e.clientX - r.left) / r.width) * 100;
        const y = ((e.clientY - r.top) / r.height) * 100;
        node.style.setProperty("--mx", x + "%");
        node.style.setProperty("--my", y + "%");
      });
    });
  }

  // ---------- Utilities ----------
  function escapeHTML(s) { return (s || "").replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;"); }
  function truncate(s, n) { s = s || ""; if (s.length <= n) return s; return s.substring(0, n) + "..."; }
  function formatDate(s) { if (!s) return "—"; const d = new Date(s); if (isNaN(d.getTime())) return s; return d.toLocaleString(); }
  function iconPlay() {
    const s = document.createElementNS("http://www.w3.org/2000/svg", "svg");
    s.setAttribute("viewBox", "0 0 24 24"); s.setAttribute("width", "14"); s.setAttribute("height", "14");
    s.setAttribute("fill", "currentColor");
    s.innerHTML = "<path d=\"M8 5v14l11-7z\"/>";
    return s;
  }

  // ---------- Boot ----------
  document.addEventListener("DOMContentLoaded", () => {
    bindLiquid();
    startQueuePolling();
    // re-bind liquid for dynamically added nodes
    const mo = new MutationObserver(() => bindLiquid());
    mo.observe(document.body, { childList: true, subtree: true });
    // start at home
    go("home");
  });
})();