<script setup lang="ts">
import { ref, computed } from "vue";
import { useRouter } from "vue-router";
import { useToast } from "@/composables/toast";
import d1 from "@/assets/img/showcase/drama1.jpg?url";
import d2 from "@/assets/img/showcase/drama2.jpg?url";
import d3 from "@/assets/img/showcase/drama3.jpg?url";
import d4 from "@/assets/img/showcase/drama4.jpg?url";
import d5 from "@/assets/img/showcase/drama5.jpg?url";
import d6 from "@/assets/img/showcase/drama6.jpg?url";
import film1 from "@/assets/img/showcase/film1.jpg?url";
import film2 from "@/assets/img/showcase/film2.jpg?url";
import film3 from "@/assets/img/showcase/film3.jpg?url";
import film4 from "@/assets/img/showcase/film4.jpg?url";
import film5 from "@/assets/img/showcase/film5.jpg?url";
import myth1 from "@/assets/img/showcase/myth1.jpg?url";
import myth2 from "@/assets/img/showcase/myth2.jpg?url";
import myth3 from "@/assets/img/showcase/myth3.jpg?url";
import myth4 from "@/assets/img/showcase/myth4.jpg?url";
import myth5 from "@/assets/img/showcase/myth5.jpg?url";

const router = useRouter();
const toast = useToast();

type Template = {
  id: string;
  title: string;
  genre: string;
  category: "urban" | "suspense" | "ancient" | "youth" | "family" | "career";
  style: string;
  episodes: number;
  words: number;
  cover: string;
  description: string;
  tags: string[];
  script: string;
};

