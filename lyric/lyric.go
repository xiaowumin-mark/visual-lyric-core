package lyrics

import (
	"syscall/js"
	"time"
)

// Lyrics 结构体表示整个歌词文件
type Lyrics struct {
	Meta       Meta              // 元数据
	Contents   []*Content        // 歌词内容
	References map[string]string // 外部引用（如音乐平台ID）
	Authors    []string          // 作者信息
}

// NewEmptyLyrics 创建一个空的 Lyrics 实例
func NewEmptyLyrics() *Lyrics {
	return &Lyrics{
		Meta:       Meta{},
		Contents:   []*Content{},
		References: make(map[string]string),
		Authors:    []string{},
	}
}

// Meta 结构体表示歌词的元数据
type Meta struct {
	MusicName []string // 音乐名称
	Artist    []string // 艺术家
	Album     []string // 专辑
	Isrc      []string // ISRC 编码
}

// Content 结构体表示一行歌词的内容
type Content struct {
	Primary         *Line   // 主要歌词行
	Backgrounds     []*Line // 背景歌词行
	BackgroundsEle  js.Value
	Ele             js.Value
	ScrollAnimation js.Value
	Agent           string
	Index           int
	//Position        int
	ShowBackgrounds bool
}

// Line 结构体表示一行歌词
type Line struct {
	Begin        time.Duration     // 开始时间（毫秒）
	End          time.Duration     // 结束时间（毫秒）
	Translates   map[string]string // 翻译内容（语言代码 -> 翻译文本）
	Duet         bool              // 是否为对唱
	Blocks       []*Block          // 歌词块
	WrdList      [][]int
	Ele          js.Value
	TranslateEle js.Value
	Position     int
	Splits       []string
}

// Block 结构体表示歌词中的一个块（单词或空格）
type Block struct {
	Begin                        time.Duration // 开始时间（毫秒）
	End                          time.Duration // 结束时间（毫秒）
	Text                         string        // 文本内容
	Ele                          js.Value
	Animation                    []js.Value
	TextUpAnimation              js.Value
	HighLightAnimations          []js.Value
	HighLightBackgroungAnimation []js.Value
	AheadtrAnimation             js.Value
}

// 定义外部引用来源
const (
	ReferenceSourceQQMusic    = "qqMusic"
	ReferenceSourceSpotify    = "spotify"
	ReferenceSourceAppleMusic = "appleMusic"
	ReferenceSourceNetease    = "netease"
)
