package main

import (
	"fmt"
	"syscall/js"

	lyrics "github.com/xiaowumin-mark/visual-lyric-core/lyric"
)

func removeLyric(index int, lrcs *lyrics.Lyrics) {

	removeIndex(index)
	if len(lrcs.Contents[index].Backgrounds) != 0 {
		lrcs.Contents[index].BackgroundsEle.Get("style").Set("display", "none")
		lrcs.Contents[index].BackgroundsEle.Get("classList").Call("remove", "bgShow")

		for _, item := range lrcs.Contents[index].Backgrounds {
			for _, word := range item.Blocks {

				//gsap.Call("to", word.Ele, map[string]interface{}{
				//	"duration": 0,
				//	"--p":      "-40%",
				//	"--rp":     "0%",
				//})
				gsap.Call("to", word.Ele, map[string]interface{}{
					"duration": 0.2,
					"--rcolor": 0.2,
					"onComplete": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
						word.Ele.Get("style").Call("setProperty", "--p", "0px")
						word.Ele.Get("style").Call("setProperty", "--rp", "0%")
						word.Ele.Get("style").Call("setProperty", "--rcolor", "1")
						return nil
					}),
				})
				word.GsapAnimation.Call("kill")
				word.GsapAnimation = js.Null()
			}
		}
	}
	words := lrcs.Contents[index].Primary.Blocks

	for _, word := range words {
		if word.Begin == 0 && word.End == 0 {
			continue
		}

		var bgChildren js.Value = word.Ele.Get("children")
		//word.Ele.Get("style").Call("setProperty", "--p", "-40%")
		//word.Ele.Get("style").Call("setProperty", "--rp", "0%")

		//gsap.Call("to", word.Ele, map[string]interface{}{
		//	"duration": 0,
		//	"--p":      "-40%",
		//	"--rp":     "0%",
		//})
		gsap.Call("to", word.Ele, map[string]interface{}{
			"duration": 0.6,
			"--rcolor": 0.2,
			"onComplete": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				word.Ele.Get("style").Call("setProperty", "--p", fmt.Sprintf("%vpx", 0))
				word.Ele.Get("style").Call("setProperty", "--rp", "0%")
				word.Ele.Get("style").Call("setProperty", "--rcolor", "1")

				word.GsapAnimation.Call("cancel")
				word.GsapAnimation = js.Null()

				word.Ele.Get("style").Set("backgroundPositionX", fmt.Sprintf("%vpx", -getFPX(word.Ele.Get("offsetWidth").Float(), word.Ele.Get("offsetHeight").Float(), fadeRatio)))
				fmt.Println("onComplete")
				for i := 0; i < len(word.HighLightBackgroungAnimation); i++ {
					ite := bgChildren.Index(i)
					word.HighLightBackgroungAnimation[i].Call("kill")
					gsap.Call("set", ite, map[string]interface{}{
						"--p": fmt.Sprintf("%vpx", 0),
					})
					fmt.Println("kill")
				}
				word.HighLightBackgroungAnimation = nil
				return nil
			}),
		})

		word.TextUpAnimation.Call("pause")

		word.Ele.Call("animate", []interface{}{
			map[string]interface{}{},
			map[string]interface{}{
				"transform": "translateY(10px)",
				//"margin-top": "0px",
			},
		},
			map[string]interface{}{
				"duration": 600,
				"easing":   "ease-out",
				"fill":     "forwards",
			},
		)

		for _, ite := range word.HighLightAnimations {
			ite.Call("pause")

		}

		for i := 0; i < bgChildren.Length(); i++ {
			ite := bgChildren.Index(i)
			ite.Call("animate", []interface{}{
				map[string]interface{}{},
				map[string]interface{}{
					"transform": "scale(1) translateY(0px) translateX(0px)",
					//transform: "matrix(1, 0, 0, 1, 0, 0)",
					"textShadow": "none",
					"easing":     "cubic-bezier(0.5, 0, 0.5, 1)",
					//"color":      "rgba(255, 255, 255,0.2)",
					"filter": "blur(0px)",
					//marginLeft: "none"
				},
			},
				map[string]interface{}{
					"duration": 600,
					"easing":   "ease",
					"fill":     "forwards",
				},
			)

		}
		//word.Ele.Get("style").Set("transform", "translateY(10px)")
	}
	//GsetTimeout(func() {
	//	if len(nowPlayingIndex) > 1 {
	//		gd(index+1, lrcs)
	//	}
	//}, 100*time.Millisecond)
	//if index+1 < len(lrcs.Contents) {
	//	isScroll := true
	//	for _, item := range nowPlayingIndex {
	//		if index+1 > item {
	//			isScroll = false
	//		}
	//	}
	//	if isScroll {
	//		if len(nowPlayingIndex) > 1 {
	//			gd(nowPlayingIndex[len(nowPlayingIndex)-1], lrcs, false)
	//			//noePlayingOne = nowPlayingIndex[len(nowPlayingIndex)-1]
	//		} else {
	//			gd(index+1, lrcs, false)
	//			//noePlayingOne = index + 1
	//		}
	//		hasScrolledInRemove = true // 设置标志位，表示已经在 remove 中触发了滚动
	//	}
	//}
}