const TEMPLATES: Template[] = [
  {
    id: "rainy-coffee",
    title: "雨夜的咖啡店",
    genre: "都市情感",
    category: "urban",
    style: "韩漫厚涂 / 都市夜景",
    episodes: 3,
    words: 720,
    cover: d1,
    description: "雨夜,独自加班的女孩,与五年前不告而别的初恋在咖啡店重逢。一杯没喝完的美式,一段没说出口的真相。",
    tags: ["情感", "重逢", "雨夜"],
    script: `第1集 雨夜的咖啡店

林夏独自坐在靠窗的位置,笔记本电脑的屏幕光映在她的脸上。
店外的雨已经下了一整晚。她点了一杯美式,加了两次冰,又倒掉了两次。
门被推开,一个男人走进来,西装湿透,眉眼间有她熟悉的影子。

她几乎认不出他了——五年,他变了很多。
他也看见了她,愣在原地。

「……林夏?」

第2集 没说完的话

两人坐在角落的卡座,各自点了一杯美式。
他说他刚从机场回来,这五年一直在国外。
她问他为什么走。
他说,因为那一年她要结婚。

她笑了。「我没结婚。我连对象都没有。」
他沉默了很久。

「……我知道了。」

第3集 雨停了

凌晨两点,雨停了。
他们走出咖啡店,街灯把两个人的影子拉得很长。
他说,我可以重新开始吗。
她说,先把这杯美式喝完。

她把手机递过去。「先加个微信。」
他笑了,像极了五年前的样子。`,
  },
  {
    id: "goodnight-sms",
    title: "陌生人的晚安短信",
    genre: "都市情感",
    category: "urban",
    style: "治愈系 / 暖色调",
    episodes: 3,
    words: 680,
    cover: d2,
    description: "失意编剧误发短信给陌生号码,对方每晚都回复「晚安」。一段意料之外的陪伴,治愈了两个人。",
    tags: ["治愈", "陪伴", "陌生人"],
    script: `第1集 错发的短信

苏念对着手机发呆,本来要发给闺蜜的吐槽,不小心发给了通讯录里一个陌生号码。
她懊恼地删掉,准备道歉。

一分钟 后,对方回了:
「谢谢,晚安。」

第2集 每晚的回复

从那天起,对方每晚十点都会发来一句「晚安」。
她不知道对方是谁,但每天的吐槽也开始发给他。
他不一定回,但一定会说晚安。

有一天她问:你是谁?
他说:你不需要知道我是谁。

第3集 原来是你

苏念去医院复查,意外看见病床上的他。
她愣住了。
他笑了。「我每天跟你说晚安,只是希望有人记得,这世上还有人要好好活着。」

她握住他的手。「那我们都要好好活着。」`,
  },
  {
    id: "ten-year-letter",
    title: "第十年的信",
    genre: "悬疑推理",
    category: "suspense",
    style: "电影感 / 暗调",
    episodes: 3,
    words: 760,
    cover: d3,
    description: "一封十年前寄出的信,十年后才被拆开。寄信人已死,收信人是凶手还是受害者?",
    tags: ["悬疑", "时光", "反转"],
    script: `第1集 迟到的信

邮递员送来一封信。寄出日期是十年前。
收件人是本市著名的律师陈默。
寄件人是他的大学同学,十年前已经「意外死亡」的林萧。

陈默拆开信,里面只有一句话:
「我知道你做了什么。」

第2集 真相的碎片

陈默开始回忆十年前那个夏天。
林萧发现了他的秘密,威胁要公开。
那之后,林萧就「意外」溺水身亡。

他没有报警。因为,他就是凶手。
但这封信是谁寄的?

第3集 第十年的审判

陈默顺着信封上的邮戳,找到了寄信的小镇。
那是一个已经废弃的邮局。墙上贴满了照片,全是林萧生前搜集的证据。

最后一页写着:
「我原谅你了。但请你自首。」

陈默跪在地上,终于哭出了声。`,
  },
  {
    id: "changan-moon",
    title: "长安月下",
    genre: "古风仙侠",
    category: "ancient",
    style: "国风水墨 / 仙侠",
    episodes: 3,
    words: 800,
    cover: myth1,
    description: "大唐盛世,新科状元夜游长安,误入一处仙家别院,遇见等了他三百年的剑仙。",
    tags: ["古风", "仙侠", "三世情缘"],
    script: `第1集 长安灯会

上元夜,新科状元沈青云骑马游长安。
灯火辉煌,他却在朱雀桥头看见一个白衣女子,独立桥上,似乎在等谁。
「姑娘,你怎么一个人?」

女子回头,眉眼如画。「我在等一个人。」
「等谁?」
「你。」

第2集 三百年的约定

女子带他走进一处隐秘的庭院。
她叫青鸾,是这庭院的剑仙,也是三百年前沈青云前世的爱人。
那一世,他为了救她,魂飞魄散。

「我等了三百年,就是为了这一世再见你一面。」

第3集 今生的选择

天将亮,青鸾必须回到她的世界。
她问沈青云:「你愿意随我走吗?」
沈青云沉默。

他说:「我有今生的父母、恩师、未竟的事业。」
「但我会记你三百年。」

青鸾笑了。「好。那我也再等你三百年。」`,
  },
  {
    id: "peach-blossom",
    title: "桃花渡",
    genre: "古风言情",
    category: "ancient",
    style: "唯美古风 / 桃花",
    episodes: 3,
    words: 740,
    cover: myth2,
    description: "江南小镇,渡船少女与落魄书生的三世情缘,从桃花盛开讲到桃花凋零。",
    tags: ["古风", "言情", "江南"],
    script: `第1集 桃花盛开

江南三月,渡口桃花盛开。
摆渡的少女阿桃,日日送客过江。
这一天,船上坐了一位书生,抱着琴,神情落寞。

「公子,你从哪里来?」
「京城。」
「京城有什么好?」
「没什么好。」

第2集 桃花盛时

书生住进了渡口的客栈。
他日日来坐船,日日弹琴,日日望着江水发呆。
阿桃日日摆渡,日日为他唱歌。

他问她:「你为什么每天这么开心?」
她说:「因为明天太阳还会升起。」

第3集 桃花凋零

书生要走了,进京赶考。
阿桃送他到江边,塞给他一个桃花香囊。
「等你金榜题名,记得回来看桃花。」

多年后,京城传出一个消息:新科状元归乡途中,病逝于桃花渡口。
没有人知道,他怀里一直抱着那个香囊。`,
  },
  {
    id: "gaosan-year",
    title: "高三那年的约定",
    genre: "校园青春",
    category: "youth",
    style: "青春剧 / 暖色调",
    episodes: 3,
    words: 700,
    cover: film1,
    description: "高三开学第一天,两个曾经的青梅竹马被分到同一张课桌。一年前的误会还没解开。",
    tags: ["青春", "校园", "约定"],
    script: `第1集 同一张课桌

高三开学,班主任念座位表。
江晓晓听到自己的同桌名字,愣了——是陆言。
一年前,她不告而别,他追到火车站也没追上。

他走过来,放下书包,没说一句话。

第2集 没说出口的话

一个月过去,他们没说过一句完整的话。
直到期中考试前一天,江晓晓发了高烧。
陆言翻墙出去给她买药,第二天被班主任罚站。

「你不该翻墙。」
「你不该生病。」

第3集 高考之后

高考结束,大家在教室里最后一次合影。
陆言递给江晓晓一张纸条。
她回家才打开,上面写着:

「七年后,我们都在这座城市。如果那时你还是一个人,来找我。」

她笑了,回了一条短信:「好。」`,
  },
  {
    id: "grandma-garden",
    title: "奶奶的菜园",
    genre: "家庭伦理",
    category: "family",
    style: "温情 / 田园",
    episodes: 3,
    words: 720,
    cover: myth3,
    description: "城里长大的孙女暑假回到乡下,与独居的奶奶在菜园里度过最后一个夏天。",
    tags: ["亲情", "乡愁", "告别"],
    script: `第1集 回乡

七岁的小鱼被妈妈送到乡下,过暑假。
奶奶一个人在老房子里,菜园种满了瓜果。
小鱼一开始嫌乡下无聊,总想回家。

奶奶说:「奶奶教你怎么种菜。」

第2集 菜园里的夏天

整个暑假,小鱼跟着奶奶在菜园里。
浇水、施肥、抓虫子、摘黄瓜。
奶奶给她讲年轻的故事,讲爷爷,讲这条河。

小鱼不知不觉,忘了要回家的事。

第3集 最后一课

暑假快结束,妈妈来接小鱼。
临走前,奶奶把她叫到菜园,递给她一颗种子。
「这是你最喜欢的那种花。带回城里种,奶奶就一直在。」

小鱼抱紧奶奶,哭着说:「奶奶,明年我还来。」

奶奶笑着,摸了摸她的头。`,
  },
  {
    id: "investment-banker",
    title: "投行的第七年",
    genre: "职场商战",
    category: "career",
    style: "都市职场 / 紧张",
    episodes: 3,
    words: 760,
    cover: film3,
    description: "投行女高管第七年,面对 IPO 前夜的突发危机,她必须在 24 小时内做出改变一生的选择。",
    tags: ["职场", "商战", "女性"],
    script: `第1集 凌晨两点的电话

凌晨两点,周宁接到合伙人的电话。
「公司的核心数据有问题,审计师明天就要签字。」
她从床上坐起来,一夜未眠。

第二天早上八点,她必须给出答案。

第2集 真相

周宁带着团队连夜核查。
原来是一家子公司的 CFO 私挪了款项,做了假账。
现在,要么公开,公司 IPO 推迟;
要么掩盖,她从此背上这个秘密。

她问自己:这家公司,值得吗?

第3集 二十四年后的选择

清晨六点,周宁敲开了合伙人的门。
「我们要公开。」
「你疯了?」

「我做了七年的投行。这一次,我选择做一个好人。」

第二天,新闻头条。
周宁辞职,带着一封公开信,离开了这个行业。`,
  },
];

