package main

import (
	"fmt"
	"strconv"
	"syscall/js"
	"time"

	lyrics "github.com/xiaowumin-mark/visual-lyric-core/lyric"
)

var (
	document js.Value = js.Global().Get("document")
	location js.Value = js.Global().Get("location")
	gsap     js.Value = js.Global().Get("gsap")
	//lrcText  string   = js.Global().Get("lrcText").String()

	window js.Value = js.Global()
	// console       js.Value = js.Global().Get("console")
	alert js.Value = js.Global().Get("alert")
	// setTimeout    js.Value = js.Global().Get("setTimeout")
	// setInterval   js.Value = js.Global().Get("setInterval")
	// clearTimeout  js.Value = js.Global().Get("clearTimeout")
	// clearInterval js.Value = js.Global().Get("clearInterval")
	bglw        = 100
	fadeRatio   = 0.5
	bgfadeRatio = fadeRatio * 2
	dialogEle   = document.Call("getElementById", "dialog")
	musicInput  = document.Call("getElementById", "ex_music")
	coverInput  = document.Call("getElementById", "ex_cover")
	ttmlInput   = document.Call("getElementById", "ex_ttml")
	compBtn     = document.Call("getElementById", "comp_ex")
)
var fr js.Func
var previousIndex = make([]int, 0)
var nowPlayingIndex = make([]int, 0)
var innerHeightShowItemNum = 10
var noePlayingOne = -1
var trans = "background 0.7s, filter 0.5s, opacity 0.5s"
var audio js.Value
var hasScrolledInRemove bool
var c = make(chan struct{}, 0)

func main() {
	//const urlParams = new URLSearchParams(window.location.search);
	//const musicName = urlParams.get("m") || "ME!";
	//const musicType = urlParams.get("t") || "mp3";
	//lrcView := document.Call("createElement", "div")
	//lrcView.Get("classList").Call("add", "lyric")
	//lrcView.Get("classList").Call("add", "scrollbar-hidden")
	lrcView := document.Call("getElementById", "lrcView")

	//document.Get("body").Call("append", lrcView)
	audio = js.Global().Get("Audio").New()
	var urlParams = js.Global().Get("URLSearchParams").New(location.Get("search"))
	musicName := urlParams.Call("get", "m").String()
	musicType := urlParams.Call("get", "t").String()
	exFile := urlParams.Call("get", "e").String()
	if musicName == js.Null().String() {
		musicName = "ME!"
	}
	if musicType == js.Null().String() {
		musicType = "mp3"
	}
	if exFile == js.Null().String() {
		exFile = "false"
	}
	bgsrc := ""
	lrcText := ""
	compBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if exFile == "false" {
			return nil
		}
		if musicInput.Get("files").Length() == 0 {
			window.Call("alert", "请选择歌词文件")
			return nil
		}
		if coverInput.Get("files").Length() == 0 {
			window.Call("alert", "请选择封面文件")
			return nil
		}
		if ttmlInput.Get("files").Length() == 0 {
			window.Call("alert", "请选择ttml文件")
			return nil
		}
		// 将音频转换为在线地址
		//reader := js.Global().Get("FileReader").New()
		//reader.Call("addEventListener", "load", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		//	var base64 = reader.Get("result").String()
		//	fmt.Println(base64)
		//	audio.Set("src", base64)
		//	i++
		//	if i == 3 {
		//
		//		go start(lrcText, bgsrc, lrcView)
		//	}
		//	return nil
		//}))
		//reader.Call("readAsDataURL", musicInput.Get("files").Index(0))

		// 读取cover
		//reader = js.Global().Get("FileReader").New()
		//reader.Call("addEventListener", "load", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		//	var base64 = reader.Get("result").String()
		//	fmt.Println(base64)
		//	bgsrc = base64
		//	i++
		//	if i == 3 {
		//
		//		go start(lrcText, bgsrc, lrcView)
		//	}
		//
		//	return nil
		//}))
		//reader.Call("readAsDataURL", coverInput.Get("files").Index(0))

		url := window.Get("URL").Call("createObjectURL", musicInput.Get("files").Index(0))
		audio.Set("src", url)
		url = window.Get("URL").Call("createObjectURL", coverInput.Get("files").Index(0))
		bgsrc = url.String()

		// 读取ttml文本
		reader := js.Global().Get("FileReader").New()
		reader.Call("addEventListener", "load", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			var base64 = reader.Get("result").String()
			//
			lrcText = base64

			go start(lrcText, bgsrc, lrcView)
			return nil
		}))

		reader.Call("readAsText", ttmlInput.Get("files").Index(0))

		dialogEle.Get("classList").Call("remove", "show")
		return nil
	}))
	if location.Get("hash").String() != "" {
		getFile := getLrcText(location.Get("hash").String()[1:], "json")
		audio.Set("src", getFile.Get("song").String())
		bgsrc = getFile.Get("meta").Get("albumImgSrc").String()
		lrcText = getLrcText(getFile.Get("meta").Get("lyrics").String(), "text").String()
		go start(lrcText, bgsrc, lrcView)
	} else {
		if exFile != "false" {
			// 让用户选中歌曲，歌词，封面
			dialogEle.Get("classList").Call("add", "show")
		} else {
			audio.Set("src", "/music/"+musicName+"."+musicType)
			audio.Call("load")
			bgsrc = "/music/" + musicName + ".png"
			lrcText = getLrcText("/music/"+musicName+".ttml", "text").String()

			go start(lrcText, bgsrc, lrcView)
		}
	}

	<-c
}

