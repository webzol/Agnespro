// styles_library.go — 新版"风格库"数据源
//
// 把视觉画风、叙事题材、比例、剧集数等所有"风格库"相关数据集中管理,
// 通过 GET /api/styles/library 一次返回给前端,前端按需渲染 modal / 下拉。
//
// 数据驱动:后续要新增风格,只需要在本文件追加条目,无需改前端代码。

package api

import (
	"net/http"

	"agnesai/studio/internal/types"
)

// visualLibraryItem:视觉画风库的一条
//   - ID: 英文短码,作为 image_prompt 的 hook
//   - Name: 中文显示名
//   - Category: "2d" / "3d" / "真人",对应前端 modal 的 tabs
//   - Desc: 一句话描述
//   - PromptHint: 英文风格关键词,会拼到 image/video prompt 末尾
//   - Preview: 缩略图 URL(占位图 / 待 AI 生成)
type visualLibraryItem struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Category   string `json:"category"`
	Desc       string `json:"desc"`
	PromptHint string `json:"prompt_hint"`
	Preview    string `json:"preview"`
}

// genreLibraryItem:叙事题材库的一条
//   - PromptHint: 给 AI 生成剧本时附加的调性描述
type genreLibraryItem struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	PromptHint string `json:"prompt_hint"`
}

// episodeCountOption:剧集数预设
type episodeCountOption struct {
	Value int    `json:"value"`
	Label string `json:"label"`
}