const CATEGORIES = [
  { id: "all",     label: "全部",     icon: "✦" },
  { id: "urban",   label: "都市情感", icon: "♡" },
  { id: "suspense",label: "悬疑推理", icon: "?" },
  { id: "ancient", label: "古风仙侠", icon: "卐" },
  { id: "youth",   label: "校园青春", icon: "✿" },
  { id: "family",  label: "家庭伦理", icon: "▲" },
  { id: "career",  label: "职场商战", icon: "◆" },
];

const category = ref<string>("all");
const search = ref("");
const previewing = ref<Template | null>(null);

const filtered = computed(() => {
  let list = TEMPLATES;
  if (category.value !== "all") list = list.filter(t => t.category === category.value);
  if (search.value.trim()) {
    const q = search.value.trim().toLowerCase();
    list = list.filter(t =>
      t.title.toLowerCase().includes(q) ||
      t.description.toLowerCase().includes(q) ||
      t.tags.some(tag => tag.toLowerCase().includes(q))
    );
  }
  return list;
});

function useTemplate(t: Template) {
  // 把模板内容存入 sessionStorage,新建任务页面读取后预填
  sessionStorage.setItem("agnes_template", JSON.stringify({ title: t.title, script: t.script }));
  toast.ok(`已加载模板「${t.title}」,跳转中…`);
  setTimeout(() => router.push({ name: "new" }), 400);
}