func start(lrcText string, bgsrc string, lrcView js.Value) {
	// 获取所有包含data-muaic-background属性的img元素
	backgroundImages := document.Call("querySelectorAll", "[data-muaic-background]")
	for i := 0; i < backgroundImages.Length(); i++ {

		img := backgroundImages.Index(i)
		img.Set("src", bgsrc)
	}

	println(lrcText)
	vld, err := ParseTTML(lrcText, lrcView)
	if err != nil {
		panic(err)
	}
	gd(0, vld, true)
	for _, item := range vld.Contents {
		// 添加点击事件
		item.Ele.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			setCurrentTime(audio, item.Primary.Begin+40*time.Millisecond)
			audio.Call("play")
			return nil
		}))
		//GsetTimeout(func() {
		//	item.Ele.Get("style").Set("transition", trans)
		//	item.Ele.Get("style").Set("opacity", "1")
		//}, 100*time.Millisecond)

	}

	document.Call("addEventListener", "keydown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// 空格
		if args[0].Get("code").String() == "Space" {
			if audio.Get("paused").Bool() {
				audio.Call("play")
				//playLrc(vld)
			} else {
				audio.Call("pause")
				//pauseLrc(vld)
			}
		} else if args[0].Get("code").String() == "ArrowRight" {
			// 快进
			audio.Set("currentTime", audio.Get("currentTime").Float()+5)
		} else if args[0].Get("code").String() == "ArrowLeft" {
			// 快退
			audio.Set("currentTime", audio.Get("currentTime").Float()-5)
		} else if args[0].Get("code").String() == "ArrowUp" {
			// 音量加
			if audio.Get("volume").Float() < 1 {
				audio.Set("volume", audio.Get("volume").Float()+0.1)
			}

		} else if args[0].Get("code").String() == "ArrowDown" {
			// 音量减
			if audio.Get("volume").Float() > 0 {
				audio.Set("volume", audio.Get("volume").Float()-0.1)
			}
		}
		return nil
	}))
	document.Call("addEventListener", "dblclick", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if audio.Get("paused").Bool() {
			audio.Call("play")
			//playLrc(vld)
		} else {
			audio.Call("pause")
			//pauseLrc(vld)
		}
		return nil
	}))
	audio.Call("addEventListener", "play", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		playLrc(vld)

		js.Global().Call("requestAnimationFrame", fr)
		return nil
	}))
	audio.Call("addEventListener", "pause", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		pauseLrc(vld)
		return nil
	}))

	// 窗口缩放事件
	window.Call("addEventListener", "resize", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		initLrcBackground(vld)
		rPosition(vld)
		return nil
	}))
	//GsetInterval(
	//	func() {
	//		if audio.Get("paused").Bool() {
	//			return
	//		}
	//		currentIndex := make([]int, 0)
	//		var currentTime time.Duration = getCurrentTime(audio)
	//
	//		for i := 0; i < len(vld.Contents); i++ {
	//			v := vld.Contents[i]
	//			if currentTime >= v.Primary.Begin && currentTime <= v.Primary.End {
	//				currentIndex = append(currentIndex, i)
	//			}
	//		}
	//
	//		// 只有当 currentIndex 发生变化时才触发歌词变化处理
	//		if !every(currentIndex, previousIndex) {
	//			handleLyricsChange(vld, currentIndex)
	//			previousIndex = make([]int, len(currentIndex)) // 深拷贝
	//			copy(previousIndex, currentIndex)
	//		}
	//	},
	//	50*time.Millisecond,
	//)

	fr = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if audio.Get("paused").Bool() {
			return nil
		}
		currentIndex := make([]int, 0)
		var currentTime time.Duration = getCurrentTime(audio)

		for i := 0; i < len(vld.Contents); i++ {
			v := vld.Contents[i]
			if currentTime >= v.Primary.Begin && currentTime <= v.Primary.End {
				currentIndex = append(currentIndex, i)
			}
		}

		// 只有当 currentIndex 发生变化时才触发歌词变化处理
		if !every(currentIndex, previousIndex) {
			handleLyricsChange(vld, currentIndex)
			previousIndex = make([]int, len(currentIndex)) // 深拷贝
			copy(previousIndex, currentIndex)
		}
		js.Global().Call("requestAnimationFrame", fr)
		return nil
	})
	js.Global().Call("requestAnimationFrame", fr)

	<-c
}

