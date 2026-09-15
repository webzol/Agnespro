// aspect.go - 视频/图片比例数据(共享给 api + pipeline 两个包)
package types

// AspectRatioOption:一个比例选项
type AspectRatioOption struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	Ratio     string `json:"ratio"`
	ImageSize string `json:"image_size"`
	VideoSize string `json:"video_size"`
	Hint      string `json:"hint"`
}

// 6 个常见比例
var AspectRatios = []AspectRatioOption{
	{ID: "9:16", Label: "9:16 竖屏", Ratio: "9:16", ImageSize: "1024x1792", VideoSize: "720x1280", Hint: "抖音/快手/小红书"},
	{ID: "16:9", Label: "16:9 横屏", Ratio: "16:9", ImageSize: "1792x1024", VideoSize: "1280x720", Hint: "YouTube/B站/西瓜"},
	{ID: "1:1", Label: "1:1 方形", Ratio: "1:1", ImageSize: "1024x1024", VideoSize: "1024x1024", Hint: "朋友圈/Instagram"},
	{ID: "3:4", Label: "3:4 轻竖屏", Ratio: "3:4", ImageSize: "896x1152", VideoSize: "720x960", Hint: "小红书/微博"},
	{ID: "4:3", Label: "4:3 传统", Ratio: "4:3", ImageSize: "1152x896", VideoSize: "960x720", Hint: "老电视/复古感"},
	{ID: "21:9", Label: "21:9 电影宽屏", Ratio: "21:9", ImageSize: "1792x768", VideoSize: "1280x544", Hint: "影院质感/大片感"},
}

// ResolveAspectSize 根据 aspect_ratio id 解析成 image/video size。
func ResolveAspectSize(ar string) (imgSize, vidSize string) {
	imgSize, vidSize = "1024x1024", "1280x720"
	for _, o := range AspectRatios {
		if o.ID == ar {
			return o.ImageSize, o.VideoSize
		}
	}
	return
}