function copyScript(t: Template) {
  navigator.clipboard.writeText(t.script).then(
    () => toast.ok("已复制到剪贴板"),
    () => toast.err("复制失败,请手动选择")
  );
}

function countEpisodes(text: string) {
  return (text.match(/第\s*\d+\s*[集回章]|Episode\s+\d+|EP\s*\d+/gi) || []).length;
}
</script>

<template>
  <header class="topbar">
    <div>
      <div class="topbar-title">剧本模板</div>
      <div class="topbar-sub">工作台 / 资源 / 精选剧本模板 · 一键应用到新建任务</div>
    </div>
    <div class="topbar-right">
      <input v-model="search" type="search" placeholder="搜索标题、标签..." class="search" />
      <button class="btn primary" @click="router.push({ name: 'new' })">
        <svg class="btn-icon" viewBox="0 0 24 24"><path d="M12 5v14M5 12h14"/></svg>
        新建空白剧本
      </button>
    </div>
  </header>

  <div>
    <!-- Category filter -->
    <div class="cat-row">
      <button
        v-for="c in CATEGORIES"
        :key="c.id"
        class="cat-chip"
        :class="{ active: category === c.id }"
        @click="category = c.id"
      >
        <span class="cat-ico">{{ c.icon }}</span>
        {{ c.label }}
      </button>
    </div>

    <!-- Template grid -->
    <div v-if="filtered.length === 0" class="empty">
      没有匹配的剧本模板。试试别的分类或搜索词。
    </div>
    <div v-else class="tpl-grid">
      <article v-for="t in filtered" :key="t.id" class="tpl-card">
        <div class="tpl-cover">
          <img :src="t.cover" :alt="t.title" loading="lazy" />
          <div class="tpl-cover-mask"></div>
          <div class="tpl-genre">{{ t.genre }}</div>
          <div class="tpl-ep">{{ t.episodes }} 集 · {{ t.words }} 字</div>
        </div>
        <div class="tpl-body">
          <h3 class="tpl-title">{{ t.title }}</h3>
          <p class="tpl-desc">{{ t.description }}</p>
          <div class="tpl-tags">
            <span v-for="tag in t.tags" :key="tag" class="tag-chip">#{{ tag }}</span>
          </div>
          <div class="tpl-meta">
            <span class="meta-item">
              <span class="meta-lbl">风格</span>
              <span class="meta-val">{{ t.style }}</span>
            </span>
          </div>
          <div class="tpl-actions">
            <button class="btn sm" @click="previewing = t">
              <svg class="btn-icon" viewBox="0 0 24 24"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>
              预览
            </button>
            <button class="btn sm" @click="copyScript(t)">
              <svg class="btn-icon" viewBox="0 0 24 24"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>
              复制
            </button>
            <button class="btn sm primary tpl-use" @click="useTemplate(t)">
              <svg class="btn-icon" viewBox="0 0 24 24"><path d="M5 12h14M13 6l6 6-6 6"/></svg>
              使用此模板
            </button>
          </div>
        </div>
      </article>
    </div>
  </div>

  <!-- Preview modal -->
  <Teleport to="body">
    <div v-if="previewing" class="modal" @click.self="previewing = null">
      <div class="modal-mask"></div>
      <div class="modal-card">
        <div class="modal-head">
          <div>
            <div class="modal-eyebrow">{{ previewing.genre }} · {{ previewing.style }}</div>
            <h3>{{ previewing.title }}</h3>
          </div>
          <button class="modal-close" @click="previewing = null">✕</button>
        </div>
        <div class="modal-body">
          <p class="modal-desc">{{ previewing.description }}</p>
          <div class="modal-tags">
            <span v-for="tag in previewing.tags" :key="tag" class="tag-chip">#{{ tag }}</span>
          </div>
          <pre class="modal-script">{{ previewing.script }}</pre>
        </div>
        <div class="modal-foot">
          <span class="muted small">
            自动检测到 {{ countEpisodes(previewing.script) }} 个剧集标记
          </span>
          <div class="foot-actions">
            <button class="btn" @click="copyScript(previewing); previewing = null">复制全文</button>
            <button class="btn primary" @click="useTemplate(previewing)">
              使用此模板
              <svg class="btn-icon" viewBox="0 0 24 24"><path d="M5 12h14M13 6l6 6-6 6"/></svg>
            </button>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