func pauseLrc(lrcs *lyrics.Lyrics) {
	for _, index := range nowPlayingIndex {
		item := lrcs.Contents[index]
		for _, item1 := range item.Primary.Blocks {

			//item1.TextUpAnimation.Call("pause")
			if !item1.TextUpAnimation.IsUndefined() && !item1.TextUpAnimation.IsNull() && item1.TextUpAnimation.Type() == js.TypeObject {
				item1.TextUpAnimation.Call("pause")
			}

			for _, item_animate := range item1.Animation {
				if !item_animate.IsUndefined() && !item_animate.IsNull() && item_animate.Type() == js.TypeObject {
					item_animate.Call("pause")
				}
			}
			//item1.GsapAnimation.Call("pause")

			// item1.HighLightAnimations []
			for _, item2 := range item1.HighLightAnimations {
				item2.Call("pause")
			}
		}

		for _, item2 := range item.Backgrounds {
			for _, item3 := range item2.Blocks {
				if !item3.TextUpAnimation.IsUndefined() && !item3.TextUpAnimation.IsNull() && item3.TextUpAnimation.Type() == js.TypeObject {
					item3.TextUpAnimation.Call("pause")
				}
				//if !item3.GsapAnimation.IsUndefined() && !item3.GsapAnimation.IsNull() && item3.GsapAnimation.Type() == js.TypeObject {
				//	item3.GsapAnimation.Call("pause")
				//}
				for _, item_animate := range item3.Animation {
					if !item_animate.IsUndefined() && !item_animate.IsNull() && item_animate.Type() == js.TypeObject {
						item_animate.Call("pause")
					}
				}
				for _, item4 := range item3.HighLightAnimations {
					item4.Call("pause")
				}
			}
		}
	}
}

func playLrc(lrcs *lyrics.Lyrics) {
	for _, index := range nowPlayingIndex {
		item := lrcs.Contents[index]
		for _, item1 := range item.Primary.Blocks {
			for _, item_animate := range item1.Animation {
				if !item_animate.IsUndefined() && !item_animate.IsNull() && item_animate.Type() == js.TypeObject {
					item_animate.Call("play")
				}
			}
			if !item1.TextUpAnimation.IsUndefined() && !item1.TextUpAnimation.IsNull() && item1.TextUpAnimation.Get("overallProgress").Float() != float64(1) {
				item1.TextUpAnimation.Call("play")
			}
			// item1.HighLightAnimations []
			for _, item2 := range item1.HighLightAnimations {
				if !item2.IsUndefined() && !item2.IsNull() && item2.Type() == js.TypeObject && item2.Get("overallProgress").Float() != float64(1) {
					item2.Call("play")
				}
			}
		}
		for _, item2 := range item.Backgrounds {
			for _, item3 := range item2.Blocks {
				if !item3.TextUpAnimation.IsUndefined() && !item3.TextUpAnimation.IsNull() {
					item3.TextUpAnimation.Call("play")
				}
				for _, item_animate := range item3.Animation {
					if !item_animate.IsUndefined() && !item_animate.IsNull() && item_animate.Type() == js.TypeObject {
						item_animate.Call("play")
					}
				}

				for _, item4 := range item3.HighLightAnimations {
					if item4.Get("overallProgress").Float() != float64(1) {
						item4.Call("play")
					}
				}
			}
		}
	}
}