// ===== 视觉画风库(30 个,5 大类) =====
//
// 1. 国风东方系
// 2. 二次元动漫系
// 3. 写实影视系(真人 AI 短剧)
// 4. 科幻奇幻 & 小众艺术
// 5. 3D 渲染类
//
var visualLibrary = []visualLibraryItem{
	// 1. 国风东方系
	{ID: "ink_chinese", Name: "水墨国风", Category: "2d", Desc: "水墨晕染、宣纸留白", PromptHint: "Chinese ink wash painting, calligraphy brush strokes, rice paper texture, negative space, wuxia xianxia ancient Chinese aesthetic", Preview: "/assets/img/library/ink_chinese.png"},
	{ID: "new_chinese", Name: "新中式厚涂", Category: "2d", Desc: "精致人像、低饱和高级感", PromptHint: "modern Chinese thick painting style, refined characters, low saturation muted tones, high-end gufeng sweet romance drama aesthetic", Preview: "/assets/img/library/new_chinese.png"},
	{ID: "dunhuang", Name: "敦煌壁画风", Category: "2d", Desc: "矿物颜料、鎏金纹理", PromptHint: "Dunhuang Mogao cave mural style, mineral pigments, gilded gold leaf texture, mythological epic goddess aesthetic", Preview: "/assets/img/library/dunhuang.png"},
	{ID: "ukiyo_e", Name: "浮世绘国风", Category: "2d", Desc: "线条浓烈,志怪、民间传说", PromptHint: "Ukiyo-e woodblock print style with Chinese folklore elements, bold outlines, flat colors, Chinese monsters and folk tales short drama", Preview: "/assets/img/library/ukiyo_e.png"},
	{ID: "chinese_3d", Name: "国漫 3D", Category: "3d", Desc: "国内动漫渲染", PromptHint: "Chinese donghua 3D render style, wuxia xianxia fantasy, modern Chinese animation cinematic look, current popular manju drama trend", Preview: "/assets/img/library/chinese_3d.png"},
	// 2. 二次元动漫系
	{ID: "cel_anime", Name: "日系赛璐璐", Category: "2d", Desc: "新海诚/吉卜力通透光影", PromptHint: "Studio Ghibli Makoto Shinkai anime cel style, translucent lighting, vibrant pastel colors, youth sweet romance healing", Preview: "/assets/img/library/cel_anime.png"},
	{ID: "shoujo_manga", Name: "少女漫", Category: "2d", Desc: "柔焦、浅色系", PromptHint: "shoujo manga style, soft focus bokeh, pastel color palette, romantic sweet drama aesthetic", Preview: "/assets/img/library/shoujo_manga.png"},
	{ID: "shounen_manga", Name: "少年漫", Category: "2d", Desc: "粗轮廓、高对比热血", PromptHint: "shounen manga style, bold outlines, high contrast, hot-blooded battle shonen anime aesthetic", Preview: "/assets/img/library/shounen_manga.png"},
	{ID: "korean_manhwa", Name: "韩漫厚涂", Category: "2d", Desc: "精致五官、氛围感阴影", PromptHint: "Korean manhwa webtoon thick painting style, refined facial features, atmospheric shadows, modern Korean romance revenge drama aesthetic, top manju traffic style", Preview: "/assets/img/library/korean_manhwa.png"},
	{ID: "chibi", Name: "Q 版卡通", Category: "2d", Desc: "头大身小,搞笑轻喜剧", PromptHint: "chibi Q-version cartoon style, big head small body, comedic lighthearted children's short drama", Preview: "/assets/img/library/chibi.png"},
	// 3. 写实影视系(真人)
	{ID: "cinema_real", Name: "电影写实风", Category: "真人", Desc: "电影布光,都市悬疑", PromptHint: "photorealistic cinematic style, dramatic movie lighting, urban suspense realistic family drama", Preview: "/assets/img/library/cinema_real.png"},
	{ID: "hk_film", Name: "港风复古胶片", Category: "真人", Desc: "暖黄颗粒,年代爱情", PromptHint: "1980s Hong Kong vintage film grain, warm amber tones, retro cinematic, period love story jianghu drama aesthetic", Preview: "/assets/img/library/hk_film.png"},
	{ID: "wkw_style", Name: "王家卫文艺风", Category: "真人", Desc: "柔光慢镜头,情绪向", PromptHint: "Wong Kar-wai art film style, soft diffused lighting, slow motion blur, emotional indie romantic drama aesthetic", Preview: "/assets/img/library/wkw_style.png"},
	{ID: "film_photo", Name: "胶片写真风", Category: "真人", Desc: "柔和自然光,日常向", PromptHint: "analog film photography style, soft natural lighting, modern sweet romance slice of life, gentle warm tones", Preview: "/assets/img/library/film_photo.png"},
	{ID: "film_noir", Name: "黑白 Noir", Category: "真人", Desc: "高对比硬阴影,悬疑犯罪", PromptHint: "black and white film noir, high contrast hard shadows, suspense crime mystery drama aesthetic", Preview: "/assets/img/library/film_noir.png"},
	// 4. 科幻奇幻 & 小众艺术
	{ID: "cyberpunk", Name: "赛博朋克", Category: "3d", Desc: "霓虹雨夜金属质感", PromptHint: "cyberpunk neon rain night, metallic future city, robotic mecha aesthetic, futuristic urban drama", Preview: "/assets/img/library/cyberpunk.png"},
	{ID: "wasteland", Name: "末世废土", Category: "3d", Desc: "灰黄破败,末日求生", PromptHint: "post-apocalyptic wasteland, desaturated gray yellow ruins, survival horror short drama aesthetic", Preview: "/assets/img/library/wasteland.png"},
	{ID: "american_comic", Name: "美漫美式漫画", Category: "2d", Desc: "粗黑轮廓网点,超级英雄", PromptHint: "American comic book style, bold black outlines, halftone dots, superhero action comedy aesthetic", Preview: "/assets/img/library/american_comic.png"},
	{ID: "claymation", Name: "黏土定格动画", Category: "3d", Desc: "手工质感,治愈轻喜剧", PromptHint: "claymation stop motion style, handcrafted texture, cozy lighthearted comedy healing drama", Preview: "/assets/img/library/claymation.png"},
	{ID: "watercolor", Name: "水彩手绘", Category: "2d", Desc: "通透水彩笔触,童话奇幻", PromptHint: "watercolor hand-painted illustration style, translucent brushstrokes, fairy tale fantasy short story aesthetic", Preview: "/assets/img/library/watercolor.png"},
	{ID: "pixel_art", Name: "像素风", Category: "2d", Desc: "复古游戏质感,穿越", PromptHint: "retro pixel art game style, 16-bit nostalgic gaming aesthetic, time-travel isekai short drama", Preview: "/assets/img/library/pixel_art.png"},
	// 5. 3D 渲染类
	{ID: "pixar_3d", Name: "皮克斯卡通 3D", Category: "3d", Desc: "柔和材质,家庭喜剧", PromptHint: "Pixar Disney 3D cartoon render style, soft materials, family comedy light fantasy aesthetic", Preview: "/assets/img/library/pixar_3d.png"},
	{ID: "realistic_3d", Name: "写实 3D CG", Category: "3d", Desc: "超精细渲染,玄幻大片", PromptHint: "hyper-detailed photorealistic 3D CG render, cinematic VFX fantasy epic short film", Preview: "/assets/img/library/realistic_3d.png"},
	{ID: "blindbox_3d", Name: "潮玩盲盒 3D", Category: "3d", Desc: "圆润 IP 感,轻松向", PromptHint: "blind box collectible figurine 3D style, rounded IP toy design, casual cheerful short video aesthetic", Preview: "/assets/img/library/blindbox_3d.png"},
}