/* === Category chips === */
.cat-row {
  display: flex; flex-wrap: wrap; gap: 8px;
  margin-bottom: 20px;
  padding: 6px;
  background: var(--card);
  border: 1px solid var(--border-soft);
  border-radius: var(--r-md);
  box-shadow: var(--shadow-xs);
  width: fit-content;
  max-width: 100%;
}
.cat-chip {
  display: inline-flex; align-items: center; gap: 6px;
  padding: 7px 14px;
  border: 0;
  background: transparent;
  color: var(--text-soft);
  font-size: 13px;
  font-weight: 500;
  border-radius: var(--r-sm);
  cursor: pointer;
  transition: all .15s var(--ease);
  white-space: nowrap;
}
.cat-chip:hover { background: var(--card-soft); color: var(--text); }
.cat-chip.active {
  background: var(--grad);
  color: #fff;
  font-weight: 600;
  box-shadow: var(--shadow-glow);
}
.cat-ico { font-size: 11px; opacity: .9; }

/* === Search === */
.search {
  background: var(--card);
  border: 1px solid var(--border);
  color: var(--text);
  border-radius: var(--r-sm);
  padding: 9px 14px;
  outline: none;
  font-size: 13px;
  min-width: 260px;
  transition: border-color .15s, box-shadow .15s;
}
.search:focus {
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(20, 184, 166, .12);
}
.search::placeholder { color: var(--text-faint); }