func getLrcText(path string, Type string) js.Value {
	done := make(chan struct{})
	var responseData js.Value

	// 使用 JavaScript 的 fetch API 发送请求
	promise := js.Global().Call("fetch", path)
	promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		response := args[0]
		response.Call(Type).Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			responseData = args[0]
			close(done) // 通知请求完成
			return nil
		}))
		return nil
	})).Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		close(done) // 即使出错也要通知完成
		return nil
	}))

	// 等待请求完成
	<-done

	// 返回响应数据
	return responseData
}
func handleLyricsChange(lrcs *lyrics.Lyrics, highlightedLyric []int) {
	addedLyrics := filterCurrentIndex(highlightedLyric, previousIndex)
	removedLyrics := filterCurrentIndex(previousIndex, highlightedLyric)
	for j := 0; j < len(removedLyrics); j++ {

		index := removedLyrics[j].(int)
		removeLyric(index, lrcs)
	}
	for i := 0; i < len(addedLyrics); i++ {

		index := addedLyrics[i].(int)

		addLyric(index, lrcs)
		for _, da := range lrcs.Contents[index].Primary.Blocks {
			print(da.Text)
		}
		println("")
	}

}

//func gd(i int, lrc *lyrics.Lyrics) {
//	//lrc.Contents[i].Ele.Get("style").Set("background-color", "rgba(255, 255, 255, 0.45)")
//	for j, item := range lrc.Contents {
//		n := mathAbs(i - 3 - j)
//		ah := time.Duration(n*70-n*40) * time.Millisecond
//
//		h := getTopHeight(lrc, i, item.Index)
//		GsetTimeout(func() {
//			//item.Ele.Get("style").Call("setProperty", "--top", getTopHeight(lrc, i, j))
//			//item.Ele.Get("style").Call("setProperty", "transform", "translateY("+strconv.Itoa(h)+"px)")
//			item.Ele.Get("style").Call("setProperty", "transform", "translateY("+strconv.Itoa(h)+"px)")
//			// body innerHtml
//		}, ah)
//	}
//}

func gd(currentIndex int, lrc *lyrics.Lyrics, init bool) {
	fmt.Println("滚动", currentIndex)
	for index, item := range lrc.Contents {

		item.Ele.Get("style").Set("filter", fmt.Sprintf("blur(%vpx)", mathAbs(index-currentIndex)))
		//if index == currentIndex {
		//	it := lrc.Contents[index]
		//	if len(it.Backgrounds) != 0 {
		//		it.BackgroundsEle.Get("style").Set("display", "block")
		//		GsetTimeout(func() {
		//			it.BackgroundsEle.Get("classList").Call("add", "bgShow")
		//		}, 10*time.Millisecond)
		//	}
		//}

		// 计算当前歌词到目标歌词的累计高度
		top := getTopHeight(lrc, currentIndex, index, -1)

		//top := js.Global().Call("getTopHeight", currentIndex, index).Int()
		item.Primary.Position = top
		// 添加弹性动画延迟（与 JS 的 elastic.out 对齐）
		var delay time.Duration
		var duration = 1.2
		if init {
			delay = 0
			duration = 0
		} else {
			delay = time.Duration(mathAbs(currentIndex-index-3)*50) * time.Millisecond

		}
		//if !item.ScrollAnimation.IsUndefined() && !item.ScrollAnimation.IsNull() {
		//	item.ScrollAnimation.Call("pause")
		//}
		rn := currentIndex - index
		if rn > -innerHeightShowItemNum && rn < innerHeightShowItemNum {

			//item.ScrollAnimation = item.Ele.Call("animate",
			//	[]interface{}{
			//		map[string]interface{}{
			//			"transform": fmt.Sprintf("translateY(%dpx)", top),
			//		},
			//	},
			//	map[string]interface{}{
			//		"duration": duration,
			//		"easing":   "cubic-bezier(0.19, 1, 0.22, 1)",
			//		"delay":    delay.Milliseconds(),
			//		"fill":     "forwards",
			//	},
			//)
			GsetTimeout(func() {
				gsap.Call("to", item.Ele, map[string]interface{}{
					"y":        fmt.Sprintf("%dpx", top),
					"duration": duration,
					"ease":     "elastic.out(1, 1.35)",
				})
			}, delay)
			offset := 0
			for bi, Bitem := range item.Backgrounds {
				//log.Println("背景", bgindex)
				/*btop := getTopHeight(lrc, currentIndex, index, bgindex)
				Bitem.Position = btop
				gsap.Call("to", Bitem.Ele, map[string]interface{}{
					"y":        fmt.Sprintf("%dpx", btop),
					"duration": 0,
					//"ease":     "elastic.out(1, 1.35)",
				})*/
				if bi == 0 {
					offset += item.Ele.Get("offsetHeight").Int()
				} else {
					offset += item.Backgrounds[bi-1].Ele.Get("offsetHeight").Int()
				}

				Bitem.Position = top + offset
				GsetTimeout(func() {
					gsap.Call("to", Bitem.Ele, map[string]interface{}{
						"y":        fmt.Sprintf("%dpx", Bitem.Position),
						"duration": 1.2,
						"ease":     "elastic.out(1, 1.35)",
					})
				}, delay)
				//Bitem.Ele.Get("style").Set("top", fmt.Sprintf("%dpx", Bitem.Position))
			}
		} else {
			//item.Ele.Get("style").Set("transition", "none")
			//item.Ele.Get("style").Call("setProperty", "--top", fmt.Sprintf("%dpx", item.Position))
			//ainat := item.Ele.Call("animate",
			//	[]interface{}{
			//		map[string]interface{}{
			//			"transform": fmt.Sprintf("translateY(%dpx)", top),
			//		},
			//	},
			//	map[string]interface{}{
			//		"duration": 0,
			//		"easing":   "cubic-bezier(0.19, 1, 0.22, 1)",
			//		"delay":    0,
			//		"fill":     "forwards",
			//	},
			//)
			//if !item.ScrollAnimation.IsUndefined() && !item.ScrollAnimation.IsNull() {
			//	item.ScrollAnimation.Call("cancel")
			//}
			//item.ScrollAnimation = ainat
			gsap.Call("to", item.Ele, map[string]interface{}{
				"y":        fmt.Sprintf("%dpx", top),
				"duration": 0,
				"ease":     "elastic.out(1, 1.35)",
			})
		}
		//item.Ele.Get("style").Call("setProperty", "--top", fmt.Sprintf("%dpx", item.Position))
		// item.ScrollAnimation 不为空的话

	}
}