// ===== 叙事题材库(8 个) =====
// 用于 AI 生成剧本时附加调性描述,影响剧本的剧情走向、冲突密度、情绪基调。
var genreLibrary = []genreLibraryItem{
	{ID: "urban_power", Name: "都市爽文风", Desc: "真假千金、豪门逆袭、重生复仇,冲突密集反转快,短视频平台最大流量赛道", PromptHint: "modern urban power-fantasy drama, fake heiress豪门逆袭, reborn revenge, domineering CEO sweet romance, dense conflicts and fast plot twists, top traffic genre on short-video platforms"},
	{ID: "ancient_xianxia", Name: "古风权谋/仙侠", Desc: "宫斗、穿越、修仙、废柴逆袭,大女主、虐恋", PromptHint: "ancient Chinese palace intrigue and xianxia fantasy, time-travel cultivation, underdog reversal, strong female lead, bittersweet romance"},
	{ID: "retro_farming", Name: "年代种田风", Desc: "七八十年代重生、乡土生活,烟火气,慢治愈", PromptHint: "1970s-80s Chinese rural rebirth and farming life, earthy nostalgia, slow healing slice-of-life aesthetic"},
	{ID: "mystery_twist", Name: "悬疑反转风", Desc: "探案、民间怪谈、密室,每集结尾留钩子", PromptHint: "mystery thriller with plot twists, detective cases, Chinese folk supernatural tales, locked room puzzles, every episode ends with a cliffhanger hook"},
	{ID: "silly_comedy", Name: "沙雕搞笑风", Desc: "穿书反套路、无厘头喜剧,轻松解压", PromptHint: "silly slapstick comedy, transmigration into novel tropes, absurdist humor, lighthearted stress-relief comedy"},
	{ID: "family_drama", Name: "现实家庭向", Desc: "婆媳矛盾、原生家庭、宝妈觉醒,强共情", PromptHint: "realistic Chinese family drama, mother-in-law conflicts, family-of-origin trauma, stay-at-home mom awakening, strong emotional resonance"},
	{ID: "fantasy_apocalypse", Name: "奇幻末日风", Desc: "异能、末世生存、吸血鬼狼人等海外短剧题材", PromptHint: "supernatural apocalyptic fantasy, superpowers, end-of-world survival, vampires and werewolves, overseas short-drama genre"},
	{ID: "school_youth", Name: "青春校园", Desc: "暗恋、校园成长,清新恋爱", PromptHint: "Chinese youth school campus romance, secret crushes, coming-of-age growth, fresh innocent love story"},
}

// ===== 比例选项(6 个常见) =====
// Agnes 视频支持 720x1280 / 1280x720,图片支持 1024x1024 / 1024x1792 / 1792x1024 等。
// 比例定义在 internal/types/aspect.go,跨包共享。
var aspectRatios = types.AspectRatios

// ===== 剧集数选项 =====
// 自定义在前端用 input 实现,这里只给预设。
var episodeCountOptions = []episodeCountOption{
	{Value: 5, Label: "05 集"},
	{Value: 10, Label: "10 集"},
	{Value: 20, Label: "20 集"},
	{Value: 50, Label: "50 集"},
	{Value: 80, Label: "80 集"},
	{Value: 100, Label: "100 集"},
}

// resolveAspectSize 把 aspect_ratio id 解析成 image/video size。
// 找不到时返回默认 1024x1024 / 1280x720。
func resolveAspectSize(ar string) (imgSize, vidSize string) {
	imgSize, vidSize = "1024x1024", "1280x720"
	for _, o := range aspectRatios {
		if o.ID == ar {
			return o.ImageSize, o.VideoSize
		}
	}
	return
}

// resolveVisualHint 根据 visual_style id 拿到 PromptHint,用于强化生成 prompt。
// 找不到时返回空串。
func resolveVisualHint(id string) string {
	for _, v := range visualLibrary {
		if v.ID == id {
			return v.PromptHint
		}
	}
	return ""
}

// resolveGenre 根据 genre_id 拿到 PromptHint 和 Name。
func resolveGenre(id string) (hint, name string) {
	for _, g := range genreLibrary {
		if g.ID == id {
			return g.PromptHint, g.Name
		}
	}
	return "", ""
}

// listStyleLibrary 一次返回视觉画风库 / 叙事题材 / 比例 / 剧集数所有数据。
func (s *Server) listStyleLibrary(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"visual_styles":  visualLibrary,
		"genres":         genreLibrary,
		"aspect_ratios":  aspectRatios,
		"episode_counts": episodeCountOptions,
	})
}