/* === Grid === */
.tpl-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 18px;
}
.tpl-card {
  background: var(--card);
  border-radius: var(--r-lg);
  border: 1px solid var(--border-soft);
  overflow: hidden;
  box-shadow: var(--shadow-xs);
  display: flex;
  flex-direction: column;
  transition: all .2s var(--ease);
}
.tpl-card:hover {
  transform: translateY(-3px);
  box-shadow: var(--shadow-md);
  border-color: rgba(20, 184, 166, .35);
}
.tpl-cover {
  position: relative;
  aspect-ratio: 16/9;
  background: var(--card-soft);
  overflow: hidden;
}
.tpl-cover img {
  width: 100%; height: 100%; object-fit: cover;
  transition: transform .4s var(--ease);
}
.tpl-card:hover .tpl-cover img { transform: scale(1.05); }
.tpl-cover-mask {
  position: absolute; inset: 0;
  background: linear-gradient(to top, rgba(20, 30, 60, .55) 0%, rgba(20, 30, 60, .12) 40%, transparent 70%);
}
.tpl-genre {
  position: absolute;
  top: 12px; left: 12px;
  padding: 3px 12px;
  background: var(--grad);
  color: #fff;
  font-size: 11.5px;
  font-weight: 600;
  border-radius: 999px;
  letter-spacing: .02em;
  box-shadow: var(--shadow-glow);
}
.tpl-ep {
  position: absolute;
  bottom: 12px; left: 12px;
  padding: 4px 10px;
  background: rgba(255, 255, 255, .92);
  color: var(--text);
  font-size: 11.5px;
  font-weight: 600;
  border-radius: 6px;
  backdrop-filter: blur(8px);
}
.tpl-body {
  padding: 16px 18px 18px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  flex: 1;
}
.tpl-title {
  font-size: 16px;
  font-weight: 700;
  color: var(--text);
  line-height: 1.35;
}
.tpl-desc {
  font-size: 12.5px;
  color: var(--text-soft);
  line-height: 1.65;
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.tpl-tags { display: flex; flex-wrap: wrap; gap: 5px; }
.tag-chip {
  padding: 2px 9px;
  background: var(--grad-soft);
  color: var(--primary-deep);
  font-size: 11px;
  font-weight: 500;
  border-radius: 999px;
  border: 1px solid rgba(20, 184, 166, .2);
}
.tpl-meta {
  padding-top: 8px;
  border-top: 1px dashed var(--border);
}
.meta-item { display: flex; flex-direction: column; gap: 2px; }
.meta-lbl { font-size: 11px; color: var(--text-muted); }
.meta-val { font-size: 12.5px; color: var(--text); font-weight: 500; }
.tpl-actions {
  display: flex; gap: 8px;
  margin-top: auto;
  padding-top: 8px;
}
.tpl-actions .btn { flex: 1; padding: 7px 8px; font-size: 12px; }
.tpl-actions .btn .btn-icon { width: 13px; height: 13px; }
.tpl-use { flex: 1.5 !important; }

/* === Modal === */
.modal {
  position: fixed; inset: 0; z-index: 100;
  display: flex; align-items: center; justify-content: center;
  padding: 24px;
}
.modal-mask {
  position: absolute; inset: 0;
  background: rgba(20, 30, 60, .35);
  backdrop-filter: blur(6px);
}
.modal-card {
  position: relative;
  width: min(760px, 100%);
  max-height: calc(100vh - 48px);
  display: flex; flex-direction: column;
  background: var(--card);
  border: 1px solid var(--border-soft);
  border-radius: var(--r-lg);
  box-shadow: var(--shadow-lg);
  overflow: hidden;
}
.modal-head {
  display: flex; align-items: flex-start; justify-content: space-between;
  padding: 20px 24px 16px;
  border-bottom: 1px solid var(--border-soft);
  gap: 16px;
}
.modal-eyebrow {
  font-size: 11.5px;
  color: var(--primary-deep);
  font-weight: 600;
  letter-spacing: .04em;
  margin-bottom: 6px;
}
.modal-head h3 {
  font-size: 18px;
  font-weight: 700;
  color: var(--text);
}
.modal-close {
  background: transparent;
  border: 0;
  color: var(--text-muted);
  font-size: 20px;
  cursor: pointer;
  padding: 4px 12px;
  border-radius: 8px;
  flex-shrink: 0;
}
.modal-close:hover { background: var(--card-soft); color: var(--text); }

.modal-body {
  flex: 1; overflow-y: auto;
  padding: 20px 24px;
  background: var(--bg-soft);
}
.modal-desc {
  font-size: 13.5px;
  color: var(--text-soft);
  line-height: 1.75;
  margin: 0 0 12px;
}
.modal-tags { display: flex; flex-wrap: wrap; gap: 5px; margin-bottom: 14px; }
.modal-script {
  background: var(--card);
  border: 1px solid var(--border-soft);
  border-radius: var(--r-sm);
  padding: 16px 18px;
  font-family: ui-monospace, "SFMono-Regular", Menlo, Consolas, monospace;
  font-size: 12.5px;
  line-height: 1.85;
  color: var(--text);
  white-space: pre-wrap;
  word-break: break-word;
  margin: 0;
  max-height: 400px;
  overflow-y: auto;
}

.modal-foot {
  display: flex; align-items: center; justify-content: space-between;
  padding: 14px 24px;
  border-top: 1px solid var(--border-soft);
  background: var(--card);
  gap: 12px;
  flex-wrap: wrap;
}
.foot-actions { display: flex; gap: 8px; }

/* === Responsive === */
@media (max-width: 720px) {
  .tpl-grid { grid-template-columns: 1fr; }
  .search { min-width: 0; width: 100%; }
}
</style>