func rPosition(lrc *lyrics.Lyrics) {
	in := bubbleSort(nowPlayingIndex)
	i := 0
	if len(in) > 0 {
		i = in[0]
	}
	for index, item := range lrc.Contents {
		top := getTopHeight(lrc, i, index, -1)
		for bgindex, Bitem := range item.Backgrounds {
			Bitem.Position = getTopHeight(lrc, i, index, bgindex)
			gsap.Call("to", item.Ele, map[string]interface{}{
				"y":        fmt.Sprintf("%dpx", Bitem.Position),
				"duration": 0,
				//"ease":     "elastic.out(1, 1.35)",
			})
		}
		item.Ele.Get("style").Call("setProperty", "transform", "translate(0px,"+strconv.Itoa(top)+"px)")
		item.Primary.Position = top
	}
}
func getCurrentTime(audio js.Value) time.Duration {
	seconds := audio.Get("currentTime").Float()
	return time.Duration(seconds * float64(time.Second))
}

func setCurrentTime(audio js.Value, currentTime time.Duration) {
	audio.Set("currentTime", currentTime.Seconds())
}

// 每个元素与另一个切片中的元素相等时返回 true，否则返回 false
func every(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// contains 函数用于检查一个元素是否存在于切片中，切片和元素都为 interface{} 类型
func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// filterCurrentIndex 函数根据条件过滤出 currentIndex 中不在 nowPlayingIndex 中的元素
func filterCurrentIndex(currentIndex []int, nowPlayingIndex []int) []interface{} {
	var result []interface{}
	for _, value := range currentIndex {
		// 如果 value 不在 nowPlayingIndex 中，添加到结果切片中
		if !contains(nowPlayingIndex, value) {
			result = append(result, value)
		}
	}
	return result
}

/*func getTopHeight(lrc *lyrics.Lyrics, now, to, bgIndex int) int {
	var res int = 0
	if to > now {
		for i := now; i < to; i++ {
			// 强制布局更新
			h := lrc.Contents[i].Ele.Get("clientHeight").Int()
			res += h + 10 // 行间距为10

			for _, Bitem := range lrc.Contents[i].Backgrounds {
				res += Bitem.Ele.Get("clientHeight").Int()
			}
			if bgIndex >= 0 && i == now {

				//log.Println("计算滚动的时候遇到我了，index", to, "滚向", now)
				//log.Println("背景数量", len(lrc.Contents[i].Backgrounds))
				//log.Println("i", i)
				for _, e := range lrc.Contents[i].Primary.Blocks {
					fmt.Print(e.Text)
				}
				fmt.Println()
				for q := 0; q < bgIndex+1; q++ {
					log.Println("top", q)
					res += lrc.Contents[i].Backgrounds[q].Ele.Get("clientHeight").Int()
				}
			}

		}
	} else {
		for j := now; j > to; j-- {
			// 强制布局更新
			h := lrc.Contents[j-1].Ele.Get("clientHeight").Int()
			res -= h + 10 // 行间距为10
			for _, Bitem := range lrc.Contents[j-1].Backgrounds {
				res -= Bitem.Ele.Get("clientHeight").Int()
			}
			if bgIndex >= 0 && j-1 == now {
				for q := 0; q < bgIndex+1; q++ {

					res -= lrc.Contents[j-1].Backgrounds[q].Ele.Get("clientHeight").Int()
				}
			}
		}
	}

	return res + 400// 偏移400
}*/

func getTopHeight(lrc *lyrics.Lyrics, now, to, bgIndex int) int {
	res := 0

	if to > now {
		for i := now; i < to; i++ {
			h := lrc.Contents[i].Ele.Get("clientHeight").Int()
			res += h
			if lrc.Contents[i].ShowBackgrounds {

				for _, Bitem := range lrc.Contents[i].Backgrounds {
					res += Bitem.Ele.Get("clientHeight").Int()
				}
			}
		}
		// 单独处理“滚到的是自己且是某个背景歌词”
		/*if bgIndex >= 0 && to < len(lrc.Contents) {
			for q := 0; q <= bgIndex && q < len(lrc.Contents[to].Backgrounds); q++ {
				res += lrc.Contents[to].Backgrounds[q].Ele.Get("clientHeight").Int()
			}
		}*/
	} else {
		for j := now; j > to; j-- {
			h := lrc.Contents[j-1].Ele.Get("clientHeight").Int()
			res -= h
			if lrc.Contents[j-1].ShowBackgrounds {
				for _, Bitem := range lrc.Contents[j-1].Backgrounds {
					res -= Bitem.Ele.Get("clientHeight").Int()
				}
			}
		}
		// 单独处理“滚到的是自己且是某个背景歌词”
		/*if bgIndex >= 0 && to >= 0 {
			for q := 0; q <= bgIndex && q < len(lrc.Contents[to].Backgrounds); q++ {
				res -= lrc.Contents[to].Backgrounds[q].Ele.Get("clientHeight").Int()
			}
		}*/
	}

	return res + 200 // 偏移
}

//func getTopHeight(lrc *lyrics.Lyrics, now, to int) int {
//	var res int = 200
//	// 确保 now 和 to 的索引有效
//	if now < 0 || to < 0 || now >= len(lrc.Contents) || to >= len(lrc.Contents) {
//		return 0
//	}
//	// 统一使用正向累计逻辑
//	step := 1
//	if to < now {
//		step = -1
//	}
//	for i := now; i != to; i += step {
//		height := lrc.Contents[i].Ele.Get("offsetHeight").Int()
//		res += height
//	}
//	return res
//}

func mathAbs(a int) int {
	if a < 0 {
		return -a
	} else {
		return a
	}
}

func bubbleSort(arr []int) []int {
	length := len(arr) //数据总长度（个数）
	for i := 0; i < length; i++ {
		for j := 0; j < length-1-i; j++ {
			if arr[j] > arr[j+1] { //和相邻的比
				arr[j], arr[j+1] = arr[j+1], arr[j] //对换位置
			}
		}
	}
	return arr
}

func addIndex(i int) {
	// 如果i在nowPlayingIndex中，则什么也不做
	if contains(nowPlayingIndex, i) {
		return
	}
	nowPlayingIndex = append(nowPlayingIndex, i)
}

func removeIndex(i int) {
	// 如果i不在nowPlayingIndex中 则什么也不做
	if !contains(nowPlayingIndex, i) {
		return
	}
	// 创建一个新的切片，将nowPlayingIndex中的元素复制到新切片中
	newSlice := make([]int, len(nowPlayingIndex))
	copy(newSlice, nowPlayingIndex)
	// 遍历新切片，将每个元素与i进行比较，如果相等则将其从新切片中删除
	for j, v := range newSlice {
		if v == i {
			newSlice = append(newSlice[:j], newSlice[j+1:]...)
		}
	}
	nowPlayingIndex = newSlice

}
