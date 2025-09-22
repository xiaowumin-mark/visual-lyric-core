package main

import (
	"fmt"
	"syscall/js"
	"time"

	lyrics "github.com/xiaowumin-mark/visual-lyric-core/lyric"
)

func removeLyric(index int, lrcs *lyrics.Lyrics) {
	haveMoreLine := false
	if len(nowPlayingIndex) > 1 {
		haveMoreLine = true
	}
	removeIndex(index)
	gsap.Call("to", lrcs.Contents[index].Ele, map[string]interface{}{
		"duration": 0.6,
		"scale":    1,
	})
	if len(lrcs.Contents[index].Backgrounds) != 0 {
		lrcs.Contents[index].ShowBackgrounds = false
		for _, item := range lrcs.Contents[index].Backgrounds {
			for _, word := range item.Blocks {
				if word.Begin == 0 && word.End == 0 {
					continue
				}
				cancelWord(word)
			}

			//item.Ele.Call("animate", []interface{}{
			//	map[string]interface{}{},
			//	map[string]interface{}{
			//		"opacity":   "0",
			//		"transform": "scale(0.8) " + item.Ele.Get("style").Get("transform").String(),
			//	},
			//}, map[string]interface{}{
			//	"duration": 600,
			//	"fill":     "forwards",
			//})

			gsap.Call("to", item.Ele, map[string]interface{}{
				"opacity":  0,
				"scale":    0.8,
				"duration": 0.5,
				"ease":     "power4.out",
			})
		}
		GsetTimeout(func() {
			bb := bubbleSort(nowPlayingIndex)
			if len(bb) > 0 {
				gd(bb[0], lrcs, false)
			} else {
				gd(index, lrcs, false)

			}
		}, 0*time.Millisecond)
	}
	words := lrcs.Contents[index].Primary.Blocks

	for _, word := range words {
		if word.Begin == 0 && word.End == 0 {
			continue
		}
		cancelWord(word)
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

	//GsetTimeout(func() {
	/*bb := bubbleSort(nowPlayingIndex)
	if len(bb) > 0 {
		gd(bb[0], lrcs, false)
	} else {
		gd(index, lrcs, false)

	}*/

	if haveMoreLine {
		bb := bubbleSort(nowPlayingIndex)
		if len(bb) > 0 {
			gd(bb[0], lrcs, false)
		} else {
			gd(index, lrcs, false)

		}
	}

	//}, 0*time.Millisecond)
}

func cancelWord(word *lyrics.Block) {
	var bgChildren js.Value = word.Ele.Get("children")

	gsap.Call("to", word.Ele, map[string]interface{}{
		"duration": 0.6,
		"--rcolor": 0.2,
		"onComplete": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			word.Ele.Get("style").Call("setProperty", "--p", fmt.Sprintf("%vpx", 0))
			word.Ele.Get("style").Call("setProperty", "--rp", "0%")
			word.Ele.Get("style").Call("setProperty", "--rcolor", "1")

			//word.GsapAnimation.Call("cancel")
			//word.GsapAnimation = js.Null()
			for _, animation := range word.Animation {
				animation.Call("cancel")
			}
			word.Animation = []js.Value{}

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

	//word.TextUpAnimation.Call("pause")

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
}
